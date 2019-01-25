package synse

// http.go implements a http client.

import (
	"fmt"
	"time"

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
		return nil, errors.Wrap(err, "failed to setup a http client")
	}

	return &httpClient{
		options: options,
		client:  client,
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

	client := resty.New()
	client = client.SetHostURL(fmt.Sprintf("http://%s/synse/", opt.Address))

	if opt.Timeout == 0 {
		// FIXME - find a better way to use default options here?
		opt.Timeout = 2 * time.Second
	}
	client = client.SetTimeout(opt.Timeout)

	// Only use retry options if set, otherwise let the resty client goes with
	// its defaults (O Count, 100 milliseconds WaitTime, 2 seconds MaxWaitTime).
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
func (c *httpClient) Status() (*Status, error) {
	out := new(Status)

	err := c.getUnversioned(testURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/status` endpoint")
	}

	return out, nil
}

// Version returns the version info.
func (c *httpClient) Version() (*Version, error) {
	out := new(Version)

	err := c.getUnversioned(versionURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/version` endpoint")
	}

	return out, nil
}

// Config returns the config info.
func (c *httpClient) Config() (*Config, error) {
	out := new(Config)

	err := c.getVersioned(configURI, out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/config` endpoint")
	}

	return out, nil
}

// getUnversioned makes an unversioned request.
func (c *httpClient) getUnversioned(path string, okResp interface{}) error {
	errResp := new(Error)

	_, err := c.client.R().
		SetResult(okResp).
		SetError(errResp).
		Get(path)
	if err != nil {
		return errors.Wrap(err, "failed to make an unversioned request to synse server")
	}

	if *errResp != (Error{}) {
		return errors.Errorf("got an error response from synse server: %+v", errResp)
	}

	return nil
}

// getVersioned makes a versioned request.
func (c *httpClient) getVersioned(path string, okResp interface{}) error {
	err := c.cacheAPIVersion()
	if err != nil {
		return errors.Wrap(err, "failed to cache api version")
	}

	versionedPath := fmt.Sprintf("%s/%s", c.apiVersion, path)
	err = c.getUnversioned(versionedPath, okResp)
	if err != nil {
		return errors.Wrap(err, "failed to make a versioned request to synse server")
	}

	return nil
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
