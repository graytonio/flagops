package provider

import (
	"errors"
	"fmt"

	flagsmithClient "github.com/Flagsmith/flagsmith-go-client/v3"
	"github.com/go-logr/logr"
	"github.com/graytonio/flagops/lib/config"
	flagsmith "github.com/open-feature/go-sdk-contrib/providers/flagsmith/pkg"
	fromEnv "github.com/open-feature/go-sdk-contrib/providers/from-env/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

var providerConfigMap = map[config.ProviderType]providerConfig{
	config.Flagsmith: configureFlagsmithProvider,
	config.FromEnv:   configureFromEnvProvider,
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
	client := flagsmithClient.NewClient(env.APIKey)
	provider := flagsmith.NewProvider(client, flagsmith.WithUsingBooleanConfigValue())
	openfeature.SetNamedProvider(fmt.Sprintf("%s-%s", name, env.Provider), provider)
	return openfeature.NewClient(fmt.Sprintf("%s-%s", name, env.Provider)).WithLogger(logr.Discard()), nil
}

func configureFromEnvProvider(name string, env config.Environment) (*openfeature.Client, error) {
	openfeature.SetNamedProvider(fmt.Sprintf("%s-%s", name, env.Provider), &fromEnv.FromEnvProvider{})
	return openfeature.NewClient(fmt.Sprintf("%s-%s", name, env.Provider)).WithLogger(logr.Discard()), nil
}
