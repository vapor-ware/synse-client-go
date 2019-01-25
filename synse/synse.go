package synse

// synse.go provides a client API for Synse Server.

// Client API for Synse Server.
type Client interface {
	// Test returns the status of Synse Sever. This is used to check if the
	// server is reachable.
	Status() (*Status, error)

	// Version returns the version info for Synse Server.
	Version() (*Version, error)

	// Config returns the config info for Synse Server.
	Config() (*Config, error)
}
