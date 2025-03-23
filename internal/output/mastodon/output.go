package mastodon

import (
	"fmt"
	"log/slog"

	"github.com/leonid-shevtsov/omniwope/internal/config"
	"github.com/leonid-shevtsov/omniwope/internal/output/mastodon/api"
	mastoConfig "github.com/leonid-shevtsov/omniwope/internal/output/mastodon/config"
	"github.com/leonid-shevtsov/omniwope/internal/store"
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

	return output, nil
}

func (o *Output) Name() string {
	return "mastodon"
}

func (o *Output) Close() {
	// noop: does not need closing
}
