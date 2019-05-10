package scheme

// Transaction describes a response for `/transaction` endpoint. It also
// describes an unit in a response for `/write/wait` endpoint.
type Transaction struct {
	ID      string    `json:"id" yaml:"id" mapstructure:"id"`
	Timeout string    `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Device  string    `json:"device" yaml:"device" mapstructure:"device"`
	Context WriteData `json:"context" yaml:"context" mapstructure:"context"`
	Status  string    `json:"status" yaml:"status" mapstructure:"status"`
	Created string    `json:"created" yaml:"created" mapstructure:"created"`
	Updated string    `json:"updated" yaml:"updated" mapstructure:"updated"`
	Message string    `json:"message" yaml:"message" mapstructure:"message"`
}
