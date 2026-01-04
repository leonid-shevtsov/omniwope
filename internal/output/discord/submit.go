package discord

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"path"
	"strings"

	"github.com/leonid-shevtsov/omniwope/internal/checksum"
	"github.com/leonid-shevtsov/omniwope/internal/content"
	"github.com/leonid-shevtsov/omniwope/internal/hashtags"
	"github.com/leonid-shevtsov/omniwope/internal/linkparser"
	"github.com/leonid-shevtsov/omniwope/internal/output/discord/api"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	"github.com/samber/lo"
	"github.com/yuin/goldmark/parser"
)

const (
	maxMessageLength    = 2000
	maxEmbedDescription = 4096
	maxEmbedTitle       = 256
	maxFileSizeBytes    = 25 * 1024 * 1024 // 25MB
)

func (o *Output) Submit(post *content.Post) error {
	existingPost, exists, err := store.Get[Post](o.store, post.URL)
	if err != nil {
		return err
	}
	if exists && existingPost.Version < VERSION {
		slog.Info("discord: Post is older than current version - not updating", "url", post.URL)
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

	// Truncate if too long (Discord has limits) - use lo.Ellipsis for UTF-8 safety
	description := lo.Ellipsis(contents, maxEmbedDescription)
	if len(contents) > maxEmbedDescription {
		slog.Warn("discord: content truncated to fit embed description limit", "url", post.URL)
	}

	// Truncate title if needed (Discord limit is 256 characters)
	embedTitle := lo.Ellipsis(post.Title, maxEmbedTitle)
	if len(post.Title) > maxEmbedTitle {
		slog.Warn("discord: title truncated to fit embed title limit", "url", post.URL)
	}

	// Create embed
	embed := api.Embed{
		Title:       embedTitle,
		Description: description,
	}

	// Handle resources
	var resource *content.Resource
	if len(post.Resources) > 0 {
		resource = &post.Resources[0]
	}

	if !exists {
		return o.createPost(post, embed, resource)
	} else {
		return o.updatePost(existingPost, post, embed, resource)
	}
}

func (o *Output) createPost(post *content.Post, embed api.Embed, resource *content.Resource) error {
	if o.config.DryRun {
		slog.Info("discord: would CREATE post", "url", post.URL, "title", embed.Title, "has_resource", resource != nil)
		return nil
	}

	if resource != nil {
		// Upload file with message
		reader, err := o.config.GetResource(resource.Path)
		if err != nil {
			return err
		}
		defer reader.Close()

		resourceBytes, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		// Validate file size (Discord limit is 25MB)
		if len(resourceBytes) > maxFileSizeBytes {
			return fmt.Errorf("file size %d bytes exceeds Discord limit of %d bytes", len(resourceBytes), maxFileSizeBytes)
		}

		// Render resource label if present and include in embed description
		if resource.Label != "" {
			var labelBuf bytes.Buffer
			if err := o.md.Convert([]byte(resource.Label), &labelBuf); err != nil {
				return err
			}
			labelContent := labelBuf.String()
			if labelContent != "" {
				combinedDescription := embed.Description + "\n\n" + labelContent
				embed.Description = lo.Ellipsis(combinedDescription, maxEmbedDescription)
			}
		}

		// Determine content type
		contentType := resource.MediaType
		if contentType == "" {
			contentType = api.GetContentType(resource.Path)
		}

		fileName := path.Base(resource.Path)
		response, err := o.client.UploadFile(o.channelID, api.UploadFileRequest{
			Embeds: []api.Embed{embed},
			Files: []api.FileAttachment{
				{
					Name:        fileName,
					Content:     resourceBytes,
					ContentType: contentType,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}

		// Store the post
		checksumValue := calculateChecksum(embed, resourceBytes)
		err = store.Set[Post](o.store, post.URL, Post{
			ID:               response.ID,
			RenderedChecksum: checksumValue,
			Version:          VERSION,
		})
		if err != nil {
			return err
		}
	} else {
		// Post message without file
		response, err := o.client.CreateMessage(o.channelID, api.CreateMessageRequest{
			Embeds: []api.Embed{embed},
		})
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}

		checksumValue := calculateChecksum(embed, nil)
		err = store.Set[Post](o.store, post.URL, Post{
			ID:               response.ID,
			RenderedChecksum: checksumValue,
			Version:          VERSION,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// calculateChecksum calculates the checksum for a post, including resource if present
func calculateChecksum(embed api.Embed, resourceBytes []byte) string {
	baseChecksum := checksum.Sum([]byte(embed.Title + embed.Description))
	if len(resourceBytes) > 0 {
		return baseChecksum + ":" + checksum.Sum(resourceBytes)
	}
	return baseChecksum
}

func (o *Output) updatePost(existingPost Post, post *content.Post, embed api.Embed, resource *content.Resource) error {
	var resourceBytes []byte
	if resource != nil {
		reader, err := o.config.GetResource(resource.Path)
		if err != nil {
			return err
		}
		defer reader.Close()

		resourceBytes, err = io.ReadAll(reader)
		if err != nil {
			return err
		}
	}
	renderedChecksum := calculateChecksum(embed, resourceBytes)

	if existingPost.RenderedChecksum == renderedChecksum {
		slog.Info("Post is unchanged - skipping", "url", post.URL)
		return nil
	}

	if o.config.DryRun {
		slog.Info("discord: would UPDATE post", "url", post.URL, "title", embed.Title)
		return nil
	}

	// Discord doesn't support editing messages with file attachments
	// If the resource changed, we'd need to delete and recreate, but for MVP
	// we'll just update the embed content
	_, err := o.client.EditMessage(o.channelID, existingPost.ID, api.EditMessageRequest{
		Embeds: []api.Embed{embed},
	})
	if err != nil {
		// Check if it's a "message not modified" error
		if strings.Contains(err.Error(), "message is not modified") {
			slog.Debug("discord: message is not modified", "url", post.URL)
		} else {
			return fmt.Errorf("failed to edit message: %w", err)
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
