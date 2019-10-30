package scheme

// ReadStreamOptions describe the query parameters for the streamed readings endpoint.
type ReadStreamOptions struct {
	Ids  []string `json:"ids,omitempty" yaml:"ids,omitempty" mapstructure:"ids"`
	Stop bool     `json:"stop,omitempty" yaml:"stop,omitempty" mapstructure:"stop"`
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty" mapstructure:"tags"`
}
