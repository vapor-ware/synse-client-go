package scheme

// Plugin describes a response for `plugin` endpoint when the `id` URI
// parameter is provided.
type Plugin struct {
	PluginMeta `mapstructure:",squash"`
	Network    NetworkOptions `json:"network" yaml:"network" mapstructure:"network"`
	Health     HealthOptions  `json:"health" yaml:"health" mapstructure:"health"`
}

// PluginMeta holds the plugin metainfo.
type PluginMeta struct {
	Active      bool           `json:"active" yaml:"active" mapstructure:"active"`
	ID          string         `json:"id" yaml:"id" mapstructure:"id"`
	Name        string         `json:"name" yaml:"name" mapstructure:"name"`
	Description string         `json:"description" yaml:"description" mapstructure:"description"`
	Maintainer  string         `json:"maintainer" yaml:"maintainer" mapstructure:"maintainer"`
	Tag         string         `json:"tag" yaml:"tag" mapstructure:"tag"`
	VCS         string         `json:"vcs" yaml:"vcs" mapstructure:"vcs"`
	Version     VersionOptions `json:"version" yaml:"version" mapstructure:"version"`
}

// VersionOptions holds the version info.
type VersionOptions struct {
	PluginVersion string `json:"plugin_version" yaml:"plugin_version" mapstructure:"plugin_version"`
	SDKVersion    string `json:"sdk_version" yaml:"sdk_version" mapstructure:"sdk_version"`
	BuildDate     string `json:"build_date" yaml:"build_date" mapstructure:"build_date"`
	GitCommit     string `json:"git_commit" yaml:"git_commit" mapstructure:"git_commit"`
	GitTag        string `json:"git_tag" yaml:"git_tag" mapstructure:"git_tag"`
	Arch          string `json:"arch" yaml:"arch" mapstructure:"arch"`
	OS            string `json:"os" yaml:"os" mapstructure:"os"`
}

// NetworkOptions holds communication protocol info.
type NetworkOptions struct {
	Protocol string `json:"protocol" yaml:"protocol" mapstructure:"protocol"`
	Address  string `json:"address" yaml:"address" mapstructure:"address"`
}

// HealthOptions holds health info.
type HealthOptions struct {
	Timestamp string         `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
	Status    string         `json:"status" yaml:"status" mapstructure:"status"`
	Message   string         `json:"message" yaml:"message" mapstructure:"message"`
	Checks    []CheckOptions `json:"checks" yaml:"checks" mapstructure:"checks"`
}

// CheckOptions holds the health check info.
type CheckOptions struct {
	Name      string `json:"name" yaml:"name" mapstructure:"name"`
	Status    string `json:"status" yaml:"status" mapstructure:"status"`
	Message   string `json:"message" yaml:"message" mapstructure:"message"`
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
	Type      string `json:"type" yaml:"type" mapstructure:"type"`
}

// PluginHealth describes a response for `plugin/health` endpoint.
type PluginHealth struct {
	Status    string   `json:"status" yaml:"status" mapstructure:"status"`
	Updated   string   `json:"updated" yaml:"updated" mapstructure:"updated"`
	Healthy   []string `json:"healthy" yaml:"healthy" mapstructure:"healthy"`
	Unhealthy []string `json:"unhealthy" yaml:"unhealthy" mapstructure:"unhealthy"`
	Active    int      `json:"active" yaml:"active" mapstructure:"active"`
	Inactive  int      `json:"inactive" yaml:"inactive" mapstructure:"inactive"`
}
