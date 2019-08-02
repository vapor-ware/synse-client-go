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

	status, err := client.Status()
	assert.NoError(t, err)
	assert.Equal(t, "ok", status.Status)
	assert.NotEmpty(t, status.Timestamp)
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

	version, err := client.Version()
	assert.NoError(t, err)
	assert.Equal(t, "v3", version.APIVersion)
	assert.NotEmpty(t, version.Version)
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

	config, err := client.Config()
	assert.NoError(t, err)
	assert.NotEmpty(t, config)
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

	plugins, err := client.Plugins()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(plugins))
	assert.NotEmpty(t, plugins[0])

	plugin, err := client.Plugin(plugins[0].ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, plugin)

	pluginHealth, err := client.PluginHealth()
	assert.NoError(t, err)
	assert.NotEmpty(t, pluginHealth)
}

func TestIntegration_Read(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ScanOptions{}
	devices, err := client.Scan(opts)
	assert.NoError(t, err)
	assert.Equal(t, 22, len(devices))

	for _, device := range devices {
		assert.NotEmpty(t, device)

		info, err := client.Info(device.ID)
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
	tags, err := client.Tags(opts)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(tags))

	for _, tag := range tags {
		assert.NotEmpty(t, tag)
	}
}
