package scheme

// Error describes an error response.
type Error struct {
	HTTPCode    int    `json:"http_code"`
	ErrorID     int    `json:"error_id"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	Context     string `json:"context"`
}
