package synse

import (
	// "encoding/json"
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
	assert.NotEmpty(t, resp.Timestamp)
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
	assert.NotEmpty(t, resp.Version)
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
	assert.NotEmpty(t, resp)
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
	assert.NotEmpty(t, resp[0])
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
	assert.NotEmpty(t, resp)
}

func TestIntegration_PluginHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.PluginHealth()
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestIntegration_Device(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ScanOptions{}
	resp, err := client.Scan(opts)
	assert.NoError(t, err)
	assert.Equal(t, 22, len(resp))

	for _, v := range resp {
		assert.NotEmpty(t, v)

		info, err := client.Info(v.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, info)
	}
}

func TestIntegration_Tags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.TagsOptions{}
	resp, err := client.Tags(opts)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(resp))

	for _, v := range resp {
		assert.NotEmpty(t, v)
	}
}
