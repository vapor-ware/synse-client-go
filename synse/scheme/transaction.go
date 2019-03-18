package scheme

// Transaction describes a response for `/transaction` endpoint. It also
// describes an unit in a response for `/write/wait` endpoint.
type Transaction struct {
	ID      string    `json:"id" mapstructure:"id"`
	Timeout string    `json:"timeout" mapstructure:"timeout"`
	Device  string    `json:"device" mapstructure:"device"`
	Context WriteData `json:"context" mapstructure:"context"`
	Status  string    `json:"status" mapstructure:"status"`
	Created string    `json:"created" mapstructure:"created"`
	Updated string    `json:"updated" mapstructure:"updated"`
	Message string    `json:"message" mapstructure:"message"`
}
