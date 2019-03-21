package scheme

// Write describes an unit in a response for the `/write` endpoint.
type Write struct {
	Context     WriteData `json:"context" mapstructure:"context"`
	Device      string    `json:"device" mapstructure:"device"`
	Transaction string    `json:"transaction" mapstructure:"transaction"`
	Timeout     string    `json:"timeout" mapstructure:"timeout"`
}

// WriteData describes an unit in the POST body for the `/write` endpoint. This
// can also be used in a websocket request event payload.
type WriteData struct {
	Transaction string      `json:"transaction"`
	Action      string      `json:"action"`
	Data        interface{} `json:"data"`
}
