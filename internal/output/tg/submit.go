package tg

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/leonid-shevtsov/omniwope/internal/checksum"
	"github.com/leonid-shevtsov/omniwope/internal/content"
	"github.com/leonid-shevtsov/omniwope/internal/hashtags"
	"github.com/leonid-shevtsov/omniwope/internal/linkparser"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	"github.com/yuin/goldmark/parser"
)

func (o *Output) Submit(post *content.Post) error {
	existingPost, exists, err := store.Get[Post](o.store, post.URL)
	if err != nil {
		return err
	}
	if exists && existingPost.Version < VERSION {
		slog.Info("tg: Post is older than current version - not updating", "url", post.URL)
		return nil
	}

	slog.Debug("rendering post", "url", post.URL)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := o.md.Convert(linkparser.PreprocessRefs([]byte(post.Content)), &buf, parser.WithContext(context)); err != nil {
		return err
	}

	contents := string(linkparser.UndoRefs(buf.Bytes()))
	contents = hashtags.Insert(post.Tags, contents)

	var titleBuf, titleOutBuf bytes.Buffer
	titleBuf.WriteString("# ")
	titleBuf.WriteString(post.Title)
	if err := o.md.Convert(titleBuf.Bytes(), &titleOutBuf); err != nil {
		return err
	}
	contents = fmt.Sprintf("%s\n%s", titleOutBuf.String(), contents)

	// Currently the resource is posted separately
	if len(post.Resources) > 0 {
		err := o.submitResource(post)
		if err != nil {
			return err
		}
	}

	if !exists {
		return o.createPost(post, contents)
	} else {
		return o.updatePost(existingPost, post, contents)
	}
}

func (o *Output) submitResource(post *content.Post) error {
	// TODO: allow posting more than 1 image
	resource := post.Resources[0]
	// For now we consider only resource per post, so store key is just the post URL
	_, exists, err := store.Get[Post](o.resourceStore, post.URL)
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
			panic(err)
		}
		captionHTML := buf.String()

		if o.config.DryRun {
			slog.Info("tg: would post image", "image_path", resource.Path, "caption", captionHTML)
			return nil
		}

		var msgConfig tgbotapi.Chattable
		if resource.IsVideo() {
			video := tgbotapi.NewVideo(o.channelID, tgbotapi.FileBytes{Name: imageName, Bytes: imageBytes})
			video.Caption = captionHTML
			video.ParseMode = "HTML"
			msgConfig = video
		} else {
			photo := tgbotapi.NewPhoto(o.channelID, tgbotapi.FileBytes{Name: imageName, Bytes: imageBytes})
			photo.Caption = captionHTML
			photo.ParseMode = "HTML"
			msgConfig = photo
		}

		msg, err := o.bot.Send(msgConfig)
		if err != nil {
			panic(err)
		}

		store.Set[Post](o.resourceStore, post.URL, Post{
			ID:               msg.MessageID,
			RenderedChecksum: checksum.Sum([]byte(captionHTML)) + ":" + checksum.Sum(imageBytes),
			Version:          VERSION,
		})
	} else {
		// TODO update resource (not implemented)
	}

	return nil
}

func (o *Output) createPost(post *content.Post, renderedContents string) error {
	if o.config.DryRun {
		slog.Info("tg: would CREATE post", "url", post.URL, "contents", renderedContents)
		return nil
	}

	msgConfig := tgbotapi.NewMessage(o.channelID, renderedContents)
	msgConfig.ParseMode = "HTML"
	msgConfig.DisableWebPagePreview = true

	msg, err := o.bot.Send(msgConfig)
	if err != nil {
		return err
	}

	err = store.Set[Post](o.store, post.URL, Post{
		ID:               msg.MessageID,
		RenderedChecksum: checksum.Sum([]byte(renderedContents)),
		Version:          VERSION,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) updatePost(existingPost Post, post *content.Post, renderedContents string) error {
	renderedChecksum := checksum.Sum([]byte(renderedContents))

	if existingPost.RenderedChecksum == renderedChecksum {
		slog.Info("Post is unchanged - skipping", "url", post.URL)
		return nil
	}

	if o.config.DryRun {
		slog.Info("tg: would UPDATE post", "url", post.URL, "contents", renderedContents)
		os.WriteFile("debug_output/"+path.Base(post.URL), []byte(renderedContents), 0644)
		return nil
	}

	msgConfig := tgbotapi.NewEditMessageText(o.channelID, existingPost.ID, renderedContents)
	msgConfig.ParseMode = "HTML"
	msgConfig.DisableWebPagePreview = true

	_, err := o.bot.Send(msgConfig)
	if err != nil {
		if !strings.Contains(err.Error(), "message is not modified") {
			slog.Error("Rendered checksum was new, but TG responded with message is not modified", "url", post.URL)
		} else {
			return err
		}
	}

	err = store.Set[Post](o.store, post.URL, Post{
		ID:               existingPost.ID,
		RenderedChecksum: renderedChecksum,
		Version:          VERSION,
	})
	if err != nil {
		return err
	}

	return nil
}
