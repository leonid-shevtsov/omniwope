package mastodon

import (
	"bytes"
	"log/slog"

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

	contents := "# " + post.Title + "\n\n" + string(linkparser.UndoRefs(buf.Bytes()))

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

	// TODO: handle resources

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
