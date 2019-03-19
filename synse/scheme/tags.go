package scheme

// TagsOptions describes the query parameters for `/tags` endpoint.
type TagsOptions struct {
	NS  []string `json:"ns" mapstructure:"ns"`
	IDs bool     `json:"ids" mapstructure:"ids"`
}
