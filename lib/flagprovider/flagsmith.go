package flagprovider

import (
	"errors"
	"fmt"

	"github.com/Flagsmith/flagsmith-go-client/v2"
	"github.com/graytonio/flagops/lib/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type FlagsmithFeatureFlagProvider struct {
    Traits map[string]*flagsmith.Trait
    Identity string

    apiKey string
    client *flagsmith.Client
    mapCache FeatureMap
}

func (fp *FlagsmithFeatureFlagProvider) New(providerConfig config.ProviderConfig, log *logrus.Entry) (FeatureFlagProvider, error) {
    log = log.WithField("provider", "flagsmith")
    if providerConfig.APIKey == "" {
        log.Error("API Key not configured")
        return &FlagsmithFeatureFlagProvider{}, errors.New("apiKey is required for provider flagsmith")
    }
    
    var options []flagsmith.Option
    if providerConfig.BaseURL != "" {
        log.WithField("base_url", providerConfig.BaseURL).Debug("Base URL configured")
        options = append(options, flagsmith.WithBaseURL(providerConfig.BaseURL))
    }
    
    client := flagsmith.NewClient(providerConfig.APIKey, options...)
    log.Debug("Client created")
    return &FlagsmithFeatureFlagProvider{
        apiKey: providerConfig.APIKey,
        client: client,
        Traits: make(map[string]*flagsmith.Trait),
    }, nil
}

func (fp *FlagsmithFeatureFlagProvider) SetIdentity(identity string) {
   fp.Identity = identity
}

func (fp *FlagsmithFeatureFlagProvider) SetTrait(key string, value string) {
    trait, ok := fp.Traits[key]
    if ok {
        trait.TraitValue = value
        return
    }

    fp.Traits[key] = &flagsmith.Trait{
        TraitKey: key,
        TraitValue: value,
    }
}

func (fp *FlagsmithFeatureFlagProvider) UnsetTrait(key string) {
    delete(fp.Traits, key)
}

func (fp *FlagsmithFeatureFlagProvider) DeleteTrait(key string) {
    fp.Traits[key] = &flagsmith.Trait{
        TraitKey: key,
        TraitValue: nil,
    }
}

func (fp *FlagsmithFeatureFlagProvider) GetTraitsSlice() (traits []*flagsmith.Trait) {
    for _, value := range fp.Traits {
        traits = append(traits, value)
    }
    return traits
}

func (fp* FlagsmithFeatureFlagProvider) parseFlagValue(value any) (string) {
    switch v := value.(type){
    case string:
        return v
    default:
        return ""
    }
}

func (fp *FlagsmithFeatureFlagProvider) GetFeatureMap(log *logrus.Entry) (FeatureMap, error) {
    if fp.Identity == "" {
        fp.Identity = "anonymous"
    }

    if len(fp.mapCache) > 0 {
        return fp.mapCache, nil
    }

    log = log.WithField("provider", "flagsmith").WithField("identity", fp.Identity)
    log.Debug("Parsing Feature Flags")
    flags, err := fp.client.GetIdentityFlags(fp.Identity, fp.GetTraitsSlice())
    if err != nil {
        return nil, err
    }

    for _, flag := range flags.AllFlags() {
        log = log.WithField("flag", flag.FeatureName)
        log.Debug("Processing Feature Flag")

        var data map[string]any
        err := yaml.Unmarshal([]byte(fp.parseFlagValue(flag.Value)), &data)
        if err != nil {
            log.WithField("value", flag.Value).Debug("Not valid yaml setting raw text value")
            fp.mapCache[flag.FeatureName] = flag.Value
        } else {
            log.WithField("value", data).Debug("Found valid yaml. Setting map")
            fp.mapCache[flag.FeatureName] = data
        }

       fp.mapCache[fmt.Sprintf("%s_enabled", flag.FeatureName)] = flag.Enabled
       log.WithField("enabled", flag.Enabled).Debug("Added feature flag")
    }

    return fp.mapCache, nil
}