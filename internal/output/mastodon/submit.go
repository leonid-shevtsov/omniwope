package mastodon

import (
	"bytes"
	"io"
	"log/slog"
	"path"

	"github.com/leonid-shevtsov/omniwope/internal/checksum"
	"github.com/leonid-shevtsov/omniwope/internal/content"
	"github.com/leonid-shevtsov/omniwope/internal/linkparser"
	"github.com/leonid-shevtsov/omniwope/internal/output/mastodon/api"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	"github.com/yuin/goldmark/parser"
)

func (o *Output) Submit(post *content.Post) error {
	mastoPost, exists, err := store.Get[Post](o.store, post.URL)
	if err != nil {
		return err
	}
	if exists && mastoPost.Version < VERSION {
		slog.Info("mastodon: Post is older than current version - not updating", "url", post.URL)
		return nil
	}

	slog.Debug("rendering post", "url", post.URL)

	// Even though we use the original markdown for Mastodon, it
	// still needs `relref`s replaced.

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := o.md.Convert(linkparser.PreprocessRefs([]byte(post.Content)), &buf, parser.WithContext(context)); err != nil {
		return err
	}

	body := string(linkparser.UndoRefs(buf.Bytes()))

	if len(post.Resources) > 0 {
		// Prepend the informal "thread marker"
		body = "1/ " + body
	}

	contents := "# " + post.Title + "\n\n" + body

	if len(post.Tags) > 0 {
		contents += "\n\n"
		for i, tag := range post.Tags {
			if i > 0 {
				contents += " "
			}
			contents += "#" + tag
		}
		contents += "\n"
	}

	if !exists {
		// create status
		response, err := o.client.CreateStatus(api.CreateStatusRequest{
			Status:      string(contents),
			Visibility:  o.mastoConfig.Visibility,
			Language:    o.mastoConfig.Language,
			ContentType: "text/markdown",
			Federated:   true,
			Boostable:   true,
			Replyable:   true,
			Likeable:    true,
		})
		if err != nil {
			return err
		}
		mastoPost.ID = response.ID
		mastoPost.URL = response.URL
		mastoPost.Version = VERSION

		if len(post.Resources) > 0 {
			err = o.submitResource(post, mastoPost.ID)
			if err != nil {
				slog.Warn("failed to submit image resource", "err", err)
			}
		}
	} else {
		// update status
		err = o.client.UpdateStatus(mastoPost.ID, api.UpdateStatusRequest{
			Status:      string(contents),
			Language:    o.mastoConfig.Language,
			ContentType: "text/markdown",
		})
		if err != nil {
			return err
		}
	}
	mastoPost.RenderedChecksum = checksum.Sum([]byte(contents))
	err = store.Set(o.store, post.URL, mastoPost)
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) submitResource(post *content.Post, inReplyToID string) error {
	// TODO: allow posting more than 1 image
	resource := post.Resources[0]
	// For now we consider only resource per post, so store key is just the post URL
	mastoPost, exists, err := store.Get[Post](o.resourceStore, post.URL)
	if err != nil {
		return err
	}

	if !exists {
		reader, err := o.config.GetResource(resource.Path)
		if err != nil {
			return err
		}
		defer reader.Close()

		imageName := path.Base(resource.Path)
		imageBytes, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := o.md.Convert([]byte(resource.Label), &buf); err != nil {
			return err
		}

		// 2/ is the informal "thread marker"
		caption := "2/ " + buf.String()

		if o.config.DryRun {
			slog.Info("mastodon: would post image", "image_path", resource.Path, "caption", caption)
			return nil
		}

		mediaID, err := o.client.CreateMedia(imageName, resource.MediaType, imageBytes)
		if err != nil {
			return err
		}

		response, err := o.client.CreateStatus(api.CreateStatusRequest{
			Status:      string(caption),
			MediaIDs:    []string{mediaID},
			InReplyToID: inReplyToID,
			Visibility:  o.mastoConfig.Visibility,
			Language:    o.mastoConfig.Language,
			ContentType: "text/markdown",
			Federated:   true,
			Boostable:   true,
			Replyable:   true,
			Likeable:    true,
		})
		if err != nil {
			return err
		}

		mastoPost.ID = response.ID
		mastoPost.URL = response.URL
		mastoPost.Version = VERSION
		mastoPost.RenderedChecksum = checksum.Sum([]byte(caption)) + ":" + checksum.Sum(imageBytes)
	} else {
		// TODO update resource (not implemented)
		slog.Info("mastodon: updating resource is not implemented")
	}

	err = store.Set(o.resourceStore, post.URL, mastoPost)
	if err != nil {
		return err
	}

	return nil
}
