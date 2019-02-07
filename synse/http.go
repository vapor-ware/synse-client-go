package synse

// http.go implements a http client.

import (
	"fmt"
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

// getUnversioned performs a GET request against the Synse Server unversioned API.
func (c *httpClient) getUnversioned(uri string, okScheme interface{}) error {
	errScheme := new(scheme.Error)
	_, err := c.setUnversioned().R().SetResult(okScheme).SetError(errScheme).Get(uri)
	return check(err, errScheme)
}

// getVersioned performs a GET request against the Synse Server versioned API.
func (c *httpClient) getVersioned(uri string, okScheme interface{}) error {
	errScheme := new(scheme.Error)
	client, err := c.setVersioned()
	if err != nil {
		return errors.Wrap(err, "failed to set a versioned host")
	}

	_, err = client.R().SetResult(okScheme).SetError(errScheme).Get(uri)
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
