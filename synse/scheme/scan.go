package scheme

// Scan describes a unit in a response of `/scan` endpoint.
type Scan struct {
	ID     string   `json:"id"`
	Info   string   `json:"info"`
	Type   string   `json:"type"`
	Plugin string   `json:"plugin"`
	Tags   []string `json:"tags"`
}

// ScanOptions holds the scan query parameters.
type ScanOptions struct {
	NS    string
	Tags  []string
	Force bool
	Sort  []string
}
