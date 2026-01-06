package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*viper.Viper)
		wantNil   bool
		wantToken string
		wantChan  string
	}{
		{
			name: "valid config",
			setup: func(v *viper.Viper) {
				v.Set("discord.bot_token", "test_token")
				v.Set("discord.channel", "123456789")
			},
			wantNil:   false,
			wantToken: "test_token",
			wantChan:  "123456789",
		},
		{
			name: "channel with # prefix",
			setup: func(v *viper.Viper) {
				v.Set("discord.bot_token", "test_token")
				v.Set("discord.channel", "#general")
			},
			wantNil:   false,
			wantToken: "test_token",
			wantChan:  "#general",
		},
		{
			name: "missing bot token",
			setup: func(v *viper.Viper) {
				v.Set("discord.channel", "123456789")
			},
			wantNil: true,
		},
		{
			name: "missing channel",
			setup: func(v *viper.Viper) {
				v.Set("discord.bot_token", "test_token")
			},
			wantNil: true,
		},
		{
			name: "empty bot token",
			setup: func(v *viper.Viper) {
				v.Set("discord.bot_token", "")
				v.Set("discord.channel", "123456789")
			},
			wantNil: true,
		},
		{
			name: "empty channel",
			setup: func(v *viper.Viper) {
				v.Set("discord.bot_token", "test_token")
				v.Set("discord.channel", "")
			},
			wantNil: true,
		},
		{
			name: "channel with whitespace",
			setup: func(v *viper.Viper) {
				v.Set("discord.bot_token", "test_token")
				v.Set("discord.channel", "  123456789  ")
			},
			wantNil:   false,
			wantToken: "test_token",
			wantChan:  "123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			tt.setup(v)

			got := Read(v)

			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tt.wantToken, got.BotToken)
				assert.Equal(t, tt.wantChan, got.Channel)
			}
		})
	}
}

