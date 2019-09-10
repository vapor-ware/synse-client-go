package scheme

// ReadCacheOptions describes the query parameters for `/readcache` endpoint.
type ReadCacheOptions struct {
	Start string `json:"start,omitempty" yaml:"start,omitempty" mapstructure:"start"`
	End   string `json:"end,omitempty" yaml:"end,omitempty" mapstructure:"end"`
}
