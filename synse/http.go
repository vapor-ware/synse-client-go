package synse

// http.go implements a http client.

import (
	"crypto/tls"
	"encoding/json"

	"github.com/vapor-ware/synse-client-go/synse/scheme"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// httpClient implements a http client.
type httpClient struct {
	// options is the global config options of the client.
	options *Options

	// client holds the resty.Client.
	client *resty.Client

	// apiVersion is the current api version of Synse Server that we are
	// communicating with.
	apiVersion string

	// scheme could either be `http` or `https`, depends on the TLS
	// configuration.
	scheme string
}

// NewHTTPClientV3 returns a new instance of a http client for v3 API.
func NewHTTPClientV3(options *Options) (Client, error) {
	scheme := "http"
	client, err := createHTTPClient(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a http client")
	}

	// Check if TLS options are set.
	if options.TLS.CertFile != "" && options.TLS.KeyFile != "" {
		// Change the scheme to `https`
		scheme = "https"

		// Register the certificates.
		cert, err := tls.LoadX509KeyPair(options.TLS.CertFile, options.TLS.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set client certificates")
		}

		client.SetCertificates(cert)

		// Set the security check.
		// FIXME - if not disable linting here, it will yield: warning: TLS
		// InsecureSkipVerify may be true.,HIGH,LOW (gosec)
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: options.TLS.SkipVerify}) // nolint
	}

	return &httpClient{
		options:    options,
		client:     client,
		apiVersion: "v3",
		scheme:     scheme,
	}, nil
}

// createHTTPClient setups the client with configured options.
func createHTTPClient(opts *Options) (*resty.Client, error) {
	if opts == nil {
		return nil, errors.New("options can not be nil")
	}

	if opts.Address == "" {
		return nil, errors.New("no address is specified")
	}

	err := defaults.Set(opts)
	if err != nil {
		return nil, errors.New("failed to set default configs")
	}

	// Create a new resty client with configured options.
	client := resty.New()
	return client.
		SetTimeout(opts.Timeout).
		SetRetryCount(opts.Retry.Count).
		SetRetryWaitTime(opts.Retry.WaitTime).
		SetRetryMaxWaitTime(opts.Retry.MaxWaitTime), nil
}

// Status returns the status info.
func (c *httpClient) Status() (*scheme.Status, error) {
	out := new(scheme.Status)
	err := c.getUnversioned(testURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/status` endpoint")
	}

	return out, nil
}

// Version returns the version info.
func (c *httpClient) Version() (*scheme.Version, error) {
	out := new(scheme.Version)
	err := c.getUnversioned(versionURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/version` endpoint")
	}

	return out, nil
}

// Config returns the config info.
func (c *httpClient) Config() (*scheme.Config, error) {
	out := new(scheme.Config)
	if err := c.getVersioned(configURI, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Plugins returns the summary of all plugins currently registered with
// Synse Server.
func (c *httpClient) Plugins() (*[]scheme.PluginMeta, error) {
	out := new([]scheme.PluginMeta)
	if err := c.getVersioned(pluginURI, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Plugin returns data from a specific plugin.
func (c *httpClient) Plugin(id string) (*scheme.Plugin, error) {
	out := new(scheme.Plugin)
	if err := c.getVersioned(makePath(pluginURI, id), out); err != nil {
		return nil, err
	}

	return out, nil
}

// PluginHealth returns the summary of the health of registered plugins.
func (c *httpClient) PluginHealth() (*scheme.PluginHealth, error) {
	out := new(scheme.PluginHealth)
	if err := c.getVersioned(pluginHealthURI, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Scan returns the list of devices that Synse knows about and can read
// from/write to via the configured plugins. It can be filtered to show
// only those devices which match a set of provided tags by using ScanOptions.
func (c *httpClient) Scan(opts scheme.ScanOptions) (*[]scheme.Scan, error) {
	out := new([]scheme.Scan)
	if err := c.getVersionedQueryParams(scanURI, opts, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Tags returns the list of all tags currently associated with devices.
// If no TagsOptions is specified, the default tag namespace will be used.
func (c *httpClient) Tags(opts scheme.TagsOptions) (*[]string, error) {
	out := new([]string)
	if err := c.getVersionedQueryParams(tagsURI, opts, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Info returns the full set of meta info and capabilities for a specific
// device.
func (c *httpClient) Info(id string) (*scheme.Info, error) {
	out := new(scheme.Info)
	if err := c.getVersioned(makePath(infoURI, id), out); err != nil {
		return nil, err
	}

	return out, nil
}

// Read returns data from devices which match the set of provided tags
// using ReadOptions.
func (c *httpClient) Read(opts scheme.ReadOptions) (*[]scheme.Read, error) {
	out := new([]scheme.Read)
	if err := c.getVersionedQueryParams(readURI, opts, out); err != nil {
		return nil, err
	}

	return out, nil
}

// ReadDevice returns data from a specific device. It is the same as Read()
// where the label matches the device id tag specified in ReadOptions.
func (c *httpClient) ReadDevice(id string, opts scheme.ReadOptions) (*[]scheme.Read, error) {
	out := new([]scheme.Read)
	if err := c.getVersionedQueryParams(makePath(readURI, id), opts, out); err != nil {
		return nil, err
	}

	return out, nil
}

// ReadCache returns stream reading data from the registered plugins.
func (c *httpClient) ReadCache(opts scheme.ReadCacheOptions) (*[]scheme.Read, error) {
	var out []scheme.Read
	errScheme := new(scheme.Error)

	resp, err := c.setVersioned().R().SetDoNotParseResponse(true).SetQueryParams(structToMapString(opts)).SetError(errScheme).Get(readcacheURI)
	if err = check(err, errScheme); err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.RawBody())
	for dec.More() {
		var read scheme.Read
		err := dec.Decode(&read)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode a JSON response into an appropriate struct")
		}

		out = append(out, read)
	}

	return &out, nil
}

// WriteAsync writes data to a device, in an asynchronous manner.
func (c *httpClient) WriteAsync(id string, opts []scheme.WriteData) (*[]scheme.Write, error) {
	out := new([]scheme.Write)
	if err := c.postVersioned(makePath(writeURI, id), opts, out); err != nil {
		return nil, err
	}

	return out, nil
}

// WriteSync writes data to a device, waiting for the write to complete.
func (c *httpClient) WriteSync(id string, opts []scheme.WriteData) (*[]scheme.Transaction, error) {
	out := new([]scheme.Transaction)
	if err := c.postVersioned(makePath(writeWaitURI, id), opts, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Transactions returns the sorted list of all cached transaction IDs.
func (c *httpClient) Transactions() (*[]string, error) {
	out := new([]string)
	if err := c.getVersioned(transactionURI, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Transaction returns the state and status of a write transaction.
func (c *httpClient) Transaction(id string) (*scheme.Transaction, error) {
	out := new(scheme.Transaction)
	if err := c.getVersioned(makePath(transactionURI, id), out); err != nil {
		return nil, err
	}

	return out, nil
}

// GetOptions returns the current config options of the client.
func (c *httpClient) GetOptions() *Options {
	return c.options
}

// Close closes the connection between the client and Synse Server. However,
// this method is not applicable for a HTTP client. Hence, it will do nothing.
func (c *httpClient) Close() {}

// getVersionedQueryParams performs a GET request using query parameters
// against the Synse Server versioned API.
func (c *httpClient) getVersionedQueryParams(uri string, params interface{}, okScheme interface{}) error {
	errScheme := new(scheme.Error)
	_, err := c.setVersioned().R().SetQueryParams(structToMapString(params)).SetResult(okScheme).SetError(errScheme).Get(uri)
	return check(err, errScheme)

}

// getVersioned performs a GET request against the Synse Server versioned API.
func (c *httpClient) getVersioned(uri string, okScheme interface{}) error {
	params := struct{}{}
	return c.getVersionedQueryParams(uri, params, okScheme)
}

// getUnversioned performs a GET request against the Synse Server unversioned API.
func (c *httpClient) getUnversioned(uri string, okScheme interface{}) error {
	errScheme := new(scheme.Error)
	_, err := c.setUnversioned().R().SetResult(okScheme).SetError(errScheme).Get(uri)
	return check(err, errScheme)
}

// postVersioned performs a POST request against the Synse Server versioned API.
func (c *httpClient) postVersioned(uri string, body interface{}, okScheme interface{}) error {
	errScheme := new(scheme.Error)
	_, err := c.setVersioned().R().SetBody(body).SetResult(okScheme).SetError(errScheme).Post(uri)
	return check(err, errScheme)
}

// setUnversioned returns a client that uses unversioned host URL.
func (c *httpClient) setUnversioned() *resty.Client {
	return c.client.SetHostURL(buildURL(c.scheme, c.options.Address))
}

// setVersioned returns a client that uses versioned host URL.
func (c *httpClient) setVersioned() *resty.Client {
	return c.client.SetHostURL(buildURL(c.scheme, c.options.Address, c.apiVersion))
}

// check validates returned response from the Synse Server.
func check(err error, errResp *scheme.Error) error {
	if err != nil {
		return errors.Wrap(err, "failed to make a request to synse server")
	}

	if *errResp != (scheme.Error{}) {
		return errors.Errorf(
			"got a %v error response from synse server at %v, saying %v, with context: %v",
			errResp.HTTPCode, errResp.Timestamp, errResp.Description, errResp.Context,
		)

	}

	return nil
}
