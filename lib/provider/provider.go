package provider

import (
	"errors"
	"fmt"
	"time"

	flagsmithClient "github.com/Flagsmith/flagsmith-go-client/v3"
	"github.com/go-logr/logr"
	"github.com/graytonio/flagops/lib/config"
	ld "github.com/launchdarkly/go-server-sdk/v6"
	flagsmith "github.com/open-feature/go-sdk-contrib/providers/flagsmith/pkg"
	fromEnv "github.com/open-feature/go-sdk-contrib/providers/from-env/pkg"
	ofld "github.com/open-feature/go-sdk-contrib/providers/launchdarkly/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

var providerConfigMap = map[config.ProviderType]providerConfig{
	config.Flagsmith:    configureFlagsmithProvider,
	config.FromEnv:      configureFromEnvProvider,
	config.LaunchDarkly: configureLaunchDarklyProvider,
}

func ConfigureProviders(envs map[string]config.Environment) (map[string]*openfeature.Client, error) {
	clients := map[string]*openfeature.Client{}
	for name, env := range envs {
		if cf, ok := providerConfigMap[env.Provider]; !ok {
			return nil, fmt.Errorf("unsupported provider type: %s", env.Provider)
		} else {
			client, err := cf(name, env)
			if err != nil {
				return nil, errors.Join(err, errors.New("error configuring provider"))
			}
			clients[name] = client
		}
	}
	return clients, nil
}

type providerConfig func(name string, env config.Environment) (*openfeature.Client, error)

func configureFlagsmithProvider(name string, env config.Environment) (*openfeature.Client, error) {
	clientOpts := []flagsmithClient.Option{}
	if env.BaseURL != "" {
		clientOpts = append(clientOpts, flagsmithClient.WithBaseURL(env.BaseURL))
	}
	
	client := flagsmithClient.NewClient(env.APIKey, clientOpts...)
	provider := flagsmith.NewProvider(client, flagsmith.WithUsingBooleanConfigValue())
	err := openfeature.SetNamedProvider(fmt.Sprintf("%s-%s", name, env.Provider), provider)
	if err != nil {
		return nil, err
	}
	return openfeature.NewClient(fmt.Sprintf("%s-%s", name, env.Provider)).WithLogger(logr.Discard()), nil
}

func configureFromEnvProvider(name string, env config.Environment) (*openfeature.Client, error) {
	err := openfeature.SetNamedProvider(fmt.Sprintf("%s-%s", name, env.Provider), &fromEnv.FromEnvProvider{})
	if err != nil {
		return nil, err
	}
	return openfeature.NewClient(fmt.Sprintf("%s-%s", name, env.Provider)).WithLogger(logr.Discard()), nil
}

func configureLaunchDarklyProvider(name string, env config.Environment) (*openfeature.Client, error) {
	client, err := ld.MakeClient(env.APIKey, 5*time.Second)
	if err != nil {
		return nil, err
	}

	err = openfeature.SetNamedProvider(fmt.Sprintf("%s-%s", name, env.Provider), ofld.NewProvider(client))
	if err != nil {
		return nil, err
	}

	return openfeature.NewClient(fmt.Sprintf("%s-%s", name, env.Provider)).WithLogger(logr.Discard()), nil
}
