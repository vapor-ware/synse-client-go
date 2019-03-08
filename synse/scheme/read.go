package scheme

// Read describes a unit in a response for `/read` endpoint.
type Read struct {
	Device     string                 `json:"device"`
	DeviceType string                 `json:"device_type"`
	Type       string                 `json:"type"`
	Value      interface{}            `json:"value"`
	Timestamp  string                 `json:"timestamp"`
	Unit       UnitOptions            `json:"unit"`
	Context    map[string]interface{} `json:"context"`
}

// ReadOptions describes the query parameters for `/read` endpoint.
type ReadOptions struct {
	NS   string   `json:"ns"`
	Tags []string `json:"tags"`
	SOM  string   `json:"som"`
}
