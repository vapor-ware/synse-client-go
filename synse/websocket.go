package synse

// websocket.go implements a websocket client.

import (
	"crypto/tls"
	"reflect"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
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
	// FIXME (etd): I think this will panic if close is called before connect.
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
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
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
			Event: requestPlugins,
		},
	}

	resp := new([]*scheme.PluginMeta)
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// Info returns the full set of meta info and capabilities for a specific
// device.
func (c *websocketClient) Info(device string) (*scheme.Info, error) {
	req := scheme.RequestInfo{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestInfo,
		},
		Data: scheme.DeviceData{
			Device: device,
		},
	}

	resp := new(scheme.Info)
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// ReadDevice returns data from a specific device.
// It is the same as Read() where the label matches the device id tag
// specified in ReadOptions.
func (c *websocketClient) ReadDevice(device string) ([]*scheme.Read, error) {
	req := scheme.RequestReadDevice{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestReadDevice,
		},
		Data: scheme.ReadDeviceData{
			Device: device,
		},
	}

	resp := new([]*scheme.Read)
	err := c.makeRequestResponse(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// ReadCache returns cached reading data from the registered plugins.
func (c *websocketClient) ReadCache(opts scheme.ReadCacheOptions, out chan<- *scheme.Read) error {
	defer close(out)

	req := scheme.RequestReadCache{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestReadCache,
		},
		Data: opts,
	}

	resp := new([]*scheme.Read)
	err := c.makeRequestResponse(req, resp)
	if err != nil {
		return err
	}
	for _, r := range *resp {
		out <- r
	}
	return nil
}

// ReadStream returns a stream of current reading data from the registered plugins.
func (c *websocketClient) ReadStream(opts scheme.ReadStreamOptions, out chan<- *scheme.Read, stop chan struct{}) error {

	req := scheme.RequestReadStream{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestReadStream,
		},
		Data: opts,
	}

	resp := &scheme.Read{}
	proxy := make(chan interface{})
	go func() {
		for {
			data, open := <-proxy
			if !open {
				return
			}
			out <- data.(*scheme.Read)
		}
	}()

	err := c.streamRequest(req, resp, proxy, stop)
	if err != nil {
		return errors.Wrap(err, "failed to stream reading data")
	}

	// If we got here, the stream has been terminated (e.g. via the `stop` channel)
	// and the WebSocket session is still active (e.g. we did not error, panic, or
	// exit the program). The client has stopped listening for readings, but we must
	// tell the server to stop sending them as well.
	req = scheme.RequestReadStream{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestReadStream,
		},
		Data: scheme.ReadStreamOptions{
			Stop: true,
		},
	}

	err = c.makeRequest(req)
	if err != nil {
		return errors.Wrap(err, "failed to stop server-side read stream")
	}
	return nil
}

// WriteAsync writes data to a device, in an asynchronous manner.
func (c *websocketClient) WriteAsync(device string, opts []scheme.WriteData) ([]*scheme.Write, error) {
	req := scheme.RequestWrite{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestWriteAsync,
		},
		Data: scheme.RequestWriteData{
			Device:  device,
			Payload: opts,
		},
	}

	resp := new([]*scheme.Write)
	err := c.makeRequestResponse(req, resp)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

// WriteSync writes data to a device, waiting for the write to complete.
func (c *websocketClient) WriteSync(device string, opts []scheme.WriteData) ([]*scheme.Transaction, error) {
	req := scheme.RequestWrite{
		EventMeta: scheme.EventMeta{
			ID:    c.addCounter(),
			Event: requestWriteSync,
		},
		Data: scheme.RequestWriteData{
			Device:  device,
			Payload: opts,
		},
	}

	resp := new([]*scheme.Transaction)
	err := c.makeRequestResponse(req, resp)
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
			Event: requestTransactions,
		},
	}

	resp := new([]string)
	err := c.makeRequestResponse(req, resp)
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
	err := c.makeRequestResponse(req, resp)
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
func (c *websocketClient) makeRequestResponse(req, resp interface{}) error {
	// Write to the connection.
	err := c.connection.WriteJSON(req)
	if err != nil {
		return errors.Wrap(err, "failed to write to connection")
	}

	return c.readResponse(req, resp)
}

// makeRequest issues a request event. It does not attempt to read back a response
// for the request.
func (c *websocketClient) makeRequest(req interface{}) error {
	err := c.connection.WriteJSON(req)
	if err != nil {
		return errors.Wrap(err, "failed to issue request")
	}
	return nil
}

// readResponse reads the response for a given request.
func (c *websocketClient) readResponse(req, resp interface{}) error {
	// Read from the connection.
	var re scheme.Response
	err := c.connection.ReadJSON(&re)
	if err != nil {
		return errors.Wrap(err, "failed to read response message")
	}
	return c.parseResponseMessage(re, req, resp)
}

func (c *websocketClient) streamRequest(req, resp interface{}, stream chan interface{}, stop chan struct{}) error {
	defer close(stream)

	err := c.connection.WriteJSON(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	for {
		// If the stop channel is closed, terminate the stream.
		select {
		case _, open := <-stop:
			if !open {
				return nil
			}
		default:
			// Do nothing and continue on. The next iteration of the
			// loop will check again whether the stream should stop.
		}

		var response scheme.Response
		err := c.connection.ReadJSON(&response)
		if err != nil {
			return errors.Wrap(err, "failed to read data from stream")
		}

		respInst := reflect.New(reflect.TypeOf(resp).Elem()).Interface()
		err = c.parseResponseMessage(response, req, respInst)
		if err != nil {
			return errors.Wrap(err, "failed to parse response message")
		}

		stream <- respInst
	}
}

func (c *websocketClient) parseResponseMessage(r scheme.Response, req, resp interface{}) error {
	if r.Event == responseError {
		var e scheme.Error

		err := mapstructure.Decode(r.Data, &e)
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
	if v.FieldByName("ID").Uint() != r.ID {
		return errors.Errorf("response id mismatch: %v != %v", v.FieldByName("ID"), r.ID)
	}

	if matchEvent(v.FieldByName("Event").String()) != r.Event {
		return errors.Errorf("(%v) %v did not match %v", v.FieldByName("Event"), matchEvent(v.FieldByName("Event").String()), r.Event)
	}

	// Handle successful response.
	err := mapstructure.Decode(r.Data, resp)
	if err != nil {
		return errors.Wrap(err, "failed to decode map into a proper scheme")
	}

	return nil
}

// matchEvent returns a corresponding response event for a given request event.
func matchEvent(reqEvent string) string {
	switch reqEvent {
	case requestStatus:
		return responseStatus
	case requestVersion:
		return responseVersion
	case requestConfig:
		return responseConfig
	case requestPlugin:
		return responsePluginInfo
	case requestPlugins:
		return responsePluginSummary
	case requestPluginHealth:
		return responsePluginHealth
	case requestScan:
		return responseDeviceSummary
	case requestTags:
		return responseTags
	case requestInfo:
		return responseDeviceInfo
	case requestRead:
		return responseReading
	case requestReadDevice:
		return responseReading
	case requestReadStream:
		return responseReading
	case requestReadCache:
		return responseReading
	case requestWriteSync:
		return responseTransactionStatus
	case requestWriteAsync:
		return responseTransactionInfo
	case requestTransactions:
		return responseTransactionList
	case requestTransaction:
		return responseTransactionStatus
	default:
		return ""
	}
}
