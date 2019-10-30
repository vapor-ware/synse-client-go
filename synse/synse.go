package synse

// synse.go provides a client API for Synse Server.

import (
	"github.com/vapor-ware/synse-client-go/synse/scheme"
)

// Client API for Synse Server.
type Client interface {
	// Status returns the status info. This is used to check if the server
	// is responsive and reachable.
	Status() (*scheme.Status, error)

	// Version returns the version info.
	Version() (*scheme.Version, error)

	// Config returns the unified configuration info.
	Config() (*scheme.Config, error)

	// Plugins returns the summary of all plugins currently registered with
	// Synse Server.
	Plugins() ([]*scheme.PluginMeta, error)

	// Plugin returns data from a specific plugin.
	Plugin(string) (*scheme.Plugin, error)

	// PluginHealth returns the summary of the health of registered plugins.
	PluginHealth() (*scheme.PluginHealth, error)

	// Scan returns the list of devices that Synse knows about and can read
	// from/write to via the configured plugins.
	// It can be filtered to show only those devices which match a set
	// of provided tags by using ScanOptions.
	Scan(scheme.ScanOptions) ([]*scheme.Scan, error)

	// Tags returns the list of all tags currently associated with devices.
	// If no TagsOptions is specified, the default tag namespace will be used.
	Tags(scheme.TagsOptions) ([]string, error)

	// Info returns the full set of meta info and capabilities for a specific
	// device.
	Info(string) (*scheme.Info, error)

	// Read returns data from devices which match the set of provided tags
	// using ReadOptions.
	Read(scheme.ReadOptions) ([]*scheme.Read, error)

	// ReadDevice returns data from a specific device.
	// It is the same as Read() where the label matches the device id tag
	// specified in ReadOptions.
	ReadDevice(string) ([]*scheme.Read, error)

	// ReadCache returns cached reading data from the registered plugins.
	ReadCache(scheme.ReadCacheOptions, chan<- *scheme.Read) error

	// ReadStream returns a stream of current reading data from the
	// registered plugins.
	ReadStream(scheme.ReadStreamOptions, chan<- *scheme.Read, chan struct{}) error

	// WriteAsync writes data to a device, in an asynchronous manner.
	WriteAsync(string, []scheme.WriteData) ([]*scheme.Write, error)

	// WriteSync writes data to a device, waiting for the write to complete.
	WriteSync(string, []scheme.WriteData) ([]*scheme.Transaction, error)

	// Transactions returns the sorted list of all cached transaction IDs.
	Transactions() ([]string, error)

	// Transaction returns the state and status of a write transaction.
	Transaction(string) (*scheme.Transaction, error)

	// GetOptions returns the current config options of the client.
	GetOptions() *Options

	// Open opens the websocket connection between the client and Synse Server.
	// As the description suggested, it is only applicable for a WebSocket
	// client. Calling this method on a HTTP Client will have no effect.
	Open() error

	// Close closes the websocket connection between the client and Synse Server.
	// It is only applicable for a WebSocket client in a sense that, one must
	// close the connection after finish using it. Calling this method on a
	// HTTP Client will have no effect.
	Close() error
}
