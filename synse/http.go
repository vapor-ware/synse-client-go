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
	client := resty.New()
	err := setupHTTPClient(options, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup a http client")
	}

	return &httpClient{
		options: options,
		client:  client,
	}, nil
}

// setupHTTPClient setups the client with configured options.
func setupHTTPClient(opt *Options, c *resty.Client) error {
	if opt.Server.Address == "" {
		return errors.New("no address is specified")
	}
	c = c.SetHostURL(fmt.Sprintf("http://%s/synse/", opt.Server.Address))

	if opt.Server.Timeout == 0 {
		opt.Server.Timeout = 2 * time.Second
		c = c.SetTimeout(opt.Server.Timeout)
	}

	if opt.Retry.Count != 0 {
		c = c.SetRetryCount(opt.Retry.Count)
	}

	if opt.Retry.WaitTime != 0 {
		c = c.SetRetryWaitTime(opt.Retry.WaitTime)
	}

	if opt.Retry.MaxWaitTime != 0 {
		c = c.SetRetryMaxWaitTime(opt.Retry.MaxWaitTime)
	}

	return nil
}

// Status returns the status info.
func (c *httpClient) Status() (*Status, error) {
	out := new(Status)
	synseError := new(Error)

	err := c.getUnversioned(testURI, out, synseError)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/status` endpoint")
	}

	return out, nil
}

// Version returns the version info.
func (c *httpClient) Version() (*Version, error) {
	out := new(Version)
	synseError := new(Error)

	err := c.getUnversioned(versionURI, out, synseError)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/version` endpoint")
	}

	return out, nil
}

// Config returns the config info.
func (c *httpClient) Config() (*Config, error) {
	out := new(Config)
	synseError := new(Error)

	err := c.getVersioned(configURI, out, synseError)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request `/config` endpoint")
	}

	return out, nil
}

// getUnversioned makes an unversioned request.
func (c *httpClient) getUnversioned(path string, okResp interface{}, errResp *Error) error {
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
func (c *httpClient) getVersioned(path string, okResp interface{}, errResp *Error) error {
	err := c.cacheAPIVersion()
	if err != nil {
		return errors.Wrap(err, "failed to cache api version")
	}

	versionedPath := fmt.Sprintf("%s/%s", c.apiVersion, path)
	err = c.getUnversioned(versionedPath, okResp, errResp)
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
