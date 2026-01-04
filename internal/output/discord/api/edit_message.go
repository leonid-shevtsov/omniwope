package api

import "net/http"

type EditMessageRequest struct {
	Content string   `json:"content,omitempty"`
	Embeds  []Embed  `json:"embeds,omitempty"`
}

type EditMessageResponse struct {
	ID string `json:"id"`
}

// EditMessage edits an existing message in a Discord channel
func (c *Client) EditMessage(channelID string, messageID string, payload EditMessageRequest) (*EditMessageResponse, error) {
	url := BaseURL + "/channels/" + channelID + "/messages/" + messageID
	return doJSONRequest[EditMessageResponse](c, "PATCH", url, payload, http.StatusOK)
}
