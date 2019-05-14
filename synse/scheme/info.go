package scheme

// Info describes a response from `/info` endpoint.
type Info struct {
	Timestamp    string              `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
	ID           string              `json:"id" yaml:"id" mapstructure:"id"`
	Type         string              `json:"type" yaml:"type" mapstructure:"type"`
	Metadata     map[string]string   `json:"metadata" yaml:"metadata" mapstructure:"metadata"`
	Plugin       string              `json:"plugin" yaml:"plugin" mapstructure:"plugin"`
	Info         string              `json:"info" yaml:"info" mapstructure:"info"`
	Tags         []string            `json:"tags" yaml:"tags" mapstructure:"tags"`
	Capabilities CapabilitiesOptions `json:"capabilities" yaml:"capabilities" mapstructure:"capabilities"`
	Output       []OutputOptions     `json:"output" yaml:"output" mapstructure:"output"`
}

// CapabilitiesOptions holds the capabilities info.
type CapabilitiesOptions struct {
	Mode  string            `json:"mode" yaml:"mode" mapstructure:"mode"`
	Read  map[string]string `json:"read" yaml:"read" mapstructure:"read"`
	Write WriteOptions      `json:"write" yaml:"write" mapstructure:"write"`
}

// WriteOptions holds the write info.
type WriteOptions struct {
	Actions []string `json:"actions" yaml:"actions" mapstructure:"actions"`
}

// OutputOptions holds the output info.
type OutputOptions struct {
	Name          string        `json:"name" yaml:"name" mapstructure:"name"`
	Type          string        `json:"type" yaml:"type" mapstructure:"type"`
	Precision     int           `json:"precision" yaml:"precision" mapstructure:"precision"`
	ScalingFactor float64       `json:"scaling_factor" yaml:"scaling_factor" mapstructure:"scaling_factor"`
	Units         []UnitOptions `json:"units" yaml:"units" mapstructure:"units"`
}

// UnitOptions holds the unit info.
type UnitOptions struct {
	Name   string `json:"name" yaml:"name" mapstructure:"name"`
	Symbol string `json:"symbol" yaml:"symbol" mapstructure:"symbol"`
}
