package scheme

// Error describes an error response.
type Error struct {
	HTTPCode    int    `json:"http_code" mapstructure:"http_code"`
	Description string `json:"description" mapstructure:"description"`
	Timestamp   string `json:"timestamp" mapstructure:"timestamp"`
	Context     string `json:"context" mapstructure:"context"`
}
