
# Template Files

Flagops uses the golang templating syntax however it uses `[{` and `}]` as deliminators. In order to use a feature flag as part of your templates you can use the special `env` template function.

## Template Functions

In addition to the `env` template function the templates have access to all of the [sprig](https://masterminds.github.io/sprig/) template functions, as well as `toYaml` and `fromYaml` functions directly from helm.

## Basic Usage

By default using the env function with a key will lookup a feature flag with a given key and interpret any value as a string.

```yaml
my_config_key: [{ env "my_feature_flag" }]
```

There are a few special triggers for the key passed to the `env` function.

## Enabled Syntax

By adding the `_enabled` suffix to a key will lookup a feature flag with the `_enabled` suffix removed and the value interpreted as a boolean. Certain providers have the ability to enable and disable feature flags this state will be used as the value when using this syntax.

```
my_config_key: [{ env "my_feature_flag_enabled" }]
```

## Object Syntax

If the value of your feature flag is a valid json object then you can use the `.` syntax to access subkeys of a json object in your templates.

```
my_config_key: [{ env "my_feature_flag.sub_key" }]
```
