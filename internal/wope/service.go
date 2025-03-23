package wope

import (
	"log/slog"
	"os"

	"github.com/leonid-shevtsov/omniwope/internal/config"
	"github.com/leonid-shevtsov/omniwope/internal/output"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	"github.com/spf13/viper"
)

type Service struct{}

// Top-level function to
// NOTE: refactor when there are more commands
func (s *Service) Execute() {
	config := config.Read(viper.GetViper())

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: config.LogLevel,
	})))

	outputs, err := output.BuildOutputs(viper.GetViper(), config)
	if err != nil {
		panic(err)
	}
	slog.Debug("Initialized outputs", "output_count", len(outputs))

	checksumStore, err := config.StoreProvider.GetKV("checksums")
	if err != nil {
		panic(err)
	}

	for _, post := range config.Content {
		checksum := post.Checksum()
		storedChecksum, _, err := store.Get[string](checksumStore, post.URL)
		if err != nil {
			panic(err)
		}
		if storedChecksum == checksum {
			// post did not change
			continue
		}

		for _, output := range outputs {
			if !output.StartDate.IsZero() && output.StartDate.After(post.Date) {
				slog.Debug("Skipping post because of start date", "post", post.URL, "output", output.Name())
				continue
			}

			err := output.Submit(&post)
			if err != nil {
				panic(err)
			}
		}

		if !config.DryRun {
			err = store.Set(checksumStore, post.URL, checksum)
			if err != nil {
				panic(err)
			}
		}
	}

	for _, output := range outputs {
		output.Close()
	}

	slog.Info("All done")
}
