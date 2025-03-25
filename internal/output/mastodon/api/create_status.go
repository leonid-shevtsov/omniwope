package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateStatusRequest struct {
	Status      string   `json:"status"`
	MediaIDs    []string `json:"media_ids,omitempty"`
	InReplyToID string   `json:"in_reply_to_id,omitempty"`
	Visibility  string   `json:"visibility"`
	Language    string   `json:"language"`
	ContentType string   `json:"content_type"`
	Federated   bool     `json:"federated"`
	Boostable   bool     `json:"boostable"`
	Replyable   bool     `json:"replyable"`
	Likeable    bool     `json:"likeable"`
}

type CreateStatusResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// https://docs.joinmastodon.org/methods/statuses/#create
func (c *Client) CreateStatus(payload CreateStatusRequest) (*CreateStatusResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest("POST", c.config.InstanceURL+"/api/v1/statuses", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		panic(fmt.Errorf("bad status code: %s %s", resp.Status, body))
	}

	var statusResponse CreateStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&statusResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &statusResponse, nil
}
