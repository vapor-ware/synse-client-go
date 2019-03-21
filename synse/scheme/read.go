package scheme

// Read describes a unit in a response for `/read` endpoint.
type Read struct {
	Device     string                 `json:"device" mapstructure:"device"`
	DeviceType string                 `json:"device_type" mapstructure:"device_type"`
	Type       string                 `json:"type" mapstructure:"type"`
	Value      interface{}            `json:"value" mapstructure:"value"`
	Timestamp  string                 `json:"timestamp" mapstructure:"timestamp"`
	Unit       UnitOptions            `json:"unit" mapstructure:"unit"`
	Context    map[string]interface{} `json:"context" mapstructure:"context"`
}

// ReadOptions describes the query parameters for `/read` endpoint.
type ReadOptions struct {
	NS   string   `json:"ns" mapstructure:"ns"`
	Tags []string `json:"tags" mapstructure:"tags"`
}
