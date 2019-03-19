package scheme

// Scan describes a unit in a response of `/scan` endpoint.
type Scan struct {
	ID     string   `json:"id" mapstructure:"id"`
	Info   string   `json:"info" mapstructure:"info"`
	Type   string   `json:"type" mapstructure:"type"`
	Plugin string   `json:"plugin" mapstructure:"plugin"`
	Tags   []string `json:"tags" mapstructure:"tags"`
}

// ScanOptions describes the query parameters for `/scan` endpoint.
type ScanOptions struct {
	NS    string
	Tags  []string
	Force bool
	Sort  []string
}
