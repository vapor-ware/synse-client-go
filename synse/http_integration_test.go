package synse

import (
	// "encoding/json"
	// "fmt"
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
	// "github.com/vapor-ware/synse-client-go/synse/scheme"
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
	assert.Equal(t, "en_US", resp.Locale)
	assert.Equal(t, "debug", resp.Logging)
	assert.Equal(t, 1, len(resp.Plugin.TCP))
	assert.Equal(t, 0, len(resp.Plugin.Unix))
	assert.Empty(t, resp.Plugin.Discover)
	assert.Equal(t, 0, resp.Cache.Device.TTL)
	assert.Equal(t, 300, resp.Cache.Transaction.TTL)
	assert.Equal(t, 3, resp.GRPC.Timeout)
	assert.Empty(t, resp.GRPC.TLS.Cert)
	// FIXME - should both transport config be True because any http or
	// websocket client could talk to it?
	assert.False(t, resp.Transport.HTTP)
	assert.False(t, resp.Transport.WebSocket)
	assert.False(t, resp.Metrics.Enabled)
	assert.True(t, resp.PrettyJSON)
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
	assert.True(t, resp[0].Active)
	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", resp[0].ID)
	assert.Equal(t, "emulator plugin", resp[0].Name)
	assert.Equal(t, "A plugin with emulated devices and data", resp[0].Description)
	assert.Equal(t, "vaporio", resp[0].Maintainer)
	assert.Equal(t, "vaporio/emulator-plugin", resp[0].Tag)
}

func TestIntegration_PluginInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Plugin("4032ffbe-80db-5aa5-b794-f35c88dff85c")
	assert.NoError(t, err)
	assert.True(t, resp.Active)
	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", resp.ID)
	assert.Equal(t, "emulator plugin", resp.Name)
	assert.Equal(t, "A plugin with emulated devices and data", resp.Description)
	assert.Equal(t, "vaporio", resp.Maintainer)
	assert.Equal(t, "vaporio/emulator-plugin", resp.Tag)
	assert.Equal(t, "github.com/vapor-ware/synse-emulator-plugin", resp.VCS)
	assert.Equal(t, "3.0.0-alpha.3", resp.Version.PluginVersion)
	assert.Equal(t, "3.0.0-alpha.1", resp.Version.SDKVersion)
	assert.NotNil(t, resp.Version.BuildDate)
	assert.NotNil(t, resp.Version.GitCommit)
	assert.Equal(t, "3.0.0-alpha.3", resp.Version.GitTag)
	assert.Equal(t, "amd64", resp.Version.Arch)
	assert.Equal(t, "linux", resp.Version.OS)
	assert.Equal(t, "tcp", resp.Network.Protocol)
	assert.Equal(t, "emulator:5001", resp.Network.Address)
	assert.NotNil(t, resp.Health.Timestamp)
	assert.Equal(t, "OK", resp.Health.Status)
	assert.Empty(t, resp.Health.Message)
	assert.Equal(t, 2, len(resp.Health.Checks))
	assert.Equal(t, "read queue health", resp.Health.Checks[0].Name)
	assert.Equal(t, "OK", resp.Health.Checks[0].Status)
	assert.Empty(t, resp.Health.Checks[0].Message)
	assert.Empty(t, resp.Health.Checks[0].Timestamp) // FIXME - should this be non-empty?
	assert.Equal(t, "periodic", resp.Health.Checks[0].Type)
	assert.Equal(t, "write queue health", resp.Health.Checks[1].Name)
	assert.Equal(t, "OK", resp.Health.Checks[1].Status)
	assert.Empty(t, resp.Health.Checks[1].Message)
	assert.Empty(t, resp.Health.Checks[1].Timestamp) // FIXME - should this be non-empty?
	assert.Equal(t, "periodic", resp.Health.Checks[1].Type)
}
