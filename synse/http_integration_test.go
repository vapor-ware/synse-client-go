package synse

import (
	"testing"
	"time"

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
	assert.NotEmpty(t, version.Version)
	assert.Equal(t, "v3", version.APIVersion)
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
	assert.Equal(t, "debug", config.Logging)
	assert.True(t, config.PrettyJSON)
	assert.Equal(t, "en_US", config.Locale)
	assert.Equal(t, 1, len(config.Plugin.TCP))
	assert.Equal(t, "emulator:5001", config.Plugin.TCP[0])
	assert.Equal(t, 0, len(config.Plugin.Unix))
	// assert.Equal(t, 180, config.Cache.Device.RebuildEvery) // TODO - update scheme
	assert.Equal(t, 300, config.Cache.Transaction.TTL)
	assert.Equal(t, 3, config.GRPC.Timeout)
	assert.False(t, config.Metrics.Enabled)
}

func TestIntegration_Plugins(t *testing.T) {
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

	plugin := plugins[0]
	assert.Equal(t, "emulator plugin", plugin.Name)
	assert.Equal(t, "vaporio", plugin.Maintainer)
	assert.Equal(t, "vaporio/emulator-plugin", plugin.Tag)
	assert.Equal(t, "A plugin with emulated devices and data", plugin.Description)
	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", plugin.ID)
	assert.True(t, plugin.Active)
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

	plugin, err := client.Plugin("4032ffbe-80db-5aa5-b794-f35c88dff85c")
	assert.NoError(t, err)
	assert.Equal(t, "emulator plugin", plugin.Name)
	assert.Equal(t, "vaporio", plugin.Maintainer)
	assert.Equal(t, "vaporio/emulator-plugin", plugin.Tag)
	assert.Equal(t, "A plugin with emulated devices and data", plugin.Description)
	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", plugin.ID)
	assert.True(t, plugin.Active)

	assert.Equal(t, "emulator:5001", plugin.Network.Address)
	assert.Equal(t, "tcp", plugin.Network.Protocol)
	assert.NotEmpty(t, plugin.Version.PluginVersion)
	assert.NotEmpty(t, plugin.Version.SDKVersion)
	assert.NotEmpty(t, plugin.Version.BuildDate)
	assert.NotEmpty(t, plugin.Version.GitCommit)
	assert.NotEmpty(t, plugin.Version.GitTag)
	assert.Equal(t, "amd64", plugin.Version.Arch)
	assert.Equal(t, "linux", plugin.Version.OS)
	assert.NotEmpty(t, plugin.Health.Timestamp)
	assert.Equal(t, "OK", plugin.Health.Status)
	assert.Equal(t, 2, len(plugin.Health.Checks))

	readCheck := plugin.Health.Checks[0]
	assert.Equal(t, "read queue health", readCheck.Name)
	assert.Equal(t, "OK", readCheck.Status)
	assert.Equal(t, "periodic", readCheck.Type)
	assert.NotEmpty(t, readCheck.Timestamp)
	assert.Empty(t, readCheck.Message)

	writeCheck := plugin.Health.Checks[1]
	assert.Equal(t, "write queue health", writeCheck.Name)
	assert.Equal(t, "OK", writeCheck.Status)
	assert.Equal(t, "periodic", writeCheck.Type)
	assert.NotEmpty(t, writeCheck.Timestamp)
	assert.Empty(t, writeCheck.Message)
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

	health, err := client.PluginHealth()
	assert.NoError(t, err)
	assert.Equal(t, "healthy", health.Status)
	assert.NotEmpty(t, health.Updated)
	assert.Equal(t, 1, len(health.Healthy))
	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", health.Healthy[0])
	assert.Equal(t, 0, len(health.Unhealthy))
	assert.Equal(t, 1, health.Active)
	assert.Equal(t, 0, health.Inactive)
}

func TestIntegration_Scan(t *testing.T) {
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
	assert.Equal(t, 1, len(devices))

	device := devices[0]
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", device.ID)
	assert.Equal(t, "emulator-led", device.Alias)
	assert.Equal(t, "Synse LED", device.Info)
	assert.Equal(t, "led", device.Type)
	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", device.Plugin)
	assert.Equal(t, 2, len(device.Tags))
	assert.Equal(t, "system/id:f041883c-cf87-55d7-a978-3d3103836412", device.Tags[0])
	assert.Equal(t, "system/type:led", device.Tags[1])
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
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, "system/type:led", tags[0])
}

func TestIntegration_Info(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	device, err := client.Info("f041883c-cf87-55d7-a978-3d3103836412")
	assert.NoError(t, err)
	assert.NotEmpty(t, device.Timestamp)
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", device.ID)
	assert.Equal(t, "led", device.Type)
	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", device.Plugin)
	assert.Equal(t, "Synse LED", device.Info)
	assert.Equal(t, "emulator-led", device.Alias)
	assert.Equal(t, map[string]string{"model": "emul8-led"}, device.Metadata)
	assert.Equal(t, "rw", device.Capabilities.Mode)
	assert.Equal(t, 0, len(device.Capabilities.Write.Actions))
	assert.Equal(t, 2, len(device.Tags))
	assert.Equal(t, "system/id:f041883c-cf87-55d7-a978-3d3103836412", device.Tags[0])
	assert.Equal(t, "system/type:led", device.Tags[1])
	assert.Equal(t, 2, len(device.Outputs))

	stateOutput := device.Outputs[0]
	assert.Equal(t, "state", stateOutput.Name)
	assert.Equal(t, "state", stateOutput.Type)
	assert.Equal(t, 0, stateOutput.Precision)
	assert.Equal(t, 0.0, stateOutput.ScalingFactor)

	colorOutput := device.Outputs[1]
	assert.Equal(t, "color", colorOutput.Name)
	assert.Equal(t, "color", colorOutput.Type)
	assert.Equal(t, 0, colorOutput.Precision)
	assert.Equal(t, 0.0, colorOutput.ScalingFactor)

	// assert.Equal(t, 0, device.SortIndex) // TODO - update scheme
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

	opts := scheme.ReadOptions{}
	readings, err := client.Read(opts)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(readings))

	stateRead := readings[0]
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", stateRead.Device)
	assert.NotEmpty(t, stateRead.Timestamp)
	assert.Equal(t, "state", stateRead.Type)
	assert.Equal(t, "led", stateRead.DeviceType)
	assert.Empty(t, stateRead.Unit)
	assert.Equal(t, "off", stateRead.Value)
	assert.Empty(t, stateRead.Context)

	colorRead := readings[1]
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", colorRead.Device)
	assert.NotEmpty(t, colorRead.Timestamp)
	assert.Equal(t, "color", colorRead.Type)
	assert.Equal(t, "led", colorRead.DeviceType)
	assert.Empty(t, colorRead.Unit)
	assert.Equal(t, "000000", colorRead.Value)
	assert.Empty(t, colorRead.Context)
}

func TestIntegration_ReadDevice(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadOptions{}
	readings, err := client.ReadDevice("f041883c-cf87-55d7-a978-3d3103836412", opts)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(readings))

	stateRead := readings[0]
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", stateRead.Device)
	assert.NotEmpty(t, stateRead.Timestamp)
	assert.Equal(t, "state", stateRead.Type)
	assert.Equal(t, "led", stateRead.DeviceType)
	assert.Empty(t, stateRead.Unit)
	assert.Equal(t, "off", stateRead.Value)
	assert.Empty(t, stateRead.Context)

	colorRead := readings[1]
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", colorRead.Device)
	assert.NotEmpty(t, colorRead.Timestamp)
	assert.Equal(t, "color", colorRead.Type)
	assert.Equal(t, "led", colorRead.DeviceType)
	assert.Empty(t, colorRead.Unit)
	assert.Equal(t, "000000", colorRead.Value)
	assert.Empty(t, colorRead.Context)
}

func TestIntegration_ReadCache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	opts := scheme.ReadCacheOptions{}
	readings := make(chan *scheme.Read, 1)

	go func() {
		err := client.ReadCache(opts, readings)
		assert.NoError(t, err)
	}()

	for {
		var done bool
		select {
		case read, open := <-readings:
			if !open {
				done = true
				break
			}

			assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", read.Device)
			assert.NotEmpty(t, read.Timestamp)
			assert.Contains(t, []string{"state", "color"}, read.Type)
			assert.Equal(t, "led", read.DeviceType)
			assert.Empty(t, read.Unit)
			assert.Contains(t, []string{"off", "000000"}, read.Value)

		case <-time.After(2 * time.Second):
			// if the test does not complete after 2s, error.
			t.Fatal("timeout: failed getting readcache data from channel")
		}

		if done {
			break
		}
	}
}

func TestIntegration_WriteAsync(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	writeData := []scheme.WriteData{
		{Action: "state", Data: "on"},
		{Action: "color", Data: "ffffff"},
	}
	writes, err := client.WriteAsync("f041883c-cf87-55d7-a978-3d3103836412", writeData)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(writes))

	stateWrite := writes[0]
	assert.NotEmpty(t, stateWrite.ID)
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", stateWrite.Device)
	assert.Equal(t, "state", stateWrite.Context.Action)
	// assert.Equal(t, "on", stateWrite.Context.Data) // FIXME - reflected data isn't decoded yet
	assert.Empty(t, stateWrite.Context.Transaction)
	assert.Equal(t, "30s", stateWrite.Timeout)

	colorWrite := writes[1]
	assert.NotEmpty(t, colorWrite.ID)
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", colorWrite.Device)
	assert.Equal(t, "color", colorWrite.Context.Action)
	// assert.Equal(t, "ffffff", colorWrite.Context.Data) // FIXME - reflected data isn't decoded yet
	assert.Empty(t, colorWrite.Context.Transaction)
	assert.Equal(t, "30s", colorWrite.Timeout)
}

func TestIntegration_WriteSync(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	writeData := []scheme.WriteData{
		{Action: "state", Data: "blink"},
		{Action: "color", Data: "0f0f0f"},
	}
	writes, err := client.WriteSync("f041883c-cf87-55d7-a978-3d3103836412", writeData)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(writes))

	stateWrite := writes[0]
	assert.NotEmpty(t, stateWrite.ID)
	assert.NotEmpty(t, stateWrite.Created)
	assert.NotEmpty(t, stateWrite.Updated)
	assert.Equal(t, "30s", stateWrite.Timeout)
	assert.Equal(t, "DONE", stateWrite.Status)
	assert.Equal(t, "state", stateWrite.Context.Action)
	// assert.Equal(t, "blink", stateWrite.Context.Data) // FIXME - reflected data isn't decoded yet
	assert.Empty(t, stateWrite.Context.Transaction)
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", stateWrite.Device)

	colorWrite := writes[1]
	assert.NotEmpty(t, colorWrite.ID)
	assert.NotEmpty(t, colorWrite.Created)
	assert.NotEmpty(t, colorWrite.Updated)
	assert.Equal(t, "30s", colorWrite.Timeout)
	assert.Equal(t, "DONE", colorWrite.Status)
	assert.Equal(t, "color", colorWrite.Context.Action)
	// assert.Equal(t, "0f0f0f", colorWrite.Context.Data) // FIXME - reflected data isn't decoded yet
	assert.Empty(t, colorWrite.Context.Transaction)
	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", colorWrite.Device)
}

func TestIntegration_Transaction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	transactions, err := client.Transactions()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(transactions))

	for _, id := range transactions {
		transaction, err := client.Transaction(id)
		assert.NoError(t, err)
		assert.NotEmpty(t, transaction.ID)
		assert.NotEmpty(t, transaction.Created)
		assert.NotEmpty(t, transaction.Updated)
		assert.Equal(t, "30s", transaction.Timeout)
		assert.Equal(t, "DONE", transaction.Status)
		assert.Contains(t, []string{"state", "color"}, transaction.Context.Action)
		// assert.Contains(t, []string{"on", "blink", "ffffff", "0f0f0f"}, transaction.Context.Data) // FIXME - reflected data isn't decoded yet
		assert.Empty(t, transaction.Message)
		assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", transaction.Device)
	}
}
