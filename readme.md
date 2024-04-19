# FlagOps

A DevOps tool for generating git ops repositories by combining application templates and environment flag providers.

## Install

```bash
go install github.com/graytonio/flagops@latest
```

## Quick Start

First create a config file where your templates will be located called `.flagops`.

```yaml
envs:
    production:
        provider: flagsmith
        apiKey: prod-key
    staging:
        provider: flagsmith
        apiKey: stg-key
paths:
    - path: apps/
      env: production
      dest:
        type: git
        repo: git@github.com:user/gitops-deployment.git
        path: apps/production
    - path: apps/
      env: staging
      dest:
        type: git
        repo: git@github.com:user/gitops-deployment.git
        path: apps/production
```

The config file is broken into two parts:

1. Environments - Defines a flag provider that can be used to fill out the templates.
2. Paths - Directories that should be run through the templater with the given environment and output to a specific directory.

Once this file is present you can run the tool:

```bash
flagops
```

This will read the config file from the current directory and output the templated files in the requested output.

## Template Files

Flagops usese the golang templating syntax however it uses `[{` and `}]` as deliminators. In order to use a feature flag as part of your templates you can use the special `env` template function.

### Template Functions

In addition to the `env` tempalte function the templates have access to all of the [sprig](https://masterminds.github.io/sprig/) template functions, as well as `toYaml` and `fromYaml` functions directly from helm.

### Basic Usage

By default using the env function with a key will lookup a feature flag with a given key and interpret any value as a string.

```
my_config_key: [{ env "my_feature_flag" }]
```

There are a few special triggers for the key passed to the `env` function.

### Enabled Syntax

By adding the `_enabled` suffix to a key will lookup a feature flag with the `_enabled` suffix removed and the value interpreted as a boolean. Certain providers have the ability to enable and disable feature flags this state will be used as the value when using this syntax.

```
my_config_key: [{ env "my_feature_flag_enabled" }]
```

### Object Syntax

If the value of your feature flag is a valid json object then you can use the `.` syntax to access subkeys of a json object in your templates.

```
my_config_key: [{ env "my_feature_flag.sub_key" }]
```

## Supported Providers

This is the list of currently supported providers:

- Flagsmith

**NOTE:** The internals of the engine use the open feature sdk as it's base so it will be very easy to add any open feature provider later down the line.

## Supported Outputs

- Git Repo
- Local Path
- Console