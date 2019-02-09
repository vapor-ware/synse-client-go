package scheme

// Info describes a response from `/info` endpoint.
type Info struct {
	Timestamp    string              `json:"timestamp"`
	ID           string              `json:"id"`
	Type         string              `json:"type"`
	Metadata     MetadataOptions     `json:"metadata"`
	Plugin       string              `json:"plugin"`
	Info         string              `json:"info"`
	Tags         []string            `json:"tags"`
	Capabilities CapabilitiesOptions `json:"capabilities"`
	Output       []OutputOptions     `json:"output"`
}

// MetadataOptions holds the metadata info.
type MetadataOptions struct {
	Model string `json:"model"`
}

// CapabilitiesOptions holds the capabilities info.
type CapabilitiesOptions struct {
	Mode  string            `json:"mode"`
	Read  map[string]string `json:"read"`
	Write WriteOptions      `json:"write"`
}
type WriteOptions struct {
	Actions []string `json:"actions"`
}

// OutputOptions holds the output info.
type OutputOptions struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Precision     int           `json:"precision"`
	ScalingFactor float64       `json:"scaling_factor"`
	Units         []UnitOptions `json:"units"`
}

// UnitOptions holds the unit info.
type UnitOptions struct {
	System string `json:"system"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
