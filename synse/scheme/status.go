package scheme

// Status describes a response for `/test` endpoint.
type Status struct {
	Status    string `json:"status" yaml:"status" mapstructure:"status"`
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
}
