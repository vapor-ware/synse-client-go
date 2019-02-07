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

	server := test.NewUnversionedHTTPServer()
	defer server.Close()

	client, err := NewHTTPClient(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	for _, tt := range tests {
		server.Serve(t, tt.path, 200, tt.in)

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

	server := test.NewUnversionedHTTPServer()
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
		server.Serve(t, tt.path, 500, in)

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
   "cache":{
      "device":{
         "ttl":20
      },
      "transaction":{
         "ttl":300
      }
   },
   "grpc":{
      "timeout":3,
      "tls":{
         "cert":"/tmp/ssl/synse.crt"
      }
   },
   "metrics":{
      "enabled":false
   },
   "transport":{
      "http":true,
      "websocket":true
   },
   "plugin":{
      "tcp":[
         "emulator-plugin:5001"
      ],
      "unix":[
         "/tmp/synse/plugin/foo.sock"
      ],
      "discover":{
         "kubernetes":{
            "namespace":"vapor",
            "endpoints":{
               "labels":{
                  "app":"synse",
                  "component":"server"
               }
            }
         }
      }
   }
}`,
			&scheme.Config{
				Logging:    "info",
				PrettyJSON: true,
				Locale:     "en_US",
				Cache: scheme.CacheOptions{
					Device: scheme.DeviceOptions{
						TTL: int(20),
					},
					Transaction: scheme.TransactionOptions{
						TTL: int(300),
					},
				},
				GRPC: scheme.GRPCOptions{
					Timeout: int(3),
					TLS: scheme.TLSOptions{
						Cert: "/tmp/ssl/synse.crt",
					},
				},
				Metrics: scheme.MetricsOptions{
					Enabled: false,
				},
				Transport: scheme.TransportOptions{
					HTTP:      true,
					WebSocket: true,
				},
				Plugin: scheme.PluginOptions{
					TCP:  []string{"emulator-plugin:5001"},
					Unix: []string{"/tmp/synse/plugin/foo.sock"},
					Discover: scheme.DiscoveryOptions{
						Kubernetes: scheme.KubernetesOptions{
							Namespace: "vapor",
							Endpoints: scheme.EndpointsOptions{
								Labels: map[string]string{
									"app":       "synse",
									"component": "server",
								},
							},
						},
					},
				},
			},
		},
		{
			"/plugin",
			`
		[
		{
		"description": "a plugin with emulated devices and data",
		"id": "12835beffd3e6c603aa4dd92127707b5",
		"name": "emulator plugin",
		"maintainer": "vapor io",
		"active": true
		},
		{
		"description": "a custom third party plugin",
		"id": "12835beffd3e6c603aa4dd92127707b6",
		"name": "custom-plugin",
		"maintainer": "third-party",
		"active": true
		}
		]`,
			&[]scheme.PluginMeta{
				scheme.PluginMeta{
					Description: "a plugin with emulated devices and data",
					ID:          "12835beffd3e6c603aa4dd92127707b5",
					Name:        "emulator plugin",
					Maintainer:  "vapor io",
					Active:      true,
				},
				scheme.PluginMeta{
					Description: "a custom third party plugin",
					ID:          "12835beffd3e6c603aa4dd92127707b6",
					Name:        "custom-plugin",
					Maintainer:  "third-party",
					Active:      true,
				},
			},
		},
	}

	server := test.NewVersionedHTTPServer()
	defer server.Close()

	client, err := NewHTTPClient(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	for _, tt := range tests {
		server.Serve(t, tt.path, 200, tt.in)

		var (
			resp interface{}
			err  error
		)
		switch tt.path {
		case "/config":
			resp, err = client.Config()
		case "/plugin":
			resp, err = client.Plugins()
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

	server := test.NewVersionedHTTPServer()
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
		server.Serve(t, tt.path, 500, in)

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
