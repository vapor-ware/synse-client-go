package scheme

// Plugin describes a response for `plugin` endpoint when the `id` URI
// parameter is provided.
type Plugin struct {
	PluginMeta
	Network NetworkOptions `json:"network"`
	Health  HealthOptions  `json:"health"`
}

// PluginMeta holds the plugin metainfo.
type PluginMeta struct {
	Active      bool           `json:"active"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Maintainer  string         `json:"maintainer"`
	Tag         string         `json:"tag"`
	VCS         string         `json:"vcs"`
	Version     VersionOptions `json:"version"`
}

// VersionOptions holds the version info.
type VersionOptions struct {
	PluginVersion string `json:"plugin_version"`
	SDKVersion    string `json:"sdk_version"`
	BuildDate     string `json:"build_date"`
	GitCommit     string `json:"git_commit"`
	GitTag        string `json:"git_tag"`
	Arch          string `json:"arch"`
	OS            string `json:"os"`
}

// NetworkOptions holds communication protocol info.
type NetworkOptions struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
}

// HealthOptions holds health info.
type HealthOptions struct {
	Timestamp string         `json:"timestamp"`
	Status    string         `json:"status"`
	Message   string         `json:"message"`
	Checks    []CheckOptions `json:"checks"`
}

// CheckOptions holds the health check info.
type CheckOptions struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}
