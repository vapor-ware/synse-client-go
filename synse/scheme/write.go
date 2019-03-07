package scheme

// Write describes an unit in a response for the `/write` endpoint.
type Write struct {
	Context     WriteData `json:"context"`
	Device      string    `json:"device"`
	Transaction string    `json:"transaction"`
	Timeout     string    `json:"timeout"`
}

// WriteData describes an unit in the POST body for the `/write` endpoint. This
// can also be used with in a websocket request event.
type WriteData struct {
	Plugin      string      `json:"plugin"`
	Device      string      `json:"device"`
	Transaction string      `json:"transaction"`
	Action      string      `json:"action"`
	Data        interface{} `json:"data"`
}
