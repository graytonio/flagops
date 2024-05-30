package config

type ProviderType string

const (
	Flagsmith    ProviderType = "flagsmith"
	LaunchDarkly ProviderType = "launchdarkly"
	FromEnv      ProviderType = "env"
)

type DestinationType string

const (
	Git     DestinationType = "git"
	File    DestinationType = "file"
	Console DestinationType = "console"
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

	// Unique name for this set of templates
	Identity string `mapstructure:"identity"`

	// Key Value pairs to add context to template engine
	Properties map[string]any `mapstructure:"properties"`
}

type Destination struct {
	// Output type (file, git)
	Type DestinationType `mapstructure:"type"`

	// For git type define output git repository
	Repo string `mapstructure:"repo"`

	// Root path of output. For git type relative to repo root
	Path string `mapstructure:"path"`

	// Do not delete any existing files in destination only template and update specified files
	UpsertMode bool `mapstructure:"upsert"`

	// Insert string at beginning of templated files
	Header string `mapstructure:"header"`

	// Insert string at end of templated files
	Footer string `mapstructure:"footer"`
}

type Config struct {
	Envs  map[string]Environment `mapstructure:"envs"`
	Paths []Path                 `mapstructure:"paths"`
}
