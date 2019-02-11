package scheme

// Transaction describes a response for `/transaction` endpoint. It also
// describes an unit in a response for `/write/wait` endpoint.
type Transaction struct {
	ID      string    `json:"id"`
	Timeout string    `json:"timeout"`
	Device  string    `json:"device"`
	Context WriteData `json:"context"`
	Status  string    `json:"status"`
	Created string    `json:"created"`
	Updated string    `json:"updated"`
	Message string    `json:"message"`
}
