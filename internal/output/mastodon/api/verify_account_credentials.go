package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Not the full set of fields but we don't need anything else.
type Account struct {
	Acct string `json:"acct"`
	URL  string `json:"url"`
}

// https://docs.joinmastodon.org/methods/accounts/#verify_credentials
func (c *Client) VerifyAccountCredentials() (*Account, error) {
	req, err := http.NewRequest(http.MethodGet, c.config.InstanceURL+"/api/v1/accounts/verify_credentials", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
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

	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &account, nil
}
