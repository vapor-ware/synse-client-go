package scheme

// TagsOptions describes the query parameters for `/tags` endpoint.
type TagsOptions struct {
	NS  []string `json:"ns"`
	IDs bool     `json:"ids"`
}
