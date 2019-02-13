package scheme

// Scan describes a unit in a response of `/scan` endpoint.
type Scan struct {
	ID     string   `json:"id"`
	Info   string   `json:"info"`
	Type   string   `json:"type"`
	Plugin string   `json:"plugin"`
	Tags   []string `json:"tags"`
}

// ScanOptions describes the query parameters for `/scan` endpoint.
type ScanOptions struct {
	NS    string
	Tags  []string
	Force bool
	Sort  []string
}
