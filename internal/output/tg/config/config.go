package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Credentials string
	ChannelName string
}

func Read(viper *viper.Viper) *Config {
	credentials := viper.GetString("tg.credentials")
	channel := "@" + strings.TrimLeft(viper.GetString("tg.channel"), "@")
	if credentials == "" || channel == "@" {
		return nil
	}

	return &Config{
		Credentials: credentials,
		ChannelName: channel,
	}
}
