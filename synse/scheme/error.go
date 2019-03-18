package scheme

// Error describes an error response.
type Error struct {
	HTTPCode    int    `json:"http_code" mapstructure:"http_code"`
	ErrorID     int    `json:"error_id" mapstructure:"error_id"`
	Description string `json:"description" mapstructure:"description"`
	Timestamp   string `json:"timestamp" mapstructure:"timestamp"`
	Context     string `json:"context" mapstructure:"context"`
}
