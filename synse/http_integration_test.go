package synse

import (
	// "encoding/json"
	// "fmt"
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
)

func TestIntegration_Status(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Status()
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
	assert.NotNil(t, resp.Timestamp)
}

func TestIntegration_Version(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Version()
	assert.NoError(t, err)
	assert.Equal(t, "v3", resp.APIVersion)
	assert.NotNil(t, resp.Version)
}

func TestIntegration_Config(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Config()
	assert.NoError(t, err)

	expected := &scheme.Config{
		Locale:  "en_US",
		Logging: "debug",
		Plugin: scheme.PluginOptions{
			TCP:  []string{"emulator:5001"},
			Unix: []string{},
			Discover: scheme.DiscoveryOptions{
				Kubernetes: scheme.KubernetesOptions{
					Namespace: "",
					Endpoints: scheme.EndpointsOptions{
						Labels: map[string]string(nil),
					},
				},
			},
		},
		Cache: scheme.CacheOptions{
			Device:      scheme.DeviceOptions{TTL: 0},
			Transaction: scheme.TransactionOptions{TTL: 300},
		},
		GRPC: scheme.GRPCOptions{
			Timeout: 3,
			TLS:     scheme.TLSOptions{Cert: ""}},
		Transport: scheme.TransportOptions{
			HTTP:      false, // FIXME - should this be true?
			WebSocket: false,
		},
		Metrics:    scheme.MetricsOptions{Enabled: false},
		PrettyJSON: true,
	}
	assert.Equal(t, expected, resp)
}

func TestIntegration_Plugin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Plugins()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp))

	expected := &scheme.PluginMeta{
		Active:      true,
		ID:          "4032ffbe-80db-5aa5-b794-f35c88dff85c",
		Name:        "emulator plugin",
		Description: "A plugin with emulated devices and data",
		Maintainer:  "vaporio",
		Tag:         "vaporio/emulator-plugin",
		VCS:         "",
		Version: scheme.VersionOptions{
			PluginVersion: "",
			SDKVersion:    "",
			BuildDate:     "",
			GitCommit:     "",
			GitTag:        "",
			Arch:          "",
			OS:            "",
		},
	}
	assert.Equal(t, expected, resp[0])
}
