package cmd

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/graytonio/flagops/lib/config"
	"github.com/graytonio/flagops/lib/filters"
	"github.com/graytonio/flagops/lib/flagprovider"
	"github.com/noirbizarre/gonja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO Add provider setting to command flags
func init() {
	rootCmd.AddCommand(generateCommand)
	generateCommand.Flags().StringP("src", "s", "", "Input file to generate or directory to start from is using recursive")
	viper.BindPFlag("src", generateCommand.Flags().Lookup("src"))

	generateCommand.Flags().BoolP("recursive", "r", false, "Enable recursive mode to apply templating to entire directory")
	viper.BindPFlag("recursive", generateCommand.Flags().Lookup("recursive"))

	generateCommand.Flags().StringP("dest", "d", "", "Write template to file instead of stdout. If recursive is enabled the output base directory (default: build)")
	viper.BindPFlag("dest", generateCommand.Flags().Lookup("dest"))

	generateCommand.Flags().StringP("provider", "p", "file", "Provider type to pull context from")
	viper.BindPFlag("provider", generateCommand.Flags().Lookup("provider"))

	generateCommand.Flags().StringArrayP("filter", "f", []string{}, "Apply a filter over the produced template")
	viper.BindPFlag("filters", generateCommand.Flags().Lookup("filter"))
}

func ensureFilePath(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) { 
		os.MkdirAll(path.Dir(filePath), 0700) // Create your file
	}
}

func getDestinationFilePath(sourceFile string, conf *config.Config) string {
	if conf.Destination == "" {
		return ""
	}

	if !conf.Recursive {
		return conf.Destination
	}

	trimmedSourceFile := strings.TrimPrefix(sourceFile, conf.Source)
	return path.Join(conf.Destination, trimmedSourceFile)
}

func generateFile(provider flagprovider.FeatureFlagProvider, sourceFile string, log *logrus.Entry) error {
	log = log.WithField("path", sourceFile)

	conf := config.GetConfig()
	tpl, err := gonja.FromFile(sourceFile)
	if err != nil {
		return err
	}

	flags, err := provider.GetFeatureMap(log)
	if err != nil {
		return err
	}

	log.Debug("Executing Template")
	output, err := tpl.Execute(flags)
	if err != nil {
		return err
	}

	log.Debug("Applying Filters")
	output, err = filters.FilterString(output, conf.Filters, log)
	if err != nil {
		return err
	}

	var destFile *os.File = os.Stdout

	destFilePath := getDestinationFilePath(sourceFile, conf)
	log.WithField("destination_file_path", destFilePath).Debug("Calculating Destination")
	if destFilePath != "" {
		log.WithField("destination_file_path", destFilePath).Debug("Writing to output file")
		ensureFilePath(destFilePath)
		destFile, err = os.OpenFile(destFilePath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
	}

	log.WithField("dest", destFile.Name()).Debug("Writing output")
	_, err = io.WriteString(destFile, output)
	return err
}

var generateCommand = &cobra.Command{
	Use: "generate",
	Short: "Run template engine over files",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()

		log := logrus.WithFields(logrus.Fields{
			"recursive": conf.Recursive,
			"provider": conf.ProviderType,
			"source": conf.Source,
			"dest": conf.Destination,
		})
		provider, err := flagprovider.GetProvider(conf.ProviderType, conf.ProviderConfig, log)
		cobra.CheckErr(err)

		if !conf.Recursive {
			err = generateFile(provider, conf.Source, log)
			cobra.CheckErr(err)
			return
		}

		err = filepath.WalkDir(conf.Source, func(filePath string, d fs.DirEntry, err error) error {
			logrus.WithField("path", filePath).Debug("Processing Path")
			if d.IsDir() {
				return nil
			}

			// TODO Support source file glob
			return generateFile(provider, filePath, log)
		})
		cobra.CheckErr(err)
	},
}