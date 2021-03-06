package scheme

// Scan describes a unit in a response of `/scan` endpoint.
type Scan struct {
	ID       string                 `json:"id" yaml:"id" mapstructure:"id"`
	Alias    string                 `json:"alias" yaml:"alias" mapstructure:"alias"`
	Info     string                 `json:"info" yaml:"info" mapstructure:"info"`
	Type     string                 `json:"type" yaml:"type" mapstructure:"type"`
	Plugin   string                 `json:"plugin" yaml:"plugin" mapstructure:"plugin"`
	Tags     []string               `json:"tags" yaml:"tags" mapstructure:"tags"`
	Metadata map[string]interface{} `json:"metadata" yaml:"metadata" mapstructure:"metadata"`
}

// ScanOptions describes the query parameters for `/scan` endpoint.
type ScanOptions struct {
	NS    string   `json:"ns,omitempty" yaml:"ns,omitempty" mapstructure:"ns"`
	Tags  []string `json:"tags,omitempty" yaml:"tags,omitempty" mapstructure:"tags"`
	Force bool     `json:"force,omitempty" yaml:"force,omitempty" mapstructure:"force"`
	Sort  []string `json:"sort,omitempty" yaml:"sort,omitempty" mapstructure:"sort"`
}
