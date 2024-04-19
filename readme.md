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
      upsert: false
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

## Configuration

### Environments

Environments are the different feature flag sources that flagops can pull from. These could be different providers for different projects or different environments such as production and staging. Each entry is named based on the dictionary key and is then configured with the provider and any authentication parameters based on each provider.

| Key      | Description                               |
| -------- | ----------------------------------------- |
| provider | Flag Provider to Configure                |
| apiKey   | API Key to authenticate with the provider |

#### Supported Providers

- `flagsmith`: Flagsmith

### Paths

Paths are directories or individual files that should be templated and placed in a location. Individual files are templated and then sent to the configured destination path. Directories are templated recursively and preserve their directory structure when sent to the configured destination path.

| Key  | Description                                                                                     |
| ---- | ----------------------------------------------------------------------------------------------- |
| path | Path relative to the `.flagops` file to source the templates from. Can be a directory or a file |
| env  | Which configured environment to use for flag evaluation                                         |

| dest | Configure destination for output templates |
| dest.type | Configure destination type |
| dest.path | Path relative to destination root to where to output the templates to. Any missing directories will automatically be created. |
| dest.git | For git destinations the repo to push the changes to |
| dest.upsert | By default the path in the destination is fully cleaned before inserting the templated files. If you would like to keep existing files in the destination that do not exist on the source set the upsert to false. |

#### Supported Destination Types

- `git`: Git Repo
- `file`: Local Path
- `console`: Console

## Template Files

Flagops uses the golang templating syntax however it uses `[{` and `}]` as deliminators. In order to use a feature flag as part of your templates you can use the special `env` template function.

### Template Functions

In addition to the `env` template function the templates have access to all of the [sprig](https://masterminds.github.io/sprig/) template functions, as well as `toYaml` and `fromYaml` functions directly from helm.

### Basic Usage

By default using the env function with a key will lookup a feature flag with a given key and interpret any value as a string.

```yaml
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
