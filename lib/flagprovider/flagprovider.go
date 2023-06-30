package flagprovider

import (
	"fmt"

	"github.com/graytonio/kustomize-templater/lib/config"
	"github.com/sirupsen/logrus"
)

var ProviderTable map[string]FeatureFlagProvider

func init() {
	ProviderTable = map[string]FeatureFlagProvider{
		"flagsmith": &FlagsmithFeatureFlagProvider{},
	}
}

func GetProvider(key string, conf config.ProviderConfig, log *logrus.Entry) (FeatureFlagProvider, error) {
	log = log.WithField("provider", key)
	log.Debug("Fetching Flag Provider")
	provider, ok := ProviderTable[key]
	if !ok {
		log.Error("Could not find provider")
		return nil, fmt.Errorf("provider %s is not supported", key)
	}

	log.Debug("Found Provider")
	provider, err := provider.New(conf, log)
	if err != nil {
		return nil, err
	}
	
	log.Debug("Provider Configured")
	return provider, nil
}

type FeatureMap map[string]any

type FeatureFlagProvider interface {
	GetFeatureMap(*logrus.Entry) (FeatureMap, error)
	New(config.ProviderConfig, *logrus.Entry) (FeatureFlagProvider, error)
}
