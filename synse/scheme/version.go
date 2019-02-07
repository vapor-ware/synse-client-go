package scheme

// Version describes a response for `version` endpoint.
type Version struct {
	Version    string `json:"version"`
	APIVersion string `json:"api_version"`
}
