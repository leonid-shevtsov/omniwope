package output

import (
	"github.com/leonid-shevtsov/omniwope/internal/config"
	"github.com/leonid-shevtsov/omniwope/internal/content"
	"github.com/leonid-shevtsov/omniwope/internal/output/tg"
	tgConfig "github.com/leonid-shevtsov/omniwope/internal/output/tg/config"
	"github.com/spf13/viper"
)

type Output interface {
	Submit(post *content.Post) error
	Close()
}

func BuildOutputs(viper *viper.Viper, config *config.Config) ([]Output, error) {
	var outputs []Output

	if tgConfig := tgConfig.Read(viper); tgConfig != nil {
		tgOutput, err := tg.NewOutput(config, tgConfig)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, tgOutput)
	}

	return outputs, nil
}
