package scheme

// Info describes a response from `/info` endpoint.
type Info struct {
	Timestamp    string              `json:"timestamp" mapstructure:"timestamp"`
	ID           string              `json:"id" mapstructure:"id"`
	Type         string              `json:"type" mapstructure:"type"`
	Metadata     MetadataOptions     `json:"metadata" mapstructure:"metadata"`
	Plugin       string              `json:"plugin" mapstructure:"plugin"`
	Info         string              `json:"info" mapstructure:"info"`
	Tags         []string            `json:"tags" mapstructure:"tags"`
	Capabilities CapabilitiesOptions `json:"capabilities" mapstructure:"capabilities"`
	Output       []OutputOptions     `json:"output" mapstructure:"output"`
}

// MetadataOptions holds the metadata info.
type MetadataOptions struct {
	Model string `json:"model" mapstructure:"model"`
}

// CapabilitiesOptions holds the capabilities info.
type CapabilitiesOptions struct {
	Mode  string            `json:"mode" mapstructure:"mode"`
	Read  map[string]string `json:"read" mapstructure:"read"`
	Write WriteOptions      `json:"write" mapstructure:"write"`
}

// WriteOptions holds the write info.
type WriteOptions struct {
	Actions []string `json:"actions" mapstructure:"actions"`
}

// OutputOptions holds the output info.
type OutputOptions struct {
	Name          string        `json:"name" mapstructure:"name"`
	Type          string        `json:"type" mapstructure:"type"`
	Precision     int           `json:"precision" mapstructure:"precision"`
	ScalingFactor float64       `json:"scaling_factor" mapstructure:"scaling_factor"`
	Units         []UnitOptions `json:"units" mapstructure:"units"`
}

// UnitOptions holds the unit info.
type UnitOptions struct {
	System string `json:"system" mapstructure:"system"`
	Name   string `json:"name" mapstructure:"name"`
	Symbol string `json:"symbol" mapstructure:"symbol"`
}
