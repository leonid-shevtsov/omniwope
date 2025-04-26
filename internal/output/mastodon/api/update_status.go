package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type UpdateStatusRequest struct {
	Status      string   `json:"status"`
	MediaIDs    []string `json:"media_ids,omitempty"`
	Language    string   `json:"language,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
}

// https://docs.joinmastodon.org/methods/statuses/#edit
func (c *Client) UpdateStatus(id string, payload UpdateStatusRequest) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, c.config.InstanceURL+"/api/v1/statuses/"+id, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		if resp.StatusCode == http.StatusUnprocessableEntity && strings.Contains(string(body), "status was not changed") {
			// status remained the same, so the request technically succeeded
			return nil
		}
		panic(fmt.Errorf("bad status code: %s %s", resp.Status, body))
	}

	return nil
}
