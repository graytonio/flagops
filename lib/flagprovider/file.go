package flagprovider

import (
	"errors"
	"os"
	"strings"

	"github.com/graytonio/flagops/lib/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type FileFeatureFlagProvider struct {
	path string
	mapCache FeatureMap
}

func (fp *FileFeatureFlagProvider) New(providerConfig config.ProviderConfig, log *logrus.Entry) (FeatureFlagProvider, error) {
	log = log.WithField("provider", "file")
	if providerConfig.Path == "" {
		log.Error("Path not configured")
		return &FileFeatureFlagProvider{}, errors.New("path is required for provider file")
	}

	return &FileFeatureFlagProvider{
		path: providerConfig.Path,
		mapCache: make(FeatureMap),
	}, nil
}

func (fp *FileFeatureFlagProvider) GetFeatureMap(log *logrus.Entry) (FeatureMap, error) {
	log = log.WithField("provider", "file")
	data, err := os.ReadFile(fp.path)
	if err != nil {
		log.Error("Could not read file")
		return nil, err
	}

	switch {
	case strings.HasSuffix(fp.path, "yaml"):
		fp.mapCache, err = parseYamlFile(data, log.WithField("file_type", "yaml"))
	default:
		err = errors.New("unsupported file type")
	}

	if err != nil {
		log.Error("Could not parse flag file")
		return nil, err
	}

	return fp.mapCache, nil
}

func parseYamlFile(data []byte, log *logrus.Entry) (FeatureMap, error) {
	featureMap := make(FeatureMap)
	log.Debug("Unmarshalling data")

	err := yaml.Unmarshal(data, &featureMap)
	if err != nil {
		log.WithField("error", err).Error("Could not unmarshall yaml file")
		return nil, err
	}

	return featureMap, nil
}