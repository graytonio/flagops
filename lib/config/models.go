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

	// Customer baseURL to send api requests for provider
	BaseURL string `mapstructure:"baseURL"`

	// Configures use of the flagops datastore
	Datastore DatastoreConfig `mapstructure:"datastore"`
}

type DatastoreConfig struct {
	// Enable using flagops datastore
	Enabled bool `mapstructure:"enabled"`

	// Base url of flagops datastore deployment
	BaseURL string `mapstructure:"baseURL"`

	// API Key to authenticate to datastore
	APIKey string `mapstructure:"apiKey"`
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

	// Config for GitRepo Type
	Git GitRepo `mapstructure:"git"`

	// Root path of output. For git type relative to repo root
	Path string `mapstructure:"path"`

	// Do not delete any existing files in destination only template and update specified files
	UpsertMode bool `mapstructure:"upsert"`

	// Insert string at beginning of templated files
	Header string `mapstructure:"header"`

	// Insert string at end of templated files
	Footer string `mapstructure:"footer"`
}

type GitRepo struct {
	// For git type define output git repository
	Repo string `mapstructure:"repo"`

	// For git type define output git branch. Uses default branch if not set.
	Branch string `mapstructure:"branch"`

	// Authentication used for git type. Possible to use for other types later.
	Auth Auth `mapstructure:"auth"`
}

type Auth struct {
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	PrivateKey string `mapstructure:"private_key"`
}

type Config struct {
	Envs  map[string]Environment `mapstructure:"envs"`
	Paths []Path                 `mapstructure:"paths"`
}
