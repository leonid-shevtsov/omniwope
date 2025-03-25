package mastodon

type Post struct {
	// ID
	ID string `json:"id"`
	// URL (in Mastodon)
	URL string `json:"url"`
	// Version of Mastodon output used
	Version int `json:"version"`
	// Avoid posting if content didn't change
	RenderedChecksum string `json:"checksum"`
}
