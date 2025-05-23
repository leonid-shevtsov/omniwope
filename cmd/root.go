package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/leonid-shevtsov/omniwope/internal/wope"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "omniwope",
	Short: "Publish your posts to all configured outputs.",
	Long:  `Omniwope - Write Once Publish Everywhere`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := wope.Service{}
		service.Execute()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string
var debug bool
var dryRun bool

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is omniwope.yml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "verbose", false, "enable debug logging")
	err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "dry-run (log changes instead of applying)")
	err = viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))
	if err != nil {
		panic(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("omniwope")
	}

	if inputPath := rootCmd.PersistentFlags().Arg(0); inputPath != "" {
		viper.Set("input", inputPath)
	}

	viper.SetEnvPrefix("omniwope")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// not using a config file is valid
		} else {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	}
}
