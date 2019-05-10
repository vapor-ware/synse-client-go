package scheme

// Version describes a response for `version` endpoint.
type Version struct {
	Version    string `json:"version" yaml:"version" mapstructure:"version"`
	APIVersion string `json:"api_version" yaml:"version" mapstructure:"api_version"`
}
