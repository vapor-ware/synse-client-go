package synse

import (
	"testing"
	"time"

	"github.com/vapor-ware/synse-client-go/internal/test"
	"github.com/vapor-ware/synse-client-go/synse/scheme"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClient_NilConfig(t *testing.T) {
	client, err := NewHTTPClient(nil)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewHTTPClient_NoAddress(t *testing.T) {
	client, err := NewHTTPClient(&Options{
		Address: "",
	})
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewHTTPClient_ValidAddress(t *testing.T) {
	client, err := NewHTTPClient(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestNewHTTPClient_ValidAddressAndTimeout(t *testing.T) {
	client, err := NewHTTPClient(&Options{
		Address: "localhost:5000",
		Timeout: 3 * time.Second,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestNewHTTPClient_ValidRetry(t *testing.T) {
	client, err := NewHTTPClient(&Options{
		Address: "localhost:5000",
		Retry: RetryOptions{
			Count:       3,
			WaitTime:    5 * time.Second,
			MaxWaitTime: 20 * time.Second,
		},
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestHTTPClient_Unversioned_200(t *testing.T) {
	tests := []struct {
		path     string
		in       string
		expected interface{}
	}{
		{
			"/test",
			`
{
  "status":"ok",
  "timestamp":"2019-01-24T14:34:24.926108Z"
}`,
			&scheme.Status{
				Status:    "ok",
				Timestamp: "2019-01-24T14:34:24.926108Z",
			},
		},
		{
			"/version",
			`
{
  "version":"3.0.0",
  "api_version":"v3"
}`,
			&scheme.Version{
				Version:    "3.0.0",
				APIVersion: "v3",
			},
		},
	}

	server := test.NewMockHTTPServer()
	defer server.Close()

	client, err := NewHTTPClient(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	for _, tt := range tests {
		server.ServeUnversionedSuccess(t, tt.path, tt.in)

		var (
			resp interface{}
			err  error
		)
		switch tt.path {
		case "/test":
			resp, err = client.Status()
		case "/version":
			resp, err = client.Version()
		}
		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, resp)
	}
}

func TestHTTPClient_Unversioned_500(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"/test"},
		{"/version"},
	}

	server := test.NewMockHTTPServer()
	defer server.Close()

	client, err := NewHTTPClient(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	in := `
{
  "http_code":500,
  "error_id":0,
  "description":"unknown error",
  "timestamp":"2019-01-24T14:36:53.166038Z",
  "context":"unknown error"
}
`
	for _, tt := range tests {
		server.ServeUnversionedFailure(t, tt.path, in)

		var (
			resp interface{}
			err  error
		)
		switch tt.path {
		case "/test":
			resp, err = client.Status()
		case "/version":
			resp, err = client.Version()
		}
		assert.Nil(t, resp)
		assert.Error(t, err)
	}
}

func TestHTTPClient_Versioned_200(t *testing.T) {
	tests := []struct {
		path     string
		in       string
		expected interface{}
	}{
		{
			"/config",
			`
{
  "logging":"info",
  "pretty_json":true,
  "locale":"en_US",
  "plugin":{
    "tcp":[

    ],
    "unix":[

    ]
  },
  "cache":{
    "meta":{
      "ttl":20
    },
    "transaction":{
      "ttl":300
    }
  },
  "grpc":{
    "timeout":3
  }
}`,
			&scheme.Config{
				Logging:    "info",
				PrettyJSON: true,
				Locale:     "en_US",
				Plugin: scheme.PluginOptions{
					TCP:  []string{},
					Unix: []string{},
				},
				Cache: scheme.CacheOptions{
					Meta: scheme.MetaOptions{
						TTL: int(20),
					},
					Transaction: scheme.TransactionOptions{
						TTL: int(300),
					},
				},
				GRPC: scheme.GRPCOptions{
					Timeout: int(3),
				},
			},
		},
	}

	server := test.NewMockHTTPServer()
	defer server.Close()

	client, err := NewHTTPClient(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	for _, tt := range tests {
		server.ServeVersionedSuccess(t, tt.path, tt.in)

		var (
			resp interface{}
			err  error
		)
		switch tt.path {
		case "/config":
			resp, err = client.Config()
		}
		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, resp)
	}
}

func TestHTTPClient_Versioned_500(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"/config"},
	}

	server := test.NewMockHTTPServer()
	defer server.Close()

	client, err := NewHTTPClient(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	in := `
{
  "http_code":500,
  "error_id":0,
  "description":"unknown error",
  "timestamp":"2019-01-24T14:36:53.166038Z",
  "context":"unknown error"
}
`
	for _, tt := range tests {
		server.ServeVersionedFailure(t, tt.path, in)

		var (
			resp interface{}
			err  error
		)
		switch tt.path {
		case "/config":
			resp, err = client.Config()
		}
		assert.Nil(t, resp)
		assert.Error(t, err)
	}
}
