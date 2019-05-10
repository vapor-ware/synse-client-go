package synse

// websocket.go implements a websocket client.

import (
	"crypto/tls"
	"encoding/json"
	"reflect"
	"sync/atomic"

	"github.com/vapor-ware/synse-client-go/synse/scheme"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type websocketClient struct {
	// options is the global config options of the client.
	options *Options

	// client holds the websocket.Dialer.
	client *websocket.Dialer

	// connection holds the websocket connection.
	connection *websocket.Conn

	// counter counts the number of request sent. It has the type uint64
	// that later be used by an atomic function, which makes it more
	// concurrency-safe.
	counter uint64

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
	if opts.TLS.Enabled {
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

	if !opts.TLS.Enabled {
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

			// NOTE - refer to #24. If not disable linting, a warning will happen:
			// TLS InsecureSkipVerify may be true.,HIGH,LOW (gosec)
			InsecureSkipVerify: opts.TLS.SkipVerify, // nolint
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
func (c *websocketClient) Close() error {
	err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return errors.Wrap(err, "failed to close the connection gracefully")
	}

	return nil
}

// Status returns the status info. This is used to check if the server
// is responsive and reachable.
func (c *websocketClient) Status() (*scheme.Status, error) {
	req := scheme.RequestStatus{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestStatus,
		},
	}

	resp := new(scheme.Status)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Version returns the version info.
func (c *websocketClient) Version() (*scheme.Version, error) {
	req := scheme.RequestVersion{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestVersion,
		},
	}

	resp := new(scheme.Version)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Config returns the unified configuration info.
func (c *websocketClient) Config() (*scheme.Config, error) {
	req := scheme.RequestConfig{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestConfig,
		},
	}

	resp := new(scheme.Config)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Plugins returns the summary of all plugins currently registered with
// Synse Server.
func (c *websocketClient) Plugins() ([]*scheme.PluginMeta, error) {
	req := scheme.RequestPlugins{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestPlugin,
		},
	}

	resp := new([]*scheme.PluginMeta)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// Plugin returns data from a specific plugin.
func (c *websocketClient) Plugin(id string) (*scheme.Plugin, error) {
	req := scheme.RequestPlugin{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestPlugin,
		},
		Data: scheme.PluginData{
			Plugin: id,
		},
	}

	resp := new(scheme.Plugin)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// PluginHealth returns the summary of the health of registered plugins.
func (c *websocketClient) PluginHealth() (*scheme.PluginHealth, error) {
	req := scheme.RequestPluginHealth{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestPluginHealth,
		},
	}

	resp := new(scheme.PluginHealth)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Scan returns the list of devices that Synse knows about and can read
// from/write to via the configured plugins.
// It can be filtered to show only those devices which match a set
// of provided tags by using ScanOptions.
func (c *websocketClient) Scan(opts scheme.ScanOptions) ([]*scheme.Scan, error) {
	req := scheme.RequestScan{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestScan,
		},
		Data: opts,
	}

	resp := new([]*scheme.Scan)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// Tags returns the list of all tags currently associated with devices.
// If no TagsOptions is specified, the default tag namespace will be used.
func (c *websocketClient) Tags(opts scheme.TagsOptions) ([]string, error) {
	req := scheme.RequestTags{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestTags,
		},
		Data: opts,
	}

	resp := new([]string)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// Info returns the full set of meta info and capabilities for a specific
// device.
func (c *websocketClient) Info(id string) (*scheme.Info, error) {
	req := scheme.RequestInfo{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestInfo,
		},
		Data: scheme.DeviceData{
			Device: id,
		},
	}

	resp := new(scheme.Info)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Read returns data from devices which match the set of provided tags
// using ReadOptions.
func (c *websocketClient) Read(opts scheme.ReadOptions) ([]*scheme.Read, error) {
	req := scheme.RequestRead{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestRead,
		},
		Data: opts,
	}

	resp := new([]*scheme.Read)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// ReadDevice returns data from a specific device.
// It is the same as Read() where the label matches the device id tag
// specified in ReadOptions.
func (c *websocketClient) ReadDevice(id string, opts scheme.ReadOptions) ([]*scheme.Read, error) {
	req := scheme.RequestReadDevice{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestReadDevice,
		},
		Data: scheme.ReadDeviceData{
			ID:          id,
			ReadOptions: opts,
		},
	}

	resp := new([]*scheme.Read)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// ReadCache returns stream reading data from the registered plugins.
func (c *websocketClient) ReadCache(opts scheme.ReadCacheOptions) ([]*scheme.Read, error) {
	req := scheme.RequestReadCache{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestReadCache,
		},
		Data: opts,
	}

	resp := new([]*scheme.Read)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// WriteAsync writes data to a device, in an asynchronous manner.
func (c *websocketClient) WriteAsync(id string, opts []scheme.WriteData) ([]*scheme.Write, error) {
	req := scheme.RequestWrite{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestWriteAsync,
		},
		Data: scheme.RequestWriteData{
			ID:      id,
			Payload: opts,
		},
	}

	resp := new([]*scheme.Write)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// WriteSync writes data to a device, waiting for the write to complete.
func (c *websocketClient) WriteSync(id string, opts []scheme.WriteData) ([]*scheme.Transaction, error) {
	req := scheme.RequestWrite{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestWriteSync,
		},
		Data: scheme.RequestWriteData{
			ID:      id,
			Payload: opts,
		},
	}

	resp := new([]*scheme.Transaction)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// Transactions returns the sorted list of all cached transaction IDs.
func (c *websocketClient) Transactions() ([]string, error) {
	req := scheme.RequestTransactions{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestTransaction,
		},
	}

	resp := new([]string)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// Transaction returns the state and status of a write transaction.
func (c *websocketClient) Transaction(id string) (*scheme.Transaction, error) {
	req := scheme.RequestTransaction{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestTransaction,
		},
		Data: scheme.WriteData{
			Transaction: id,
		},
	}

	resp := new(scheme.Transaction)
	err := c.makeRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetOptions returns the current config options of the client.
func (c *websocketClient) GetOptions() *Options {
	return c.options
}

// addCounter safely increases the counter by 1.
func (c *websocketClient) addCounter() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}

// makeRequest issues a request event, reads its response event and parse the
// response back.
// FIXME - refer to #22. Need to think more about how async will work in this case.
func (c *websocketClient) makeRequest(req, resp interface{}) error {
	// Write to the connection.
	err := c.connection.WriteJSON(req)
	if err != nil {
		return errors.Wrap(err, "failed to write to connection")
	}

	// Read from the connection.
	_, r, err := c.connection.NextReader()
	if err != nil {
		return errors.Wrap(err, "failed to get the next data message")
	}

	// Decode the response into a proper scheme.
	var re scheme.Response
	err = json.NewDecoder(r).Decode(&re)
	if err != nil {
		return errors.Wrap(err, "failed to decode json message into a proper scheme")
	}

	// Handle error response.
	if re.Event == responseError {
		var e scheme.Error

		err = mapstructure.Decode(re.Data, &e)
		if err != nil {
			return errors.Wrap(err, "failed to decode map into a proper scheme")
		}

		return errors.Errorf(
			"got a %v error response from synse server at %v, saying %v, with context: %v",
			e.HTTPCode, e.Timestamp, e.Description, e.Context,
		)
	}

	// Verify if the request and response metadata are matched.
	v := reflect.ValueOf(req)
	if v.FieldByName("ID").Uint() != re.ID {
		return errors.Errorf("%v did not match %v", v.FieldByName("ID"), re.ID)
	}

	if matchEvent(v.FieldByName("Event").String()) != re.Event {
		return errors.Errorf("%v did not match %v", v.FieldByName("Event"), re.Event)
	}

	// Handle successful response.
	err = mapstructure.Decode(re.Data, resp)
	if err != nil {
		return errors.Wrap(err, "failed to decode map into a proper scheme")
	}

	return nil
}

// matchEvent returns a corresponding response event for a given request event.
// FIXME - disable linting because of cyclomatic complexity due to gocyclo.
func matchEvent(reqEvent string) string { // nolint
	var respEvent string

	switch reqEvent {
	case requestStatus:
		respEvent = responseStatus
	case requestVersion:
		respEvent = responseVersion
	case requestConfig:
		respEvent = responseConfig
	case requestPlugin:
		respEvent = responsePlugin
	case requestPluginHealth:
		respEvent = responsePluginHealth
	case requestScan:
		respEvent = responseDeviceSummary
	case requestTags:
		respEvent = responseTags
	case requestInfo:
		respEvent = responseDevice
	case requestRead:
		respEvent = responseReading
	case requestReadDevice:
		respEvent = responseReading
	case requestReadCache:
		respEvent = responseReading
	case requestWriteSync:
		respEvent = responseWriteSync
	case requestWriteAsync:
		respEvent = responseWriteAsync
	case requestTransaction:
		respEvent = responseTransaction
	default:
		respEvent = ""
	}

	return respEvent
}
