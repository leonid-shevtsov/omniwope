package api

import "net/http"

type Embed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int    `json:"color,omitempty"` // Optional: color as integer
}

type CreateMessageRequest struct {
	Content string   `json:"content,omitempty"`
	Embeds  []Embed  `json:"embeds,omitempty"`
}

type CreateMessageResponse struct {
	ID string `json:"id"`
}

// CreateMessage creates a message in a Discord channel
func (c *Client) CreateMessage(channelID string, payload CreateMessageRequest) (*CreateMessageResponse, error) {
	url := BaseURL + "/channels/" + channelID + "/messages"
	return doJSONRequest[CreateMessageResponse](c, "POST", url, payload, http.StatusOK, http.StatusCreated)
}
