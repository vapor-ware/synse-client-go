package synse

import (
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
	assert.Equal(t, "3.0.0-alpha.3", version.Version)
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
	assert.Equal(t, "3.0.0-alpha.3", plugin.Version.PluginVersion)
	assert.Equal(t, "3.0.0-alpha.1", plugin.Version.SDKVersion)
	assert.NotEmpty(t, plugin.Version.BuildDate)
	assert.Equal(t, "4234777", plugin.Version.GitCommit)
	assert.Equal(t, "3.0.0-alpha.3", plugin.Version.GitTag)
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

// func TestIntegration_ReadCache(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	opts := scheme.ReadCacheOptions{}
// 	reads := make(chan *scheme.Read, 1)

// 	go func() {
// 		err := client.ReadCache(opts, reads)
// 		assert.NoError(t, err)
// 	}()

// 	for {
// 		var done bool
// 		select {
// 		case read, open := <-reads:
// 			if !open {
// 				done = true
// 				break
// 			}

// 			assert.NotEmpty(t, read.Device)
// 			assert.NotEmpty(t, read.Timestamp)
// 			// assert.NotEmpty(t, read.Type) // FIXME - fan and airflow types are empty?
// 			assert.NotEmpty(t, read.DeviceType)

// 			// led and lock devices don't produce unit output
// 			if read.DeviceType == "led" || read.DeviceType == "lock" {
// 				assert.Empty(t, read.Unit.Name)
// 				assert.Empty(t, read.Unit.Symbol)
// 			} else {
// 				assert.NotEmpty(t, read.Unit.Name)
// 				assert.NotEmpty(t, read.Unit.Symbol)
// 			}

// 			// NOTE - read.Value could be 0 so no need to check that

// 		case <-time.After(2 * time.Second):
// 			// if the test does not complete after 2s, error.
// 			t.Fatal("timeout: failed getting readcache data from channel")
// 		}

// 		if done {
// 			break
// 		}
// 	}
// }

// func TestIntegration_WriteAsync(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	// collect all writable devices.
// 	// FIXME - refer to #69, can only query one type atm.
// 	opts := scheme.ScanOptions{
// 		Tags: []string{
// 			"system/type:fan",
// 			// "system/type:led",
// 			// "system/type:lock",
// 			// "system/type:power",
// 		},
// 	}
// 	devices, err := client.Scan(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(devices))

// 	for _, device := range devices {
// 		assert.NotEmpty(t, device.ID)

// 		writeData := []scheme.WriteData{
// 			{Action: "speed", Data: "101"}, // FIXME - if Data is int, get 500 from server.
// 		}
// 		writes, err := client.WriteAsync(device.ID, writeData)
// 		assert.NoError(t, err)

// 		assert.Equal(t, 1, len(writes))

// 		for _, write := range writes {
// 			assert.NotEmpty(t, write.ID)
// 			assert.NotEmpty(t, write.Device)
// 			assert.Equal(t, "speed", write.Context.Action)
// 			// FIXME - reflected data is not decoded yet
// 			// assert.Equal(t, "101"), write.Context.Data)
// 			assert.Empty(t, write.Context.Transaction)
// 			assert.NotEmpty(t, write.Timeout)
// 		}
// 	}
// }

// func TestIntegration_WriteSync(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	// collect all writable devices.
// 	// FIXME - refer to #69, can only query one type atm.
// 	opts := scheme.ScanOptions{
// 		Tags: []string{
// 			"system/type:fan",
// 			// "system/type:led",
// 			// "system/type:lock",
// 			// "system/type:power",
// 		},
// 	}
// 	devices, err := client.Scan(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(devices))

// 	for _, device := range devices {
// 		assert.NotEmpty(t, device.ID)

// 		writeData := []scheme.WriteData{
// 			{Action: "speed", Data: "101"}, // FIXME - if Data is int, get 500 from server.
// 		}
// 		writes, err := client.WriteSync(device.ID, writeData)
// 		assert.NoError(t, err)

// 		assert.Equal(t, 1, len(writes))

// 		for _, write := range writes {
// 			assert.NotEmpty(t, write.ID)
// 			assert.NotEmpty(t, write.Created)
// 			assert.NotEmpty(t, write.Updated)
// 			assert.NotEmpty(t, write.Timeout)
// 			assert.Equal(t, "DONE", write.Status)

// 			assert.Equal(t, "speed", write.Context.Action)
// 			// FIXME - reflected data is not decoded yet
// 			// assert.Equal(t, "101"), write.Context.Data)
// 			assert.Empty(t, write.Context.Transaction)
// 			assert.NotEmpty(t, write.Timeout)
// 			assert.Equal(t, device.ID, write.Device)
// 			assert.Empty(t, write.Message)
// 		}
// 	}
// }

// func TestIntegration_Transactions(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	transactions, err := client.Transactions()
// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(transactions))
// }

// func TestIntegration_Transaction(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	transactions, err := client.Transactions()
// 	assert.NoError(t, err)

// 	for _, id := range transactions {
// 		transaction, err := client.Transaction(id)
// 		assert.NoError(t, err)

// 		assert.NotEmpty(t, transaction.ID)
// 		assert.NotEmpty(t, transaction.Created)
// 		assert.NotEmpty(t, transaction.Updated)
// 		assert.NotEmpty(t, transaction.Timeout)
// 		assert.Equal(t, "DONE", transaction.Status)
// 		assert.NotEmpty(t, transaction.Context.Action)
// 		assert.NotEmpty(t, transaction.Context.Data)
// 		assert.Empty(t, transaction.Context.Transaction)
// 		assert.NotEmpty(t, transaction.Timeout)
// 		assert.NotEmpty(t, transaction.Device)
// 		assert.Empty(t, transaction.Message)
// 	}
// }
