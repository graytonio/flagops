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

## Documentation

- [Getting Started](https://www.graytonward.com/blog/hello-world-for-flagops/)
- [Configuration](/docs/configuration.md)
- [Templates](/docs/templates.md)
- [ArgoCD Plugin](/docs/installation.md)
