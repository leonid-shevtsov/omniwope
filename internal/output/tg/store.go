package tg

type Post struct {
	// ID in channel
	ID int `json:"id"`
	// Version of TG output used
	Version int `json:"version"`
	// Avoid posting if content didn't change
	RenderedChecksum string `json:"checksum"`
}
