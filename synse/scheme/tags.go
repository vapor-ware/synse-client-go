package scheme

// TagsOptions describes the query parameters for `/tags` endpoint.
type TagsOptions struct {
	NS  []string `json:"ns" yaml:"ns" mapstructure:"ns"`
	IDs bool     `json:"ids" yaml:"ids" mapstructure:"ids"`
}
