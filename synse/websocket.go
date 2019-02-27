package synse

// websocket.go implements a websocket client.

import (
	"crypto/tls"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
)

type websocketClient struct {
	// options is the global config options of the client.
	options *Options

	// client holds the websocket.Dialer.
	client *websocket.Dialer

	// apiVersion is the current api version of Synse Server that we are
	// communicating with.
	apiVersion string
}

// NewWebSocketClientV3 returns a new instance of a websocket client for v3.
func NewWebSocketClientV3(opts *Options) (Client, error) {
	c, err := createWebSocketClient(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a websocket client")
	}

	return &websocketClient{
		options:    opts,
		client:     c,
		apiVersion: "v3",
	}, nil
}

// createWebSocketClient setups a websocket dialer with configured options.
func createWebSocketClient(opts *Options) (*websocket.Dialer, error) {
	err := setDefaults(opts)
	if err != nil {
		return nil, err
	}

	if opts.TLS.Enabled == false {
		return &websocket.Dialer{
			HandshakeTimeout: opts.WebSocket.HandshakeTimeout,
		}, nil
	}

	// Only return a TLS client if its config option is enable.
	cert, err := setTLS(opts)
	if err != nil {
		return nil, err
	}

	return &websocket.Dialer{
		HandshakeTimeout: opts.WebSocket.HandshakeTimeout,
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}, nil
}

// Status returns the status info. This is used to check if the server
// is responsive and reachable.
func (c websocketClient) Status() (*scheme.Status, error) {
	return nil, nil
}

// Version returns the version info.
func (c websocketClient) Version() (*scheme.Version, error) {
	return nil, nil
}

// Config returns the unified configuration info.
func (c websocketClient) Config() (*scheme.Config, error) {
	return nil, nil
}

// Plugins returns the summary of all plugins currently registered with
// Synse Server.
func (c websocketClient) Plugins() (*[]scheme.PluginMeta, error) {
	return nil, nil
}

// Plugin returns data from a specific plugin.
func (c websocketClient) Plugin(id string) (*scheme.Plugin, error) {
	return nil, nil
}

// PluginHealth returns the summary of the health of registered plugins.
func (c websocketClient) PluginHealth() (*scheme.PluginHealth, error) {
	return nil, nil
}

// Scan returns the list of devices that Synse knows about and can read
// from/write to via the configured plugins.
// It can be filtered to show only those devices which match a set
// of provided tags by using ScanOptions.
func (c websocketClient) Scan(opts scheme.ScanOptions) (*[]scheme.Scan, error) {
	return nil, nil
}

// Tags returns the list of all tags currently associated with devices.
// If no TagsOptions is specified, the default tag namespace will be used.
func (c websocketClient) Tags(opts scheme.TagsOptions) (*[]string, error) {
	return nil, nil
}

// Info returns the full set of meta info and capabilities for a specific
// device.
func (c websocketClient) Info(id string) (*scheme.Info, error) {
	return nil, nil
}

// Read returns data from devices which match the set of provided tags
// using ReadOptions.
func (c websocketClient) Read(opts scheme.ReadOptions) (*[]scheme.Read, error) {
	return nil, nil
}

// ReadDevice returns data from a specific device.
// It is the same as Read() where the label matches the device id tag
// specified in ReadOptions.
func (c websocketClient) ReadDevice(id string, opts scheme.ReadOptions) (*[]scheme.Read, error) {
	return nil, nil
}

// ReadCache returns stream reading data from the registered plugins.
func (c websocketClient) ReadCache(opts scheme.ReadCacheOptions) (*[]scheme.Read, error) {
	return nil, nil
}

// WriteAsync writes data to a device, in an asynchronous manner.
func (c websocketClient) WriteAsync(id string, opts []scheme.WriteData) (*[]scheme.Write, error) {
	return nil, nil
}

// WriteSync writes data to a device, waiting for the write to complete.
func (c websocketClient) WriteSync(id string, opts []scheme.WriteData) (*[]scheme.Transaction, error) {
	return nil, nil
}

// Transactions returns the sorted list of all cached transaction IDs.
func (c websocketClient) Transactions() (*[]string, error) {
	return nil, nil
}

// Transaction returns the state and status of a write transaction.
func (c websocketClient) Transaction(id string) (*scheme.Transaction, error) {
	return nil, nil
}

// GetOptions returns the current config options of the client.
func (c websocketClient) GetOptions() *Options {
	return nil
}

// Close closes the connection between the client and Synse Server. This
// is only applicable for a WebSocket client in a sense that, one must
// close the connection after finish using it. Calling this method on a
// HTTP Client will have no effect.
func (c websocketClient) Close() {}
