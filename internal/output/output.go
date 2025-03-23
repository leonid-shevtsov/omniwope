package output

import (
	"time"

	"github.com/leonid-shevtsov/omniwope/internal/config"
	"github.com/leonid-shevtsov/omniwope/internal/content"
	"github.com/leonid-shevtsov/omniwope/internal/output/mastodon"
	mastoConfig "github.com/leonid-shevtsov/omniwope/internal/output/mastodon/config"
	"github.com/leonid-shevtsov/omniwope/internal/output/tg"
	tgConfig "github.com/leonid-shevtsov/omniwope/internal/output/tg/config"
	"github.com/spf13/viper"
)

type Output interface {
	Name() string
	Submit(post *content.Post) error
	Close()
}

type OutputConfig struct {
	Output
	StartDate time.Time
}

func BuildOutputs(viper *viper.Viper, config *config.Config) ([]OutputConfig, error) {
	var outputs []OutputConfig

	if tgConfig := tgConfig.Read(viper); tgConfig != nil {
		tgOutput, err := tg.NewOutput(config, tgConfig)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, OutputConfig{
			Output:    tgOutput,
			StartDate: viper.GetTime("tg.start_date"),
		})
	}

	if mastoConfig := mastoConfig.Read(viper); mastoConfig != nil {
		mastoOutput, err := mastodon.NewOutput(config, mastoConfig)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, OutputConfig{
			Output:    mastoOutput,
			StartDate: viper.GetTime("mastodon.start_date"),
		})
	}

	return outputs, nil
}
