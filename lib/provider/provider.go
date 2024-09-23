package provider

import (
	"errors"
	"fmt"
	"time"

	flagsmithClient "github.com/Flagsmith/flagsmith-go-client/v3"
	"github.com/go-logr/logr"
	flagopsds "github.com/graytonio/flagops-data-store-openfeature-provider-go/pkg"
	"github.com/graytonio/flagops/lib/config"
	ld "github.com/launchdarkly/go-server-sdk/v6"
	flagsmith "github.com/open-feature/go-sdk-contrib/providers/flagsmith/pkg"
	fromEnv "github.com/open-feature/go-sdk-contrib/providers/from-env/pkg"
	ofld "github.com/open-feature/go-sdk-contrib/providers/launchdarkly/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

type providerConfig func(name string, env config.Environment) (openfeature.FeatureProvider, error)

var providerConfigMap = map[config.ProviderType]providerConfig{
	config.Flagsmith:    buildFlagsmithProvider,
	config.FromEnv:      buildEnvProvider,
	config.LaunchDarkly: buildLaunchDarklyProvider,
}

func ConfigureProviders(envs map[string]config.Environment) (map[string]*openfeature.Client, error) {
	clients := map[string]*openfeature.Client{}
	for name, env := range envs {
		if cf, ok := providerConfigMap[env.Provider]; !ok {
			return nil, fmt.Errorf("unsupported provider type: %s", env.Provider)
		} else {
			provider, err := cf(name, env)
			if err != nil {
				return nil, errors.Join(err, errors.New("error configuring provider"))
			}

			if env.Datastore.Enabled {
				provider, err = flagopsds.NewProvider(env.Datastore.BaseURL, provider)
				if err != nil {
				  return nil, errors.Join(err, errors.New("error configuring datastore"))
				}
			}

			err = openfeature.SetNamedProvider(fmt.Sprintf("%s-%s", name, env.Provider), provider)
			if err != nil {
				return nil, errors.Join(err, errors.New("error setting provider"))
			}
			clients[name] = openfeature.NewClient(fmt.Sprintf("%s-%s", name, env.Provider)).WithLogger(logr.Discard())
		}
	}
	return clients, nil
}

func buildFlagsmithProvider(name string, env config.Environment) (openfeature.FeatureProvider, error) {
	clientOpts := []flagsmithClient.Option{}
	if env.BaseURL != "" {
		clientOpts = append(clientOpts, flagsmithClient.WithBaseURL(env.BaseURL))
	}
	
	client := flagsmithClient.NewClient(env.APIKey, clientOpts...)
	return flagsmith.NewProvider(client, flagsmith.WithUsingBooleanConfigValue()), nil
}

func buildEnvProvider(name string, env config.Environment) (openfeature.FeatureProvider, error) {
	return &fromEnv.FromEnvProvider{}, nil
}

func buildLaunchDarklyProvider(name string, env config.Environment) (openfeature.FeatureProvider, error) {
	client, err := ld.MakeClient(env.APIKey, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return ofld.NewProvider(client), nil
}
