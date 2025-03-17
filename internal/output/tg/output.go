package tg

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/leonid-shevtsov/omniwope/internal/config"
	"github.com/leonid-shevtsov/omniwope/internal/linkparser"
	tgConfig "github.com/leonid-shevtsov/omniwope/internal/output/tg/config"
	"github.com/leonid-shevtsov/omniwope/internal/output/tg/telegold"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

const VERSION = 2

type Output struct {
	bot           *tgbotapi.BotAPI
	store         store.KV
	resourceStore store.KV
	channelID     int64
	channelName   string
	config        *config.Config
	md            goldmark.Markdown
}

func NewOutput(config *config.Config, tgConfig *tgConfig.Config) (*Output, error) {
	store, err := config.StoreProvider.GetKV("tg")
	if err != nil {
		return nil, err
	}
	resourceStore, err := config.StoreProvider.GetKV("tg_resource_posts")
	if err != nil {
		return nil, err
	}

	bot, err := tgbotapi.NewBotAPI(tgConfig.Credentials)
	if err != nil {
		return nil, err
	}
	slog.Debug("tg: authorized", "user_name", bot.Self.UserName)

	bot.Debug = config.LogLevel <= slog.LevelDebug

	channel, err := bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{SuperGroupUsername: tgConfig.ChannelName},
	})
	if err != nil {
		return nil, err
	}
	slog.Debug("tg: posting to channel", "channel_name", channel.UserName)

	output := &Output{
		bot:           bot,
		channelID:     channel.ID,
		channelName:   channel.UserName,
		store:         store,
		resourceStore: resourceStore,
		config:        config,
	}
	output.buildMarkdown()

	return output, nil
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
			return fmt.Sprintf("https://t.me/%s/%d", o.channelName, postInfo.ID)
		},
	)

	o.md = goldmark.New(
		goldmark.WithRenderer(telegold.NewRenderer()),
		goldmark.WithParserOptions(parser.WithASTTransformers(util.Prioritized(refTransformer, 0))),
	)
}
