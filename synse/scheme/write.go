package scheme

// Write describes an unit in a response for the `/write` endpoint.
type Write struct {
	ID      string    `json:"id" yaml:"id" mapstructure:"id"`
	Device  string    `json:"device" yaml:"device" mapstructure:"device"`
	Context WriteData `json:"context" yaml:"context" mapstructure:"context"`
	Timeout string    `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

// WriteData describes an unit in the POST body for the `/write` endpoint. This
// can also be used in a websocket request event payload.
type WriteData struct {
	Transaction string `json:"transaction,omitempty" yaml:"transaction,omitempty" mapstructure:"transaction"`
	Action      string `json:"action" yaml:"action" mapstructure:"action"`

	// data is always string, as the conversion happens on the plugin side.
	Data string `json:"data,omitempty" yaml:"data,omitempty" mapstructure:"data"`
}
