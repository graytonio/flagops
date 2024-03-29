package config

type ProviderType string

const (
	Flagsmith ProviderType = "flagsmith"
)

type Environment struct {
	// Provider to configure for this environment
	Provider ProviderType `mapstructure:"provider"`
	
	// Api key to pass to provider
	APIKey string `mapstructure:"apiKey"`

	// Environment variable to source the api key from instead of apiKey
	EnvKey string `mapstructure:"envKey"`
}

type Path struct {
	// File glob pattern to template
	Path string `mapstructure:"path"`

	// Environment to use for feature flag population
	Env string `mapstructure:"env"`

	// Output definition
	Destination Destination `mapstructure:"dest"`
}

type Destination struct {
	// Output type (file, git)
	Type string `mapstructure:"type"`

	// For git type define output git repository
	Repo string `mapstructure:"repo"`

	// Root path of output. For git type relative to repo root
	Path string `mapstructure:"path"`
}

type Config struct {
	Envs map[string]Environment `mapstructure:"envs"`
	Paths []Path `mapstructure:"paths"`
}