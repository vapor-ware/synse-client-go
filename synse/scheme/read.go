package scheme

// Read describes a unit in a response for `/read` endpoint.
type Read struct {
	Device     string                 `json:"device" yaml:"device" mapstructure:"device"`
	DeviceType string                 `json:"device_type" yaml:"device_type" mapstructure:"device_type"`
	Timestamp  string                 `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
	Unit       UnitOptions            `json:"unit" yaml:"unit" mapstructure:"unit"`
	Context    map[string]interface{} `json:"context" yaml:"context" mapstructure:"context"`
}

// ReadOptions describes the query parameters for `/read` endpoint.
type ReadOptions struct {
	NS   string   `json:"ns" yaml:"ns"`
	Tags []string `json:"tags" yaml:"tags"`
}
