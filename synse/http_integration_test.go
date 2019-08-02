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
	assert.Equal(t, "emulator:5001", resp.Plugin.TCP[0])
	assert.Empty(t, resp.Plugin.Unix)
	assert.Empty(t, resp.Plugin.Discover)
	assert.Equal(t, 0, resp.Cache.Device.TTL)
	assert.Equal(t, 300, resp.Cache.Transaction.TTL)
	assert.Equal(t, 3, resp.GRPC.Timeout)
	assert.Empty(t, resp.GRPC.TLS.Cert)
	assert.Equal(t, false, resp.Transport.HTTP) // FIXME - should this be true?
	assert.Equal(t, false, resp.Transport.WebSocket)
	assert.Equal(t, false, resp.Metrics.Enabled)
	assert.Equal(t, true, resp.PrettyJSON)
}
