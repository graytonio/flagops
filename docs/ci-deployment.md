# CI Pipeline

Flagops can be deployed as a CI pipeline that generate the GitOps manifests watched by your other tools like ArgoCD. This allows you to keep the benefits of GitOps while offering an easier configuration tool to external teams however it can require more manual configuration to maintain.

## Implementation

The recommended setup for this deployment is to have your main branch contain the template files and have FlagOps configured to deploy each environment into a new branch.

For example:

```yaml
envs:
    production:
        provider: flagsmith
        apiKey: aaaabbbbccccdddd
    qa:
        provider: flagsmith
        apiKey: aaaabbbbccccdddd
    dev:
        provider: flagsmith
        apiKey: aaaabbbbccccdddd

paths:
    - path: apps/
      env: production
      dest:
        type: git
        repo: git@github.com:graytonio/my-gitops-repo.git
        path: apps/
        branch: production
    - path: apps/
      env: qa
      dest:
        type: git
        repo: git@github.com:graytonio/my-gitops-repo.git
        path: apps/
        branch: qa
    - path: apps/
      env: dev
      dest:
        type: git
        repo: git@github.com:graytonio/my-gitops-repo.git
        path: apps/
        branch: dev
```

Now as an example in ArgoCD you can create 3 appsets that resolve each of these branches to the production, qa, and dev clusters where each one will be updated as their environment is updated.

This allows you to keep the atomic commits of git with easy rollback to a previous hash if something goes wrong but maintaining the simple interface for users to configure their apps.