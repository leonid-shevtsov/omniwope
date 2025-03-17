package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/leonid-shevtsov/omniwope/internal/content"
	"github.com/leonid-shevtsov/omniwope/internal/store"
	jsonStore "github.com/leonid-shevtsov/omniwope/internal/store/json"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel slog.Level
	DryRun   bool

	GetResource   func(string) (io.ReadCloser, error)
	RefNameToURL  func(string) string
	Content       []content.Post
	StoreProvider store.Provider
}

func Read(viper *viper.Viper) *Config {
	refPattern := viper.GetString("relref.pattern")
	if refPattern == "" {
		refPattern = "%s"
	}
	resourceBasePath := viper.GetString("resources.base_path")

	logLevel := slog.LevelWarn
	if viper.GetBool("verbose") {
		logLevel = slog.LevelDebug
	}

	return &Config{
		LogLevel: logLevel,
		DryRun:   viper.GetBool("dry_run"),
		RefNameToURL: func(refName string) string {
			return fmt.Sprintf(refPattern, refName)
		},
		GetResource: func(resourcePath string) (io.ReadCloser, error) {
			return os.Open(path.Join(resourceBasePath, resourcePath))
		},
		Content:       readInput(viper),
		StoreProvider: jsonStore.NewProvider(viper.GetString("store.path")),
	}
}

func readInput(viper *viper.Viper) []content.Post {
	var inputFile io.Reader
	if inputPath := viper.GetString("input"); inputPath != "" {
		file, err := os.Open(inputPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		inputFile = file
	} else {
		inputFile = os.Stdin
	}

	var content []content.Post
	err := json.NewDecoder(inputFile).Decode(&content)
	if err != nil {
		panic(err)
	}

	return content
}
