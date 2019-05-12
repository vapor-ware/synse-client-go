package scheme

// Write describes an unit in a response for the `/write` endpoint.
type Write struct {
	Context     WriteData `json:"context" yaml:"context" mapstructure:"context"`
	Device      string    `json:"device" yaml:"device" mapstructure:"device"`
	Transaction string    `json:"transaction" yaml:"transaction" mapstructure:"transaction"`
	Timeout     string    `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

// WriteData describes an unit in the POST body for the `/write` endpoint. This
// can also be used in a websocket request event payload.
type WriteData struct {
	Transaction string      `json:"transaction" yaml:"transaction" mapstructure:"transaction"`
	Action      string      `json:"action" yaml:"action" mapstructure:"action"`
	Data        interface{} `json:"data" yaml:"data" mapstructure:"data"`
}
