package synse

// http.go implements a http client.

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/vapor-ware/synse-client-go/synse/scheme"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

// httpClient implements a http client.
type httpClient struct {
	options    *Options
	client     *resty.Client
	apiVersion string
}

// NewHTTPClient returns a new instance of a http client.
func NewHTTPClient(options *Options) (Client, error) {
	client, err := createHTTPClient(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a http client")
	}

	return &httpClient{
		options: options,
		client:  client,
	}, nil
}

// NewHTTPClientV3 returns a new instance of a http client with API v3.
func NewHTTPClientV3(options *Options) (Client, error) {
	client, err := createHTTPClient(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a http client")
	}

	return &httpClient{
		options:    options,
		client:     client,
		apiVersion: "v3",
	}, nil
}

// createHTTPClient setups the client with configured options.
func createHTTPClient(opt *Options) (*resty.Client, error) {
	if opt == nil {
		return nil, errors.New("options can not be nil")
	}

	if opt.Address == "" {
		return nil, errors.New("no address is specified")
	}

	if opt.Timeout == 0 {
		// FIXME - find a better way to use default options here.
		opt.Timeout = 2 * time.Second
	}

	// FIXME - better way to handle this?
	client := resty.New()
	client = client.SetTimeout(opt.Timeout)

	// Only use retry strategy if its options are set.
	if opt.Retry.Count != 0 {
		client = client.SetRetryCount(opt.Retry.Count)
	}

	if opt.Retry.WaitTime != 0 {
		client = client.SetRetryWaitTime(opt.Retry.WaitTime)
	}

	if opt.Retry.MaxWaitTime != 0 {
		client = client.SetRetryMaxWaitTime(opt.Retry.MaxWaitTime)
	}

	return client, nil
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
	err := c.getVersioned(configURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/config` endpoint")
	}

	return out, nil
}

// Plugins returns the summary of all plugins currently registered with
// Synse Server.
func (c *httpClient) Plugins() (*[]scheme.PluginMeta, error) {
	out := new([]scheme.PluginMeta)
	err := c.getVersioned(pluginURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/plugin` endpoint")
	}

	return out, nil
}

// Plugin returns data from a specific plugin.
func (c *httpClient) Plugin(id string) (*scheme.Plugin, error) {
	out := new(scheme.Plugin)
	err := c.getVersioned(makeURI(pluginURI, id), out)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to request `/plugin/%v` endpoint", id))
	}

	return out, nil
}

// PluginHealth returns the summary of the health of registered plugins.
func (c *httpClient) PluginHealth() (*scheme.PluginHealth, error) {
	out := new(scheme.PluginHealth)
	err := c.getVersioned(pluginHealthURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/plugin/health` endpoint")
	}

	return out, nil
}

// Scan returns the list of devices that Synse knows about and can read
// from/write to via the configured plugins. It can be filtered to show
// only those devices which match a set of provided tags by using ScanOptions.
func (c *httpClient) Scan(opts scheme.ScanOptions) (*[]scheme.Scan, error) {
	out := new([]scheme.Scan)
	err := c.getVersionedQueryParams(scanURI, opts, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/scan` endpoint")
	}

	return out, nil
}

// Tags returns the list of all tags currently associated with devices.
// If no TagsOptions is specified, the default tag namespace will be used.
func (c *httpClient) Tags(opts scheme.TagsOptions) (*[]string, error) {
	out := new([]string)
	err := c.getVersionedQueryParams(tagsURI, opts, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/tags` endpoint")
	}

	return out, nil
}

// Info returns the full set of meta info and capabilities for a specific
// device.
func (c *httpClient) Info(id string) (*scheme.Info, error) {
	out := new(scheme.Info)
	err := c.getVersioned(makeURI(infoURI, id), out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/info` endpoint")
	}

	return out, nil
}

// Read returns data from devices which match the set of provided tags
// using ReadOptions.
func (c *httpClient) Read(opts scheme.ReadOptions) (*[]scheme.Read, error) {
	out := new([]scheme.Read)
	err := c.getVersionedQueryParams(readURI, opts, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/read` endpoint")
	}

	return out, nil
}

// ReadDevice returns data from a specific device. It is the same as Read()
// where the label matches the device id tag specified in ReadOptions.
func (c *httpClient) ReadDevice(id string, opts scheme.ReadOptions) (*[]scheme.Read, error) {
	out := new([]scheme.Read)
	err := c.getVersionedQueryParams(makeURI(readURI, id), opts, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/read` endpoint")
	}

	return out, nil
}

// ReadCache returns stream reading data from the registered plugins.
func (c *httpClient) ReadCache(opts scheme.ReadCacheOptions) (*[]scheme.Read, error) {
	var out []scheme.Read
	errScheme := new(scheme.Error)

	client, err := c.setVersioned()
	if err != nil {
		return nil, errors.Wrap(err, "failed to set a versioned host while performing a GET request")
	}

	resp, err := client.R().SetDoNotParseResponse(true).SetQueryParams(structToMapString(opts)).SetError(errScheme).Get(readcacheURI)
	e := check(err, errScheme)
	if e != nil {
		return nil, errors.Wrap(err, "failed to request `/readcache` endpoint")
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

// Write writes data to a device, in an asynchronous manner.
func (c *httpClient) Write(id string, opts []scheme.WriteData) (*[]scheme.Write, error) {
	out := new([]scheme.Write)
	err := c.postVersioned(makeURI(writeURI, id), opts, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/write` endpoint")
	}

	return out, nil
}

// WriteWait writes data to a device, waiting for the write to complete.
func (c *httpClient) WriteWait(id string, opts []scheme.WriteData) (*[]scheme.Transaction, error) {
	out := new([]scheme.Transaction)
	err := c.postVersioned(makeURI(writeWaitURI, id), opts, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/write/wait` endpoint")
	}

	return out, nil
}

// Transactions returns the sorted list of all cached transaction IDs.
func (c *httpClient) Transactions() (*[]string, error) {
	out := new([]string)
	err := c.getVersioned(transactionURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/transaction` endpoint")
	}

	return out, nil
}

// Transaction returns the state and status of a write transaction.
func (c *httpClient) Transaction(id string) (*scheme.Transaction, error) {
	out := new(scheme.Transaction)
	err := c.getVersioned(makeURI(transactionURI, id), out)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to request `/transaction/%v` endpoint", id))
	}

	return out, nil
}

// Metrics returns the application-based metrics.
func (c *httpClient) Metrics() (*scheme.Metrics, error) {
	return nil, errors.New("not yet implemented")
}

// getVersionedQueryParams performs a GET request using query parameters
// against the Synse Server versioned API.
func (c *httpClient) getVersionedQueryParams(uri string, params interface{}, okScheme interface{}) error {
	errScheme := new(scheme.Error)
	client, err := c.setVersioned()
	if err != nil {
		return errors.Wrap(err, "failed to set a versioned host while performing a GET request")
	}

	_, err = client.R().SetQueryParams(structToMapString(params)).SetResult(okScheme).SetError(errScheme).Get(uri)
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
	client, err := c.setVersioned()
	if err != nil {
		return errors.Wrap(err, "failed to set a versioned host while performing a POST request")
	}

	_, err = client.R().SetBody(body).SetResult(okScheme).SetError(errScheme).Post(uri)
	return check(err, errScheme)
}

// setUnversioned returns a client that uses unversioned host URL.
func (c *httpClient) setUnversioned() *resty.Client {
	return c.client.SetHostURL(fmt.Sprintf("http://%s/", c.options.Address))
}

// setVersioned returns a client that uses versioned host URL.
func (c *httpClient) setVersioned() (*resty.Client, error) {
	err := c.cacheAPIVersion()
	if err != nil {
		return nil, errors.Wrap(err, "failed to cache api version")
	}

	return c.client.SetHostURL(fmt.Sprintf("http://%s/%s/", c.options.Address, c.apiVersion)), nil
}

// cacheAPIVersion caches the api version if not already.
func (c *httpClient) cacheAPIVersion() error {
	if c.apiVersion == "" {
		client, err := c.Version()
		if err != nil {
			return errors.Wrap(err, "failed to get synse version for caching")
		}

		c.apiVersion = client.APIVersion
	}

	return nil
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

// makeURI joins the given components into a string, delimited with '/' which
// can then be used as the URI for API requests.
func makeURI(components ...string) string {
	return strings.Join(components, "/")
}

// structToMapString decodes a struct value into a map[string]string that can
// be used as query parameters. It assumes that the struct fields follow one of
// these types: bool, string, int, float, slice.
func structToMapString(s interface{}) map[string]string {
	out := map[string]string{}
	v := ""

	fields := reflect.TypeOf(s)
	values := reflect.ValueOf(s)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)

		switch value.Kind() {
		case reflect.Slice:
			s := []string{}
			for i := 0; i < value.Len(); i++ {
				s = append(s, fmt.Sprint(value.Index(i)))
			}

			v = strings.Join(s, ",")
		default:
			v = fmt.Sprint(value)
		}

		out[strings.ToLower(field.Name)] = v
	}

	return out
}
