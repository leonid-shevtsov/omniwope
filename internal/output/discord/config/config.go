package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	BotToken string
	Channel  string // Discord channel URL: https://discord.com/channels/{guild_id}/{channel_id}
}

func Read(viper *viper.Viper) *Config {
	botToken := viper.GetString("discord.bot_token")
	channel := strings.TrimSpace(viper.GetString("discord.channel"))
	if botToken == "" || channel == "" {
		return nil
	}

	return &Config{
		BotToken: botToken,
		Channel:  channel,
	}
}
