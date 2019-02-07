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
		{
			"/plugin/12835beffd3e6c603aa4dd92127707b5",
			`
{
   "active":true,
   "id":"12835beffd3e6c603aa4dd92127707b5",
   "tag":"vaporio\/emulator-plugin",
   "name":"emulator plugin",
   "description":"A plugin with emulated devices and data",
   "maintainer":"vaporio",
   "vcs":"github.com\/vapor-ware\/synse-emulator-plugin",
   "version":{
      "plugin_version":"2.0.0",
      "sdk_version":"1.0.0",
      "build_date":"2018-06-14T16:24:09",
      "git_commit":"13e6478",
      "git_tag":"1.0.2-5-g13e6478",
      "arch":"amd64",
      "os":"linux"
   },
   "network":{
      "protocol":"tcp",
      "address":"emulator-plugin:5001"
   },
   "health":{
      "timestamp":"2018-06-15T20:04:33Z",
      "status":"ok",
      "message":"",
      "checks":[
         {
            "name":"read buffer health",
            "status":"ok",
            "message":"",
            "timestamp":"2018-06-15T20:04:06Z",
            "type":"periodic"
         },
         {
            "name":"write buffer health",
            "status":"ok",
            "message":"",
            "timestamp":"2018-06-15T20:04:06Z",
            "type":"periodic"
         }
      ]
   }
}`,
			&scheme.Plugin{
				PluginMeta: scheme.PluginMeta{
					Active:      true,
					ID:          "12835beffd3e6c603aa4dd92127707b5",
					Tag:         "vaporio/emulator-plugin",
					Name:        "emulator plugin",
					Description: "A plugin with emulated devices and data",
					Maintainer:  "vaporio",
					VCS:         "github.com/vapor-ware/synse-emulator-plugin",
					Version: scheme.VersionOptions{
						PluginVersion: "2.0.0",
						SDKVersion:    "1.0.0",
						BuildDate:     "2018-06-14T16:24:09",
						GitCommit:     "13e6478",
						GitTag:        "1.0.2-5-g13e6478",
						Arch:          "amd64",
						OS:            "linux",
					},
				},
				Network: scheme.NetworkOptions{
					Protocol: "tcp",
					Address:  "emulator-plugin:5001",
				},
				Health: scheme.HealthOptions{
					Timestamp: "2018-06-15T20:04:33Z",
					Status:    "ok",
					Message:   "",
					Checks: []scheme.CheckOptions{
						scheme.CheckOptions{
							Name:      "read buffer health",
							Status:    "ok",
							Message:   "",
							Timestamp: "2018-06-15T20:04:06Z",
							Type:      "periodic",
						},
						scheme.CheckOptions{
							Name:      "write buffer health",
							Status:    "ok",
							Message:   "",
							Timestamp: "2018-06-15T20:04:06Z",
							Type:      "periodic",
						},
					},
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
		case "/plugin/12835beffd3e6c603aa4dd92127707b5":
			resp, err = client.Plugin("12835beffd3e6c603aa4dd92127707b5")
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
		{"/plugin"},
		{"/plugin/12835beffd3e6c603aa4dd92127707b5"},
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
		case "/plugin":
			resp, err = client.Plugins()
		case "/plugin/12835beffd3e6c603aa4dd92127707b5":
			resp, err = client.Plugin("12835beffd3e6c603aa4dd92127707b5")
		}
		assert.Nil(t, resp)
		assert.Error(t, err)
	}
}
