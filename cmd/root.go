package cmd

import (
	"github.com/graytonio/flagops/lib/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "flagops",
		Short: "A cli tool for using feature flags to control kustomize output",
	}
)

func init() {
	rootCmd.PersistentFlags().String("config", "", "Path to configuration file")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging")
}

func Execute() error {

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		configFile, err := rootCmd.PersistentFlags().GetString("config")
		if err != nil {
			return err
		}
	
		verbose, err := rootCmd.PersistentFlags().GetBool("verbose")
		if err != nil {
			return err
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
	
		config.LoadConfig(configFile)
		return nil
	}
	return rootCmd.Execute()
}

