package scheme

// Version describes a response for `version` endpoint.
type Version struct {
	Version    string `json:"version" mapstructure:"version"`
	APIVersion string `json:"api_version" mapstructure:"api_version"`
}
