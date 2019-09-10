package scheme

// TagsOptions describes the query parameters for `/tags` endpoint.
type TagsOptions struct {
	NS  []string `json:"ns,omitempty" yaml:"ns,omitempty" mapstructure:"ns"`
	IDs bool     `json:"ids,omitempty" yaml:"ids,omitempty" mapstructure:"ids"`
}
