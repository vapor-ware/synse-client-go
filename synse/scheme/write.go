package scheme

// Write describes an unit in a response for the `/write` endpoint.
type Write struct {
	Context     WriteData `json:"context"`
	Device      string    `json:"device"`
	Transaction string    `json:"transaction"`
	Timeout     string    `json:"timeout"`
}

// WriteData describes an unit in the POST body for the `/write` endpoint and
// request/write event.
type WriteData struct {
	Device      string      `json:"device"`
	Transaction string      `json:"transaction"`
	Action      string      `json:"action"`
	Data        interface{} `json:"data"`
}
