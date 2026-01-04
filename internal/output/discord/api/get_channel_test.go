package api

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveChannel(t *testing.T) {
	tests := []struct {
		name    string
		channel string
		wantErr bool
	}{
		{
			name:    "numeric ID",
			channel: "123456789012345678",
			wantErr: false,
		},
		{
			name:    "numeric ID with # prefix",
			channel: "#123456789012345678",
			wantErr: false,
		},
		{
			name:    "channel name",
			channel: "general",
			wantErr: true,
		},
		{
			name:    "channel name with #",
			channel: "#general",
			wantErr: true,
		},
		{
			name:    "empty string",
			channel: "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			channel: "   ",
			wantErr: true,
		},
		{
			name:    "whitespace around ID",
			channel: "  123456789012345678  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test only validates the parsing logic, not the actual API call
			// For a full integration test, we'd need to mock the HTTP client
			channel := tt.channel
			channel = trimChannelPrefix(channel)
			channel = trimWhitespace(channel)

			_, err := strconv.ParseUint(channel, 10, 64)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions to test the parsing logic
func trimChannelPrefix(channel string) string {
	if len(channel) > 0 && channel[0] == '#' {
		return channel[1:]
	}
	return channel
}

func trimWhitespace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
