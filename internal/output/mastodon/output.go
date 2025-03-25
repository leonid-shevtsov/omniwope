package mastodon

import (
	"fmt"
	"log/slog"

	"github.com/leonid-shevtsov/omniwope/internal/config"
	"github.com/leonid-shevtsov/omniwope/internal/linkparser"
	"github.com/leonid-shevtsov/omniwope/internal/output/mastodon/api"
	mastoConfig "github.com/leonid-shevtsov/omniwope/internal/output/mastodon/config"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

const VERSION = 1

// Supports any server compatible with the Mastodon API, like GoToSocial
// TODO: list specific endpoints
// TODO: list specific servers
type Output struct {
	store       store.KV
	config      *config.Config
	mastoConfig *mastoConfig.Config
	client      *api.Client
	md          goldmark.Markdown
}

func NewOutput(config *config.Config, mastoConfig *mastoConfig.Config) (*Output, error) {
	store, err := config.StoreProvider.GetKV("mastodon")
	if err != nil {
		return nil, err
	}

	client := api.NewClient(mastoConfig)
	account, err := client.VerifyAccountCredentials()
	if err != nil {
		return nil, fmt.Errorf("the access token did not work: %w", err)
	}
	slog.Debug("mastodon: setup successful", "account", account.Acct, "url", account.URL)

	output := &Output{
		store:       store,
		config:      config,
		mastoConfig: mastoConfig,
		client:      client,
	}
	output.buildMarkdown()

	return output, nil
}

func (o *Output) Name() string {
	return "mastodon"
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
			return postInfo.URL
		},
	)

	o.md = goldmark.New(
		// TODO: because Mastodon doesn't support text/markdown posts, actually - only GoToSocial does -
		// need to be able to render into text instead. The main challenge would be to render links as nice footnotes.
		goldmark.WithRenderer(markdown.NewRenderer()),
		goldmark.WithParserOptions(parser.WithASTTransformers(util.Prioritized(refTransformer, 0))),
	)
}
