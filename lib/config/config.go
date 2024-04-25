package config

import (
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("/config")
	v.SetConfigName(".flagops")
	v.SetConfigType("yaml")

	if path != "" {
		v.SetConfigFile(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	conf := Config{}
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
