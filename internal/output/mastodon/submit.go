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

	statusID, err := o.client.CreateStatus(api.CreateStatusRequest{
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

	// TODO: publish resources

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
