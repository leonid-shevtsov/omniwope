package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AccessToken string
	InstanceURL string
	Visibility  string
	Language    string
}

func Read(viper *viper.Viper) *Config {
	accessToken := viper.GetString("mastodon.access_token")
	instanceURL := strings.TrimRight(viper.GetString("mastodon.instance_url"), "/")
	if accessToken == "" || instanceURL == "" {
		return nil
	}

	return &Config{
		AccessToken: accessToken,
		InstanceURL: instanceURL,
		Visibility:  viper.GetString("mastodon.visibility"),
		Language:    viper.GetString("mastodon.language"),
	}
}
