package discord

import (
	"fmt"
	"log/slog"

	"github.com/leonid-shevtsov/omniwope/internal/config"
	"github.com/leonid-shevtsov/omniwope/internal/linkparser"
	"github.com/leonid-shevtsov/omniwope/internal/output/discord/api"
	discordConfig "github.com/leonid-shevtsov/omniwope/internal/output/discord/config"
	"github.com/leonid-shevtsov/omniwope/internal/output/discord/discordmarkdown"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

const VERSION = 1

type Output struct {
	store        store.KV
	config       *config.Config
	discordConfig *discordConfig.Config
	client       *api.Client
	channelID    string
	guildID      string
	md           goldmark.Markdown
}

func NewOutput(config *config.Config, discordConfig *discordConfig.Config) (*Output, error) {
	store, err := config.StoreProvider.GetKV("discord")
	if err != nil {
		return nil, err
	}

	client := api.NewClient(discordConfig)

	// Resolve channel URL - this validates and returns channel info
	channel, err := client.ResolveChannel(discordConfig.Channel)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve channel: %w", err)
	}
	slog.Debug("discord: setup successful", "channel_id", channel.ID, "channel_name", channel.Name, "guild_id", channel.GuildID)

	output := &Output{
		store:        store,
		config:       config,
		discordConfig: discordConfig,
		client:       client,
		channelID:    channel.ID,
		guildID:      channel.GuildID,
	}
	output.buildMarkdown()

	return output, nil
}

func (o *Output) Name() string {
	return "discord"
}

func (o *Output) Close() {
	// noop: does not need closing
}

func (o *Output) buildMarkdown() {
	refTransformer := linkparser.NewRefTransformer(
		func(refName string) string {
			return o.config.RefNameToURL(refName)
		},
		func(refName string) string {
			url := o.config.RefNameToURL(refName)
			postInfo, exists, err := store.Get[Post](o.store, url)
			if err != nil {
				panic(err)
			}
			if !exists {
				slog.Error("missing post in mapping - not replacing", "ref_name", refName)
				return url
			}
			// Discord message link format: https://discord.com/channels/{guild_id}/{channel_id}/{message_id}
			if o.guildID != "" {
				return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", o.guildID, o.channelID, postInfo.ID)
			}
			// Fallback to URL if guild_id is not available (DM channels don't have guild_id)
			return url
		},
	)

	o.md = goldmark.New(
		goldmark.WithRenderer(discordmarkdown.NewRenderer()),
		goldmark.WithParserOptions(parser.WithASTTransformers(util.Prioritized(refTransformer, 0))),
	)
}
