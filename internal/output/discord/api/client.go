package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/leonid-shevtsov/omniwope/internal/output/discord/config"
)

// Basic client for the Discord API related to posting.
// Not intended as a complete implementation.
type Client struct {
	config *config.Config
}

const BaseURL = "https://discord.com/api/v10"

func NewClient(config *config.Config) *Client {
	return &Client{config: config}
}

// doRequest performs an HTTP request with authentication and common error handling
func (c *Client) doRequest(method, url string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+c.config.BotToken)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	return resp, nil
}

// checkResponse checks the HTTP response status and returns an error if not successful
// For error status codes, it reads and closes the body. For success, it leaves the body for the caller.
func checkResponse(resp *http.Response, allowedStatusCodes ...int) error {
	// Check if status code is in allowed list
	for _, code := range allowedStatusCodes {
		if resp.StatusCode == code {
			return nil
		}
	}

	// Handle error cases - read body for error message
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("bad request: %s", body)
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized: invalid bot token")
	case http.StatusForbidden:
		return fmt.Errorf("forbidden: bot lacks permissions")
	case http.StatusNotFound:
		return fmt.Errorf("not found")
	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")
		return fmt.Errorf("rate limited: retry after %s seconds", retryAfter)
	default:
		return fmt.Errorf("bad status code %d: %s", resp.StatusCode, body)
	}
}

// decodeJSONResponse decodes a JSON response from the HTTP response body
func decodeJSONResponse[T any](resp *http.Response) (*T, error) {
	defer func() { _ = resp.Body.Close() }()
	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result, nil
}

// doJSONRequest performs a JSON request and decodes the response
func doJSONRequest[T any](c *Client, method, url string, payload interface{}, allowedStatusCodes ...int) (*T, error) {
	var body io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewReader(payloadBytes)
	}

	resp, err := c.doRequest(method, url, body, "application/json")
	if err != nil {
		return nil, err
	}

	// Check response - if error, body is already read and closed
	// If success, body is still available for decoding
	if err := checkResponse(resp, allowedStatusCodes...); err != nil {
		return nil, err
	}

	return decodeJSONResponse[T](resp)
}
