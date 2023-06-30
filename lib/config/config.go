package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ProviderConfig struct {
	APIKey string `mapstructure:"apiKey"`
	BaseURL string `mapstructure:"baseUrl"`
}

type Config struct {
	Source string `mapstructure:"src"`
	Destination string `mapstructure:"dest"`
	Recursive bool `mapstructure:"recursive"`
	Filters []string `mapstructure:"filters"`
	ProviderType string `mapstructure:"provider"`
	ProviderConfig ProviderConfig `mapstructure:"providerConfig"`
}

var conf Config

func LoadConfig(configFile string) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".flagops.yaml")
	}

	log.WithField("config_file", viper.ConfigFileUsed()).Debug("Loading Config File")

	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	log.WithField("config_file", viper.ConfigFileUsed()).Debug("Config Loaded")

	err = viper.Unmarshal(&conf)
	cobra.CheckErr(err)

	if conf.Recursive && conf.Destination == "" {
		conf.Destination = "./build"
	}

	log.WithField("config_file", viper.ConfigFileUsed()).Debugf("Current Config: %+v", conf)
}

func GetConfig() *Config {
	return &conf
}