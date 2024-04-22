package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/graytonio/flagops/lib/config"
	"github.com/graytonio/flagops/lib/provider"
	"github.com/graytonio/flagops/lib/templ"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	configPath string

	useEnv bool
)

var rootCmd = &cobra.Command{
	Use:   "flagops",
	Short: "Generate files based on the templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.SetReportCaller(true)
		}

		if useEnv {
			return executeEnvConfig()
		}

		return executeConfigFilePaths()
	},
}

func executeEnvConfig() error {
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	providers, err := provider.ConfigureProviders(conf.Envs)
	if err != nil {
		return err
	}

	upsertMode := false
	rawUpsertMode := os.Getenv("FLAGOPS_UPSERT")
	
	if rawUpsertMode != "" {
		upsertMode, err = strconv.ParseBool(os.Getenv("FLAGOPS_UPSERT"))
		if err != nil {
		  return err
		}
	}
	

	source := config.Path{
		Path: os.Getenv("FLAGOPS_SOURCE_PATH"),
		Env: os.Getenv("FLAGOPS_ENVIRONMENT"),
		Identity: os.Getenv("FLAGOPS_IDENTITY"),
		Properties: parseEnvProps(),
		Destination: config.Destination{
			Type: config.DestinationType(os.Getenv("FLAGOPS_DESTINATION_TYPE")),
			Path: os.Getenv("FLAGOPS_DESTINATION_PATH"),
			Repo: os.Getenv("FLAGOPS_DESTINATION_REPO"),
			UpsertMode: upsertMode,
		},
	}

	engine, err := templ.NewTemplateEngine(source, providers[source.Env])
	if err != nil {
	  return err
	}

	return engine.Execute()
}

func parseEnvProps() map[string]any {
	env := os.Environ()
	props := map[string]any{}
	for _, e := range env {
		if !strings.HasPrefix(e, "FLAGOPS_PROP_") {
			continue
		}
		
		parts := strings.Split(e, "=")
		key := strings.ToLower(strings.TrimPrefix(parts[0], "FLAGOPS_PROP_"))
		value := parts[1]
		props[key] = value
	}
	return props
}

func executeConfigFilePaths() error {
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
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file (default ./.flagops)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&useEnv, "use-env", false, "Use env variables to configure path instead of config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
