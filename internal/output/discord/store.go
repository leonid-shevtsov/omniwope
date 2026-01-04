package discord

type Post struct {
	// ID in Discord channel (snowflake as string)
	ID string `json:"id"`
	// Version of Discord output used
	Version int `json:"version"`
	// Avoid posting if content didn't change
	RenderedChecksum string `json:"checksum"`
}
