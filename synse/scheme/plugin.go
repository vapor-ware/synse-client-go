package scheme

// Plugin describes a response for `plugin` endpoint when the `id` URI
// parameter is provided.
type Plugin struct {
	PluginMeta `mapstructure:",squash"`
	Network    NetworkOptions `json:"network" mapstructure:"network"`
	Health     HealthOptions  `json:"health" mapstructure:"health"`
}

// PluginMeta holds the plugin metainfo.
type PluginMeta struct {
	Active      bool           `json:"active" mapstructure:"active"`
	ID          string         `json:"id" mapstructure:"id"`
	Name        string         `json:"name" mapstructure:"name"`
	Description string         `json:"description" mapstructure:"description"`
	Maintainer  string         `json:"maintainer" mapstructure:"maintainer"`
	Tag         string         `json:"tag" mapstructure:"tag"`
	VCS         string         `json:"vcs" mapstructure:"vcs"`
	Version     VersionOptions `json:"version" mapstructure:"version"`
}

// VersionOptions holds the version info.
type VersionOptions struct {
	PluginVersion string `json:"plugin_version" mapstructure:"plugin_version"`
	SDKVersion    string `json:"sdk_version" mapstructure:"sdk_version"`
	BuildDate     string `json:"build_date" mapstructure:"build_date"`
	GitCommit     string `json:"git_commit" mapstructure:"git_commit"`
	GitTag        string `json:"git_tag" mapstructure:"git_tag"`
	Arch          string `json:"arch" mapstructure:"arch"`
	OS            string `json:"os" mapstructure:"os"`
}

// NetworkOptions holds communication protocol info.
type NetworkOptions struct {
	Protocol string `json:"protocol" mapstructure:"protocol"`
	Address  string `json:"address" mapstructure:"address"`
}

// HealthOptions holds health info.
type HealthOptions struct {
	Timestamp string         `json:"timestamp" mapstructure:"timestamp"`
	Status    string         `json:"status" mapstructure:"status"`
	Message   string         `json:"message" mapstructure:"message"`
	Checks    []CheckOptions `json:"checks" mapstructure:"checks"`
}

// CheckOptions holds the health check info.
type CheckOptions struct {
	Name      string `json:"name" mapstructure:"name"`
	Status    string `json:"status" mapstructure:"status"`
	Message   string `json:"message" mapstructure:"message"`
	Timestamp string `json:"timestamp" mapstructure:"timestamp"`
	Type      string `json:"type" mapstructure:"type"`
}

// PluginHealth describes a response for `plugin/health` endpoint.
type PluginHealth struct {
	Status    string   `json:"status" mapstructure:"status"`
	Updated   string   `json:"updated" mapstructure:"updated"`
	Healthy   []string `json:"healthy" mapstructure:"healthy"`
	Unhealthy []string `json:"unhealthy" mapstructure:"unhealthy"`
	Active    int      `json:"active" mapstructure:"active"`
	Inactive  int      `json:"inactive" mapstructure:"inactive"`
}
