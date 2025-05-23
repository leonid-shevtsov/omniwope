package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type mediaResponse struct {
	ID string `json:"id"`
}

// https://docs.joinmastodon.org/methods/media/#v2
func (c *Client) CreateMedia(filename string, contentType string, contents []byte) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// adapted from multipart.Writer#CreateFormFile
	// how else to create a file part with a proper content type?
	// this SO answer agrees: https://stackoverflow.com/a/21133168/6678
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(filename)))
	header.Set("Content-Type", contentType)

	fileWriter, err := writer.CreatePart(header)
	if err != nil {
		return "", fmt.Errorf("failed to open form writer: %w", err)
	}
	_, err = fileWriter.Write(contents)
	if err != nil {
		return "", fmt.Errorf("failed to write file to form: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close form writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.InstanceURL+"/api/v2/media", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}
		panic(fmt.Errorf("bad status code: %s %s", resp.Status, body))
	}

	var mediaResponse mediaResponse
	err = json.NewDecoder(resp.Body).Decode(&mediaResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return mediaResponse.ID, nil
}

// copied from mime/multipart package
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
