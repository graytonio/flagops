package cmd

import (
	"fmt"
	"os"

	"github.com/graytonio/flagops/lib/config"
	"github.com/graytonio/flagops/lib/provider"
	"github.com/graytonio/flagops/lib/templ"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use: "flagops",
	Short: "Generate files based on the templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}

		conf, err := config.LoadConfig(configPath)
		if err != nil {
		  return err
		}

		providers, err := provider.ConfigureProviders(conf.Envs)
		if err != nil {
		  return err
		}

		engines, err := templ.CreateEngines(conf.Paths, providers)
		if err != nil {
		  return err
		}

		for _, engine := range engines {
			err = engine.Execute()
			if err != nil {
			  return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file (default ./.flagops)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}