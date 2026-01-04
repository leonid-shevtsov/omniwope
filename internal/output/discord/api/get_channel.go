package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Channel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    int    `json:"type"` // 0 = text channel, 2 = voice, etc.
	GuildID string `json:"guild_id,omitempty"` // Guild ID for building message links
}


// ParseChannelURL parses a Discord channel URL and extracts guild_id and channel_id
// Supports formats:
//   - https://discord.com/channels/{guild_id}/{channel_id}
//   - https://discordapp.com/channels/{guild_id}/{channel_id}
// Returns guild_id and channel_id, or error if URL is invalid
func ParseChannelURL(channelURL string) (guildID string, channelID string, err error) {
	channelURL = strings.TrimSpace(channelURL)
	if channelURL == "" {
		return "", "", fmt.Errorf("empty channel URL")
	}

	parsedURL, err := url.Parse(channelURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}

	// Check if it's a Discord channel URL
	host := strings.ToLower(parsedURL.Host)
	if host != "discord.com" && host != "discordapp.com" {
		return "", "", fmt.Errorf("not a Discord URL: expected discord.com or discordapp.com, got %s", host)
	}

	// Parse path: /channels/{guild_id}/{channel_id}
	pathParts := strings.FieldsFunc(parsedURL.Path, func(c rune) bool { return c == '/' })
	if len(pathParts) < 3 || pathParts[0] != "channels" {
		return "", "", fmt.Errorf("invalid Discord channel URL format: expected /channels/{guild_id}/{channel_id}")
	}

	guildID = pathParts[1]
	channelID = pathParts[2]

	// Validate that both are numeric (Discord snowflakes)
	if _, err := strconv.ParseUint(guildID, 10, 64); err != nil {
		return "", "", fmt.Errorf("invalid guild ID in URL: %s", guildID)
	}
	if _, err := strconv.ParseUint(channelID, 10, 64); err != nil {
		return "", "", fmt.Errorf("invalid channel ID in URL: %s", channelID)
	}

	return guildID, channelID, nil
}

// GetChannel retrieves channel information by ID
func (c *Client) GetChannel(channelID string) (*Channel, error) {
	url := BaseURL + "/channels/" + channelID
	channel, err := doJSONRequest[Channel](c, "GET", url, nil, http.StatusOK)
	if err != nil {
		// Provide more context for channel-specific errors
		errStr := err.Error()
		if strings.Contains(errStr, "not found") {
			return nil, fmt.Errorf("channel not found: %s", channelID)
		}
		if strings.Contains(errStr, "forbidden") {
			return nil, fmt.Errorf("forbidden: bot lacks permissions to access channel - ensure the bot has 'View Channels' permission and is allowed in the channel's permission settings")
		}
		return nil, err
	}
	return channel, nil
}

// ResolveChannel resolves a Discord channel URL to channel information
// The URL must be in the format: https://discord.com/channels/{guild_id}/{channel_id}
func (c *Client) ResolveChannel(channelURL string) (*Channel, error) {
	channelURL = strings.TrimSpace(channelURL)
	if channelURL == "" {
		return nil, fmt.Errorf("channel URL is required")
	}

	// Parse the Discord channel URL
	guildID, channelID, err := ParseChannelURL(channelURL)
	if err != nil {
		return nil, fmt.Errorf("invalid channel URL: %w", err)
	}

	// Get channel info and set guild_id
	ch, err := c.GetChannel(channelID)
	if err != nil {
		return nil, err
	}
	// Ensure guild_id is set (it might not be in the API response)
	if ch.GuildID == "" {
		ch.GuildID = guildID
	}
	return ch, nil
}
