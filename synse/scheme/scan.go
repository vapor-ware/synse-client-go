package scheme

// Scan describes a unit in a response of `/scan` endpoint.
type Scan struct {
	ID     string   `json:"id" yaml:"id" mapstructure:"id"`
	Info   string   `json:"info" yaml:"info" mapstructure:"info"`
	Type   string   `json:"type" yaml:"type" mapstructure:"type"`
	Plugin string   `json:"plugin" yaml:"plugin" mapstructure:"plugin"`
	Tags   []string `json:"tags" yaml:"tags" mapstructure:"tags"`
}

// ScanOptions describes the query parameters for `/scan` endpoint.
type ScanOptions struct {
	NS    string   `json:"ns" yaml:"ns" mapstructure:"ns"`
	Tags  []string `json:"tags" yaml:"tags" mapstructure:"tags"`
	Force bool     `json:"force" yaml:"force" mapstructure:"force"`
	Sort  []string `json:"sort" yaml:"sort" mapstructure:"sort"`
}
