package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
)

type UploadFileRequest struct {
	Content string           `json:"content,omitempty"`
	Embeds  []Embed          `json:"embeds,omitempty"`
	Files   []FileAttachment `json:"-"`
}

type FileAttachment struct {
	Name        string
	Content     []byte
	ContentType string
}

type UploadFileResponse struct {
	ID string `json:"id"`
}

// UploadFile creates a message with file attachments in a Discord channel
func (c *Client) UploadFile(channelID string, payload UploadFileRequest) (*UploadFileResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add payload_json field with the message content/embeds
	payloadJSON := map[string]interface{}{}
	if payload.Content != "" {
		payloadJSON["content"] = payload.Content
	}
	if len(payload.Embeds) > 0 {
		payloadJSON["embeds"] = payload.Embeds
	}

	if len(payloadJSON) > 0 {
		payloadBytes, err := json.Marshal(payloadJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		payloadField, err := writer.CreateFormField("payload_json")
		if err != nil {
			return nil, fmt.Errorf("failed to create form field: %w", err)
		}
		_, err = payloadField.Write(payloadBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to write form field: %w", err)
		}
	}

	// Add file attachments
	for i, file := range payload.Files {
		part, err := writer.CreateFormFile(fmt.Sprintf("files[%d]", i), file.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		_, err = part.Write(file.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	resp, err := c.doRequest("POST", BaseURL+"/channels/"+channelID+"/messages", &body, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp, http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}

	return decodeJSONResponse[UploadFileResponse](resp)
}

// Helper function to determine content type from file extension
func GetContentType(filename string) string {
	ext := path.Ext(filename)
	// Convert to lowercase for case-insensitive matching
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}
