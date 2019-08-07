package synse

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-client-go/internal/test"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
)

func TestNewHTTPClientV3_NilConfig(t *testing.T) {
	client, err := NewHTTPClientV3(nil)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewHTTPClientV3_NoAddress(t *testing.T) {
	client, err := NewHTTPClientV3(&Options{
		Address: "",
	})
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewHTTPClientV3_NoTLSCertificates(t *testing.T) {
	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
		TLS: TLSOptions{
			// Enable TLS but not provide the certificates.
			Enabled: true,
		},
	})
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewHTTPClientV3_defaults(t *testing.T) {
	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	assert.Equal(t, "localhost:5000", client.GetOptions().Address)
	assert.Equal(t, 2*time.Second, client.GetOptions().HTTP.Timeout)
	assert.Equal(t, uint(3), client.GetOptions().HTTP.Retry.Count)
	assert.Equal(t, 100*time.Millisecond, client.GetOptions().HTTP.Retry.WaitTime)
	assert.Equal(t, 2*time.Second, client.GetOptions().HTTP.Retry.MaxWaitTime)
	assert.Empty(t, client.GetOptions().TLS.CertFile)
	assert.Empty(t, client.GetOptions().TLS.KeyFile)
	assert.False(t, client.GetOptions().TLS.Enabled)
	assert.False(t, client.GetOptions().TLS.SkipVerify)
}

func TestNewHTTPClientV3_ValidAddress(t *testing.T) {
	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestNewHTTPClientV3_ValidAddressAndTimeout(t *testing.T) {
	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
		HTTP: HTTPOptions{
			Timeout: 3 * time.Second,
		},
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestNewHTTPClientV3_ValidRetry(t *testing.T) {
	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
		HTTP: HTTPOptions{
			Retry: RetryOptions{
				Count:       3,
				WaitTime:    5 * time.Second,
				MaxWaitTime: 20 * time.Second,
			},
		},
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestHTTPClientV3_Status_200(t *testing.T) {
	in := `
{
  "status":"ok",
  "timestamp":"2019-03-20T17:37:07Z"
}`

	expected := &scheme.Status{
		Status:    "ok",
		Timestamp: "2019-03-20T17:37:07Z",
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeUnversioned(t, "/test", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Status()
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Status_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeUnversioned(t, "/test", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Status()
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Version_200(t *testing.T) {
	in := `
{
  "version":"3.0.0",
  "api_version":"v3"
}`

	expected := &scheme.Version{
		Version:    "3.0.0",
		APIVersion: "v3",
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeUnversioned(t, "/version", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Version()
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Version_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeUnversioned(t, "/version", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Version()
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Config_200(t *testing.T) {
	in := `
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
}`

	expected := &scheme.Config{
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
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/config", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Config()
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Config_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/config", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Config()
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Plugins_200(t *testing.T) {
	in := `
[
  {
    "description": "a plugin with emulated devices and data",
    "id": "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
    "name": "emulator plugin",
    "maintainer": "vapor io",
    "active": true
  },
  {
    "description": "a custom third party plugin",
    "id": "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
    "name": "custom-plugin",
    "maintainer": "third-party",
    "active": true
  }
]`

	expected := []*scheme.PluginMeta{
		{
			Description: "a plugin with emulated devices and data",
			ID:          "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			Name:        "emulator plugin",
			Maintainer:  "vapor io",
			Active:      true,
		},
		{
			Description: "a custom third party plugin",
			ID:          "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
			Name:        "custom-plugin",
			Maintainer:  "third-party",
			Active:      true,
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/plugin", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Plugins()
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Plugins_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/plugin", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Plugins()
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Plugin_200(t *testing.T) {
	in := `
{
   "active":true,
   "id":"1b714cf2-cc56-5c36-9741-fd6a483b5f10",
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
      "timestamp":"2019-03-20T17:37:07Z",
      "status":"ok",
      "message":"",
      "checks":[
         {
            "name":"read buffer health",
            "status":"ok",
            "message":"",
            "timestamp":"2019-03-20T17:37:07Z",
            "type":"periodic"
         },
         {
            "name":"write buffer health",
            "status":"ok",
            "message":"",
            "timestamp":"2019-03-20T17:37:07Z",
            "type":"periodic"
         }
      ]
   }
}`

	expected := &scheme.Plugin{
		PluginMeta: scheme.PluginMeta{
			Active:      true,
			ID:          "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
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
			Timestamp: "2019-03-20T17:37:07Z",
			Status:    "ok",
			Message:   "",
			Checks: []scheme.CheckOptions{
				{
					Name:      "read buffer health",
					Status:    "ok",
					Message:   "",
					Timestamp: "2019-03-20T17:37:07Z",
					Type:      "periodic",
				},
				{
					Name:      "write buffer health",
					Status:    "ok",
					Message:   "",
					Timestamp: "2019-03-20T17:37:07Z",
					Type:      "periodic",
				},
			},
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/plugin/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Plugin("1b714cf2-cc56-5c36-9741-fd6a483b5f10")
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Plugin_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/plugin/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Plugin("1b714cf2-cc56-5c36-9741-fd6a483b5f10")
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_PluginHealth_200(t *testing.T) {
	in := `
{
  "status": "healthy",
  "updated": "2018-06-15T20:04:33Z",
  "healthy": [
    "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
    "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
    "1b714cf2-cc56-5c36-9741-fd6a483b5f12"
  ],
  "unhealthy": [],
  "active": 3,
  "inactive": 0
}`

	expected := &scheme.PluginHealth{
		Status:  "healthy",
		Updated: "2018-06-15T20:04:33Z",
		Healthy: []string{
			"1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			"1b714cf2-cc56-5c36-9741-fd6a483b5f11",
			"1b714cf2-cc56-5c36-9741-fd6a483b5f12",
		},
		Unhealthy: []string{},
		Active:    int(3),
		Inactive:  int(0),
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/plugin/health", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.PluginHealth()
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_PluginHealth_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/plugin/health", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.PluginHealth()
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Scan_200(t *testing.T) {
	in := `
[
  {
    "id": "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
    "alias": "",
    "info": "Synse Temperature Sensor",
    "type": "temperature",
    "plugin": "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
    "tags": [
      "type:temperature",
      "temperature",
      "vio/fan-sensor"
    ]
  },
  {
    "id": "1b714cf2-cc56-5c36-9741-fd6a483b5f12",
    "alias": "",
    "info": "Synse LED",
    "type": "led",
    "plugin": "1b714cf2-cc56-5c36-9741-fd6a483b5f13",
    "tags": [
      "type:led",
      "led"
    ]
  }
]`

	expected := []*scheme.Scan{
		{
			ID:     "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			Alias:  "",
			Info:   "Synse Temperature Sensor",
			Type:   "temperature",
			Plugin: "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
			Tags: []string{
				"type:temperature",
				"temperature",
				"vio/fan-sensor",
			},
		},
		{
			ID:     "1b714cf2-cc56-5c36-9741-fd6a483b5f12",
			Alias:  "",
			Info:   "Synse LED",
			Type:   "led",
			Plugin: "1b714cf2-cc56-5c36-9741-fd6a483b5f13",
			Tags: []string{
				"type:led",
				"led",
			},
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/scan", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ScanOptions{}
	resp, err := client.Scan(opts)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Scan_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/scan", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ScanOptions{}
	resp, err := client.Scan(opts)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Tags_200(t *testing.T) {
	in := `
[
  "default/tag1",
  "default/type:temperature"
]`

	expected := []string{
		"default/tag1",
		"default/type:temperature",
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/tags", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.TagsOptions{}
	resp, err := client.Tags(opts)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Tags_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/tags", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.TagsOptions{}
	resp, err := client.Tags(opts)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Info_200(t *testing.T) {
	in := `
{
   "timestamp":"2019-03-20T17:37:07Z",
   "id":"1b714cf2-cc56-5c36-9741-fd6a483b5f10",
   "alias":"",
   "type":"humidity",
   "metadata":{
      "model":"emul8-humidity"
   },
   "plugin":"1b714cf2-cc56-5c36-9741-fd6a483b5f11",
   "info":"Synse Humidity Sensor",
   "tags":[
      "type:humidity",
      "humidity",
      "vio/fan-sensor"
   ],
   "capabilities":{
      "mode":"rw",
      "read":{

      },
      "write":{
         "actions":[
            "color",
            "state"
         ]
      }
   },
   "outputs":[
      {
         "name":"humidity",
         "type":"humidity",
         "precision":3,
         "scaling_factor":1.0,
         "unit":{
            "name":"percent humidity",
            "symbol":"%"
         }
      },
      {
         "name":"temperature",
         "type":"temperature",
         "precision":3,
         "scaling_factor":1.0,
         "unit":{
            "name":"celsius",
            "symbol":"C"
         }
      }
   ]
}`

	expected := &scheme.Info{
		Timestamp: "2019-03-20T17:37:07Z",
		ID:        "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
		Alias:     "",
		Type:      "humidity",
		Metadata: map[string]string{
			"model": "emul8-humidity",
		},
		Plugin: "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
		Info:   "Synse Humidity Sensor",
		Tags: []string{
			"type:humidity",
			"humidity",
			"vio/fan-sensor",
		},
		Capabilities: scheme.CapabilitiesOptions{
			Mode: "rw",
			Read: map[string]string{},
			Write: scheme.WriteOptions{
				Actions: []string{
					"color",
					"state",
				},
			},
		},
		Outputs: []scheme.OutputOptions{
			{
				Name:          "humidity",
				Type:          "humidity",
				Precision:     int(3),
				ScalingFactor: float64(1.0),
				Unit: scheme.UnitOptions{
					Name:   "percent humidity",
					Symbol: "%",
				},
			},
			{
				Name:          "temperature",
				Type:          "temperature",
				Precision:     int(3),
				ScalingFactor: float64(1.0),
				Unit: scheme.UnitOptions{
					Name:   "celsius",
					Symbol: "C",
				},
			},
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/info/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Info("1b714cf2-cc56-5c36-9741-fd6a483b5f10")
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Info_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/info/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Info("1b714cf2-cc56-5c36-9741-fd6a483b5f10")
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Read_200(t *testing.T) {
	in := `
[
   {
      "device":"1b714cf2-cc56-5c36-9741-fd6a483b5f10",
      "device_type":"temperature",
      "type":"temperature",
      "value":20.3,
      "timestamp":"2019-03-20T17:37:07Z",
      "unit":{
         "symbol":"C",
         "name":"degrees celsius"
      },
      "context":{
         "host":"127.0.0.1",
         "sample_rate":8
      }
   },
   {
      "device":"1b714cf2-cc56-5c36-9741-fd6a483b5f11",
      "device_type":"led",
      "type":"state",
      "value":"off",
      "timestamp":"2019-03-20T17:37:07Z",
      "unit":null
   },
   {
      "device":"1b714cf2-cc56-5c36-9741-fd6a483b5f12",
      "device_type":"led",
      "type":"color",
      "value":"000000",
      "timestamp":"2019-03-20T17:37:07Z",
      "unit":null
   },
   {
      "device":"1b714cf2-cc56-5c36-9741-fd6a483b5f13",
      "device_type":"door_lock",
      "type":"status",
      "value":"locked",
      "timestamp":"2019-03-20T17:37:07Z",
      "unit":null,
      "context":{
         "wedge":1,
         "zone":"6B"
      }
   }
]`

	expected := []*scheme.Read{
		{
			Device:     "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			DeviceType: "temperature",
			Type:       "temperature",
			Value:      float64(20.3),
			Timestamp:  "2019-03-20T17:37:07Z",
			Unit: scheme.UnitOptions{
				Symbol: "C",
				Name:   "degrees celsius",
			},
			Context: map[string]interface{}{
				"host":        "127.0.0.1",
				"sample_rate": float64(8),
			},
		},
		{
			Device:     "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
			DeviceType: "led",
			Type:       "state",
			Value:      "off",
			Timestamp:  "2019-03-20T17:37:07Z",
			Unit:       scheme.UnitOptions{},
		},
		{
			Device:     "1b714cf2-cc56-5c36-9741-fd6a483b5f12",
			DeviceType: "led",
			Type:       "color",
			Value:      "000000",
			Timestamp:  "2019-03-20T17:37:07Z",
			Unit:       scheme.UnitOptions{},
		},
		{
			Device:     "1b714cf2-cc56-5c36-9741-fd6a483b5f13",
			DeviceType: "door_lock",
			Type:       "status",
			Value:      "locked",
			Timestamp:  "2019-03-20T17:37:07Z",
			Unit:       scheme.UnitOptions{},
			Context: map[string]interface{}{
				"wedge": float64(1),
				"zone":  "6B",
			},
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/read", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadOptions{}
	resp, err := client.Read(opts)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Read_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/read", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadOptions{}
	resp, err := client.Read(opts)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_ReadDevice_200(t *testing.T) {
	in := `
[
  {
    "device": "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
    "device_type": "temperature",
    "type": "temperature",
    "value": 20.3,
    "timestamp": "2019-03-20T17:37:07Z",
    "unit": {
      "symbol": "C",
      "name": "degrees celsius"
    },
    "context": {
      "host": "127.0.0.1",
      "sample_rate": 8
    }
  }
]`

	expected := []*scheme.Read{
		{
			Device:     "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			DeviceType: "temperature",
			Type:       "temperature",
			Value:      float64(20.3),
			Timestamp:  "2019-03-20T17:37:07Z",
			Unit: scheme.UnitOptions{
				Symbol: "C",
				Name:   "degrees celsius",
			},
			Context: map[string]interface{}{
				"host":        "127.0.0.1",
				"sample_rate": float64(8),
			},
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/read/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadOptions{}
	resp, err := client.ReadDevice("1b714cf2-cc56-5c36-9741-fd6a483b5f10", opts)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_ReadDevice_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/read/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadOptions{}
	resp, err := client.ReadDevice("1b714cf2-cc56-5c36-9741-fd6a483b5f10", opts)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_ReadCache_200(t *testing.T) {
	in := `
{
  "device":"1b714cf2-cc56-5c36-9741-fd6a483b5f10",
  "device_type":"led",
  "type":"state",
  "value":"off",
  "timestamp":"2019-03-20T17:37:07Z",
  "unit":null
}
{
  "device":"1b714cf2-cc56-5c36-9741-fd6a483b5f11",
  "device_type":"led",
  "type":"color",
  "value":"000000",
  "timestamp":"2019-03-20T17:37:07Z",
  "unit":null
}`

	expected := []*scheme.Read{
		{
			Device:     "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			DeviceType: "led",
			Type:       "state",
			Value:      "off",
			Timestamp:  "2019-03-20T17:37:07Z",
			Unit:       scheme.UnitOptions{},
		},
		{
			Device:     "1b714cf2-cc56-5c36-9741-fd6a483b5f11",
			DeviceType: "led",
			Type:       "color",
			Value:      "000000",
			Timestamp:  "2019-03-20T17:37:07Z",
			Unit:       scheme.UnitOptions{},
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/readcache", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadCacheOptions{}
	readings := make(chan *scheme.Read, 1)

	go func() {
		err := client.ReadCache(opts, readings)
		assert.NoError(t, err)
	}()

	var results []*scheme.Read

	for {
		var done bool
		select {
		case r, open := <-readings:
			if !open {
				done = true
				break
			}
			results = append(results, r)

		case <-time.After(2 * time.Second):
			// If the test does not complete after 2s, error.
			t.Fatal("timeout: failed getting readcache data from channel")
		}

		if done {
			break
		}
	}

	assert.Equal(t, expected, results)
}

func TestHTTPClientV3_ReadCache_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/readcache", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadCacheOptions{}
	readings := make(chan *scheme.Read, 1)

	go func() {
		err := client.ReadCache(opts, readings)
		assert.Error(t, err)
	}()

	var results []*scheme.Read

	for {
		var done bool
		select {
		case r, open := <-readings:
			if !open {
				done = true
				break
			}
			results = append(results, r)

		case <-time.After(2 * time.Second):
			// If the test does not complete after 2s, error.
			t.Fatal("timeout: failed getting readcache data from channel")
		}

		if done {
			break
		}
	}
	assert.Empty(t, results)
}

func TestHTTPClientV3_WriteAsync_200(t *testing.T) {
	in := `
[
  {
    "id": "56a32eba-1aa6-4868-84ee-fe01af8b2e6d",
    "device": "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
    "context": {
      "action": "color",
      "data": "f38ac2"
    },
    "timeout": "10s"
  },
  {
    "id": "56a32eba-1aa6-4868-84ee-fe01af8b2e6e",
    "device": "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
    "context": {
      "action": "state",
      "data": "blink"
    },
    "timeout": "10s"
  }
]`

	expected := []*scheme.Write{
		{
			ID:     "56a32eba-1aa6-4868-84ee-fe01af8b2e6d",
			Device: "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			Context: scheme.WriteData{
				Action: "color",
				Data:   "f38ac2",
			},
			Timeout: "10s",
		},
		{
			ID:     "56a32eba-1aa6-4868-84ee-fe01af8b2e6e",
			Device: "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			Context: scheme.WriteData{
				Action: "state",
				Data:   "blink",
			},
			Timeout: "10s",
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/write/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	var opts []scheme.WriteData
	resp, err := client.WriteAsync("1b714cf2-cc56-5c36-9741-fd6a483b5f10", opts)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_WriteAsync_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/write/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	var opts []scheme.WriteData
	resp, err := client.WriteAsync("1b714cf2-cc56-5c36-9741-fd6a483b5f10", opts)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_WriteSync_200(t *testing.T) {
	in := `
[
 {
   "id": "56a32eba-1aa6-4868-84ee-fe01af8b2e6d",
   "timeout": "10s",
   "device": "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
   "context": {
     "action": "color",
     "data": "f38ac2"
   },
   "status": "done",
   "created": "2018-02-01T15:00:51Z",
   "updated": "2018-02-01T15:00:51Z",
   "message": ""
 }
]`

	expected := []*scheme.Transaction{
		{
			ID:      "56a32eba-1aa6-4868-84ee-fe01af8b2e6d",
			Timeout: "10s",
			Device:  "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
			Context: scheme.WriteData{
				Action: "color",
				Data:   "f38ac2",
			},
			Status:  "done",
			Created: "2018-02-01T15:00:51Z",
			Updated: "2018-02-01T15:00:51Z",
			Message: "",
		},
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/write/wait/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	var opts []scheme.WriteData
	resp, err := client.WriteSync("1b714cf2-cc56-5c36-9741-fd6a483b5f10", opts)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_WriteSync_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/write/wait/1b714cf2-cc56-5c36-9741-fd6a483b5f10", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	var opts []scheme.WriteData
	resp, err := client.WriteSync("1b714cf2-cc56-5c36-9741-fd6a483b5f10", opts)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Transactions_200(t *testing.T) {
	in := `
[
  "56a32eba-1aa6-4868-84ee-fe01af8b2e6b",
  "56a32eba-1aa6-4868-84ee-fe01af8b2e6c",
  "56a32eba-1aa6-4868-84ee-fe01af8b2e6d"
]`

	expected := []string{
		"56a32eba-1aa6-4868-84ee-fe01af8b2e6b",
		"56a32eba-1aa6-4868-84ee-fe01af8b2e6c",
		"56a32eba-1aa6-4868-84ee-fe01af8b2e6d",
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/transaction", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Transactions()
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Transactions_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/transaction", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Transactions()
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_Transaction_200(t *testing.T) {
	in := `
{
  "id": "56a32eba-1aa6-4868-84ee-fe01af8b2e6b",
  "timeout": "10s",
  "device": "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
  "context": {
    "action": "color",
    "data": "f38ac2"
  },
  "status": "done",
  "created": "2018-02-01T15:00:51Z",
  "updated": "2018-02-01T15:00:51Z",
  "message": ""
}`

	expected := &scheme.Transaction{
		ID:      "56a32eba-1aa6-4868-84ee-fe01af8b2e6b",
		Timeout: "10s",
		Device:  "1b714cf2-cc56-5c36-9741-fd6a483b5f10",
		Context: scheme.WriteData{
			Action: "color",
			Data:   "f38ac2",
		},
		Status:  "done",
		Created: "2018-02-01T15:00:51Z",
		Updated: "2018-02-01T15:00:51Z",
		Message: "",
	}

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/transaction/56a32eba-1aa6-4868-84ee-fe01af8b2e6b", 200, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Transaction("56a32eba-1aa6-4868-84ee-fe01af8b2e6b")
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_Transaction_500(t *testing.T) {
	in := `
{
  "http_code":500,
  "description":"unknown error",
  "timestamp":"2019-03-20T17:37:07Z",
  "context":"unknown error"
}`

	server := test.NewHTTPServerV3()
	defer server.Close()

	server.ServeVersioned(t, "/transaction/56a32eba-1aa6-4868-84ee-fe01af8b2e6b", 500, in)

	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Transaction("56a32eba-1aa6-4868-84ee-fe01af8b2e6b")
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestHTTPClientV3_TLS(t *testing.T) {
	// certFile and keyFile are self-signed test certificates' locations.
	certFile, keyFile := "testdata/cert.pem", "testdata/key.pem"

	// Parse the certificates.
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	assert.NotNil(t, cert)
	assert.NoError(t, err)

	// Create a mock https server and let it use the certificates.
	server := test.NewHTTPSServerV3()
	defer server.Close()

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	server.SetTLS(cfg)
	assert.NotNil(t, server.GetCertificates())

	// Setup a client that also uses the certificates.
	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
		TLS: TLSOptions{
			Enabled:    true,
			CertFile:   certFile,
			KeyFile:    keyFile,
			SkipVerify: true, // skip CA known authority check
		},
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	// Only need to setup one arbitrary unversioned endpoint and another
	// versioned one to make requests against since we already have tests for
	// all the endpoints and the works there are pretty much same.
	in := `
{
  "status":"ok",
  "timestamp":"2019-03-20T17:37:07Z"
}`

	expected := &scheme.Status{
		Status:    "ok",
		Timestamp: "2019-03-20T17:37:07Z",
	}

	server.ServeUnversioned(t, "/test", 200, in)

	resp, err := client.Status()
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestHTTPClientV3_TLS_UnknownCA(t *testing.T) {
	// certFile and keyFile are self-signed test certificates' locations.
	certFile, keyFile := "testdata/cert.pem", "testdata/key.pem"

	// Parse the certificates.
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	assert.NotNil(t, cert)
	assert.NoError(t, err)

	// Create a mock https server and let it use the certificates.
	server := test.NewHTTPSServerV3()
	defer server.Close()

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	server.SetTLS(cfg)
	assert.NotNil(t, server.GetCertificates())

	// Setup a client that also uses the certificates. However, this time we
	// don't skip the CA known security check.
	client, err := NewHTTPClientV3(&Options{
		Address: server.URL,
		TLS: TLSOptions{
			Enabled:  true,
			CertFile: certFile,
			KeyFile:  keyFile,
		},
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	// Only need to setup one arbitrary unversioned endpoint and another
	// versioned one to make requests against since we already have tests for
	// all the endpoints and the works there are pretty much same.
	in := `
{
  "status":"ok",
  "timestamp":"2019-03-20T17:37:07Z"
}`

	server.ServeUnversioned(t, "/test", 200, in)

	resp, err := client.Status()
	assert.Nil(t, resp)
	assert.Error(t, err)
}
