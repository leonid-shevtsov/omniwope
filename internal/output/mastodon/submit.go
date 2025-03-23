package mastodon

import (
	"log/slog"

	"github.com/leonid-shevtsov/omniwope/internal/checksum"
	"github.com/leonid-shevtsov/omniwope/internal/content"
	"github.com/leonid-shevtsov/omniwope/internal/output/mastodon/api"
	"github.com/leonid-shevtsov/omniwope/internal/store"
)

func (o *Output) Submit(post *content.Post) error {
	existingPost, exists, err := store.Get[Post](o.store, post.URL)
	if err != nil {
		return err
	}
	if exists && existingPost.Version < VERSION {
		slog.Info("mastodon: Post is older than current version - not updating", "url", post.URL)
		return nil
	}

	slog.Debug("rendering post", "url", post.URL)

	contents := "# " + post.Title + "\n\n" + post.Content

	// TODO: handle resources

	var statusID string

	if !exists {
		// create status
		statusID, err = o.client.CreateStatus(api.CreateStatusRequest{
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
	} else {
		// update status
		statusID = existingPost.ID
		err = o.client.UpdateStatus(statusID, api.UpdateStatusRequest{
			Status:      string(contents),
			Language:    o.mastoConfig.Language,
			ContentType: "text/markdown",
		})
		if err != nil {
			return err
		}
	}
	err = store.Set(o.store, post.URL, Post{
		ID:               statusID,
		RenderedChecksum: checksum.Sum([]byte(contents)),
		Version:          VERSION,
	})
	if err != nil {
		return err
	}
	return nil
}
