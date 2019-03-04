package synse

// websocket.go implements a websocket client.

import (
	"crypto/tls"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
)

// counter counts the number of request sent.
var counter uint64

type websocketClient struct {
	// options is the global config options of the client.
	options *Options

	// client holds the websocket.Dialer.
	client *websocket.Dialer

	// connection holds the websocket connection.
	connection *websocket.Conn

	// apiVersion is the current api version of Synse Server that we are
	// communicating with.
	apiVersion string

	// entryRoute is the entry route to start the websocket connection.
	entryRoute string

	// scheme could either be ws or wss depending on TLS configuration.
	scheme string
}

// NewWebSocketClientV3 returns a new instance of a websocket client for v3.
func NewWebSocketClientV3(opts *Options) (Client, error) {
	c, err := createWebSocketClient(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a websocket client")
	}

	s := "ws"
	if opts.TLS.Enabled == true {
		s = "wss"
	}

	return &websocketClient{
		options:    opts,
		client:     c,
		apiVersion: "v3",
		entryRoute: "connect",
		scheme:     s,
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

// Open opens the websocket connection between the client and Synse Server.
func (c *websocketClient) Open() error {
	conn, _, err := c.client.Dial(buildURL(c.scheme, c.options.Address, c.apiVersion, c.entryRoute), nil)
	if err != nil {
		return errors.Wrap(err, "failed to open the websocket connection")
	}

	c.connection = conn
	return nil
}

// Close closes the websocket connection between the client and Synse Server.
// It's up to the user to close the connection after finish using it.
// TODO - look into how defer work with returned error?
func (c *websocketClient) Close() error {
	err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return errors.Wrap(err, "failed to close the connection gracefully")
	}

	return nil
}

// Status returns the status info. This is used to check if the server
// is responsive and reachable.
// Note: this method is not applicable for websocket client.
func (c *websocketClient) Status() (*scheme.Status, error) {
	return nil, nil
}

// Version returns the version info.
func (c *websocketClient) Version() (*scheme.Version, error) {
	req := scheme.RequestVersion{
		EventMeta: scheme.EventMeta{
			ID:    addCounter(),
			Event: requestVersion,
		},
	}

	err := c.connection.WriteJSON(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write to connection")
	}

	resp := new(scheme.ResponseVersion)
	err = c.connection.ReadJSON(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from connection")
	}

	// Compare the id between request and response event.
	if req.ID != resp.ID {
		return nil, errors.Wrap(err, "request id doesn't match")
	}

	return &resp.Data, nil
}

// Config returns the unified configuration info.
func (c *websocketClient) Config() (*scheme.Config, error) {
	return nil, nil
}

// Plugins returns the summary of all plugins currently registered with
// Synse Server.
// Note: this method is not applicable for websocket client.
func (c *websocketClient) Plugins() (*[]scheme.PluginMeta, error) {
	return nil, nil
}

// Plugin returns data from a specific plugin.
func (c *websocketClient) Plugin(id string) (*scheme.Plugin, error) {
	req := scheme.RequestPlugin{
		EventMeta: scheme.EventMeta{
			ID:    addCounter(),
			Event: requestVersion,
		},
		Data: scheme.PluginData{
			Plugin: id,
		},
	}

	err := c.connection.WriteJSON(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write to connection")
	}

	resp := new(scheme.ResponsePlugin)
	err = c.connection.ReadJSON(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from connection")
	}

	// Compare the id between request and response event.
	if req.ID != resp.ID {
		return nil, errors.Wrap(err, "request id doesn't match")
	}

	return &resp.Data, nil

}

// PluginHealth returns the summary of the health of registered plugins.
func (c *websocketClient) PluginHealth() (*scheme.PluginHealth, error) {
	return nil, nil
}

// Scan returns the list of devices that Synse knows about and can read
// from/write to via the configured plugins.
// It can be filtered to show only those devices which match a set
// of provided tags by using ScanOptions.
func (c *websocketClient) Scan(opts scheme.ScanOptions) (*[]scheme.Scan, error) {
	return nil, nil
}

// Tags returns the list of all tags currently associated with devices.
// If no TagsOptions is specified, the default tag namespace will be used.
func (c *websocketClient) Tags(opts scheme.TagsOptions) (*[]string, error) {
	return nil, nil
}

// Info returns the full set of meta info and capabilities for a specific
// device.
func (c *websocketClient) Info(id string) (*scheme.Info, error) {
	return nil, nil
}

// Read returns data from devices which match the set of provided tags
// using ReadOptions.
func (c *websocketClient) Read(opts scheme.ReadOptions) (*[]scheme.Read, error) {
	return nil, nil
}

// ReadDevice returns data from a specific device.
// It is the same as Read() where the label matches the device id tag
// specified in ReadOptions.
func (c *websocketClient) ReadDevice(id string, opts scheme.ReadOptions) (*[]scheme.Read, error) {
	return nil, nil
}

// ReadCache returns stream reading data from the registered plugins.
func (c *websocketClient) ReadCache(opts scheme.ReadCacheOptions) (*[]scheme.Read, error) {
	return nil, nil
}

// WriteAsync writes data to a device, in an asynchronous manner.
func (c *websocketClient) WriteAsync(id string, opts []scheme.WriteData) (*[]scheme.Write, error) {
	return nil, nil
}

// WriteSync writes data to a device, waiting for the write to complete.
func (c *websocketClient) WriteSync(id string, opts []scheme.WriteData) (*[]scheme.Transaction, error) {
	return nil, nil
}

// Transactions returns the sorted list of all cached transaction IDs.
func (c *websocketClient) Transactions() (*[]string, error) {
	return nil, nil
}

// Transaction returns the state and status of a write transaction.
func (c *websocketClient) Transaction(id string) (*scheme.Transaction, error) {
	return nil, nil
}

// GetOptions returns the current config options of the client.
func (c *websocketClient) GetOptions() *Options {
	return c.options
}

// addCounter safely increases the counter by 1.
func addCounter() uint64 {
	return atomic.AddUint64(&counter, 1)
}

// getCounter safely gets current value of counter.
func getCounter() uint64 {
	return atomic.LoadUint64(&counter)
}
