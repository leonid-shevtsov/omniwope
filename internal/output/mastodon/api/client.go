package api

import "github.com/leonid-shevtsov/omniwope/internal/output/mastodon/config"

// Very basic client for the Mastodon API related to posting.
// Not intended as a complete implementation.
type Client struct {
	config *config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{config: config}
}
