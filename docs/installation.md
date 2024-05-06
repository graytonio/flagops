# ArgoCD Plugin

Flagops can be natively integrated into ArgoCD as a config plugin. This allows ArgoCD to pass manifests through flagops to natively integrate feature flags into the gitops process.

## Installation

You can install flagops as a sidecar container plugin for argocd using the manifests in the `/manifests` directory.

### CMP Plugin Config

The configuration of the plugin types is controlled via the config maps in `cmp-plugin.yaml`. You can edit the `flagops-env` config map to contain your environment configuration for the plugins.

### Sidecar Apps

The `argocd-repo-server.yaml` file contains a kustomize patch for the argocd-repo-server to add the sidecar applications and configs. This patch should not need to be configured and can be applied as is.

### Applying Plugin

Once the manifests have been configured you can install the plugin with `kustomize build manifests | kubectl apply -`. This will spin up an argocd server with the plugin installed.