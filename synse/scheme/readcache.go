package scheme

// ReadCacheOptions describes the query parameters for `/readcache` endpoint.
type ReadCacheOptions struct {
	Start string `json:"start" yaml:"start" mapstructure:"start"`
	End   string `json:"end" yaml:"end" mapstructure:"end"`
}
