# Configuration

## Environments

Environments are the different feature flag sources that flagops can pull from. These could be different providers for different projects or different environments such as production and staging. Each entry is named based on the dictionary key and is then configured with the provider and any authentication parameters based on each provider.

| Key      | Description                               |
| -------- | ----------------------------------------- |
| provider | Flag Provider to Configure                |
| apiKey   | API Key to authenticate with the provider |

### Supported Providers

- `flagsmith`: [Flagsmith](https://www.flagsmith.com/)
- `launchdarkly`: [LaunchDarkly](https://launchdarkly.com/) **BETA**
- `env`: [FromEnv](https://github.com/open-feature/go-sdk-contrib/tree/main/providers/from-env)

## Paths

Paths are directories or individual files that should be templated and placed in a location. Individual files are templated and then sent to the configured destination path. Directories are templated recursively and preserve their directory structure when sent to the configured destination path.

| Key         | Description                                                                                                                                                                                                        |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| path        | Path relative to the `.flagops` file to source the templates from. Can be a directory or a file                                                                                                                    |
| env         | Which configured environment to use for flag evaluation                                                                                                                                                            |
| identity    | An identifier for this path execution. Used when evaluating flag values                                                                                                                                            |
| properties  | A dictionary of key value pairs to add context when evaluating flag values                                                                                                                                         |
| dest        | Configure destination for output templates                                                                                                                                                                         |
| dest.type   | Configure destination type                                                                                                                                                                                         |
| dest.path   | Path relative to destination root to where to output the templates to. Any missing directories will automatically be created.                                                                                      |
| dest.git    | For git destinations the repo to push the changes to                                                                                                                                                               |
| dest.upsert | By default the path in the destination is fully cleaned before inserting the templated files. If you would like to keep existing files in the destination that do not exist on the source set the upsert to false. |

### Supported Destination Types

- `git`: Git Repo
- `file`: Local Path
- `console`: Console

## Dynamic Paths

FlagOps can also be configured to pull it's path configuration from the environment instead of a config file. There are a few restrictions with this mode of operation:

1. A config file is still needed to define the environments
2. Only one path can be configured per run

The config for a path are mapped to env variables at runtime

| Config Key  | Env                      |
| ----------- | ------------------------ |
| path        | FLAGOPS_SOURCE_PATH      |
| env         | FLAGOPS_ENVIRONMENT      |
| identity    | FLAGOPS_IDENTITY         |
| dest.type   | FLAGOPS_DESTINATION_TYPE |
| dest.path   | FLAGOPS_DESTINATION_PATH |
| dest.git    | FLAGOPS_DESTINATION_REPO |
| dest.upsert | FLAGOPS_UPSERT           |

`properties` are defined by the prefix `FLAGOPS_PROP_` so for example the env `FLAGOPS_PROP_REGION` with a value of `aws-usw-2` would be parsed to a property of `region` with a value of `aws-usw-2`
