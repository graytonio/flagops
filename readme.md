# FlagOps

A cli designed to help use feature flag services in git-ops deployment structures.

## Usage

The main command of the cli is the generate command which takes input templates and applies the feature flags from your provider to the template.

```bash
$ flagops generate --help
Run template engine over files

Usage:
  kustomize-flag generate [flags]

Flags:
  -d, --dest string          Write template to file instead of stdout. If recursive is enabled the output base directory (default: build)
  -f, --filter stringArray   Apply a filter over the produced template
  -h, --help                 help for generate
  -p, --provider string      Provider type to pull context from (default "flagsmith")
  -r, --recursive            Enable recursive mode to apply templating to entire directory
  -s, --src string           Input file to generate or directory to start from is using recursive

Global Flags:
      --config string   Path to configuration file (default "./.kustomgen.yaml")
      --verbose         Enable verbose logging
```

**Example**

```
$ cat my-service-template.yaml.j2
apiVersion: v1
kind: Service
metadata:
  name: hello-world
spec:
  selector:
    app: hello-world
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
{% if service_lb_enabled %}
    type: LoadBalancer
{% endif %}
{% if service_node_port_enabled %}
      nodePort: {{ service_node_port }}
  type: NodePort
{% endif %}


$ flagops generate --provider flagsmith -s my-deployment-template.yaml.j2
apiVersion: v1
kind: Service
metadata:
  name: hello-world
spec:
  selector:
    app: hello-world
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
    type: LoadBalancer
```

### Configurations

Configurations can be pulled either from either the cli arguments or from a config file. By default it will look for a file in the local directory named `.flagops.yaml`

**Example**

```yaml
src: my-kustomize-templates
recursive: true
filters:
  - blank
provider: flagsmith
providerConfig:
  apiKey: my-api-key
  baseUrl: http://localhost:8000/api/v1/
```

## Feature Flag Providers

FlagOps has the ability to integrate with many different Feature Flag providers however currently it is only integrated Flagsmith.

If there is a provider you would like to see integrated please open an issue about it.

## Filters

FlagOps also provides post processing filters for you templates. They are documented in this table

| Filter Name       | ID    | Description                                                    |
| ----------------- | ----- | -------------------------------------------------------------- |
| Blank Line Filter | blank | Removes any blank lines from the resulting template execution. |
