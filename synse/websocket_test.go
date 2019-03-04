package synse

import (
	// "crypto/tls"
	"testing"
	"time"

	"github.com/vapor-ware/synse-client-go/internal/test"
	"github.com/vapor-ware/synse-client-go/synse/scheme"

	"github.com/stretchr/testify/assert"
)

func TestNewWebSocketClientV3_NilConfig(t *testing.T) {
	client, err := NewWebSocketClientV3(nil)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewWebSocketClientV3_NoAddress(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "",
	})
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewWebSocketClientV3_NoTLSCertificates(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
		TLS: TLSOptions{
			// Enable TLS but not provide the certificates.
			Enabled: true,
		},
	})
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewWebSocketClientV3_defaults(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	assert.Equal(t, "localhost:5000", client.GetOptions().Address)
	assert.Equal(t, 45*time.Second, client.GetOptions().WebSocket.HandshakeTimeout)
	assert.Empty(t, client.GetOptions().TLS.CertFile)
	assert.Empty(t, client.GetOptions().TLS.KeyFile)
	assert.False(t, client.GetOptions().TLS.SkipVerify)
}

func TestNewWebSocketClientV3_ValidAddress(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestNewWebSocketClientV3_ValidAddressAndTimeout(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
		WebSocket: WebSocketOptions{
			HandshakeTimeout: 46 * time.Second,
		},
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestWebSocketClientV3_Version_200(t *testing.T) {
	expected := &scheme.Version{
		Version:    "3.0.0",
		APIVersion: "v3",
	}

	s := test.NewWebSocketServerV3()
	defer s.Close()

	s.Serve(expected)

	client, err := NewWebSocketClientV3(&Options{
		Address: s.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	v, err := client.Version()
	assert.NotNil(t, v)
	assert.NoError(t, err)
	assert.Equal(t, expected, v)
}

func TestWebSocketClientV3_Plugin_200(t *testing.T) {
	expected := &scheme.Plugin{
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
	}

	s := test.NewWebSocketServerV3()
	defer s.Close()

	s.Serve(expected)

	client, err := NewWebSocketClientV3(&Options{
		Address: s.URL,
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	v, err := client.Plugin("12835beffd3e6c603aa4dd92127707b5")
	assert.NotNil(t, v)
	assert.NoError(t, err)
	assert.Equal(t, expected, v)
}
