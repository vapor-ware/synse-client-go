package synse

import (
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
	assert.Equal(t, "debug", config.Logging)
	assert.True(t, config.PrettyJSON)
	assert.Equal(t, "en_US", config.Locale)
	assert.Equal(t, 1, len(config.Plugin.TCP))
	assert.Equal(t, "emulator:5001", config.Plugin.TCP[0])
	assert.Equal(t, 0, len(config.Plugin.Unix))
	// assert.Equal(t, 180, config.Cache.Device.RebuildEvery) // TODO: update new schema changes from upstream
	assert.Equal(t, 300, config.Cache.Transaction.TTL)
	assert.Equal(t, 3, config.GRPC.Timeout)
	assert.False(t, config.Metrics.Enabled)
}

// func TestIntegration_Plugins(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	plugins, err := client.Plugins()
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(plugins))

// 	for _, plugin := range plugins {
// 		assert.NotEmpty(t, plugin.Name)
// 		assert.NotEmpty(t, plugin.Maintainer)
// 		assert.NotEmpty(t, plugin.Tag)
// 		assert.NotEmpty(t, plugin.Description)
// 		assert.NotEmpty(t, plugin.ID)
// 		assert.True(t, plugin.Active)
// 	}
// }

// func TestIntegration_PluginInfo(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	plugins, err := client.Plugins()
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(plugins))

// 	for _, p := range plugins {
// 		plugin, err := client.Plugin(p.ID)
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, plugin.Name)
// 		assert.NotEmpty(t, plugin.Maintainer)
// 		assert.NotEmpty(t, plugin.Tag)
// 		assert.NotEmpty(t, plugin.Description)
// 		assert.NotEmpty(t, plugin.VCS)
// 		assert.NotEmpty(t, plugin.ID)
// 		assert.True(t, plugin.Active)
// 		assert.NotEmpty(t, plugin.Network.Address)
// 		assert.NotEmpty(t, plugin.Network.Protocol)
// 		assert.NotEmpty(t, plugin.Version.PluginVersion)
// 		assert.NotEmpty(t, plugin.Version.SDKVersion)
// 		assert.NotEmpty(t, plugin.Version.BuildDate)
// 		assert.NotEmpty(t, plugin.Version.GitCommit)
// 		assert.NotEmpty(t, plugin.Version.GitTag)
// 		assert.NotEmpty(t, plugin.Version.Arch)
// 		assert.NotEmpty(t, plugin.Version.OS)
// 		assert.NotEmpty(t, plugin.Health.Timestamp)
// 		assert.Equal(t, "OK", plugin.Health.Status)
// 		assert.Equal(t, 2, len(plugin.Health.Checks))

// 		for _, check := range plugin.Health.Checks {
// 			assert.NotEmpty(t, check.Name)
// 			assert.Equal(t, "OK", check.Status)
// 			assert.NotEmpty(t, check.Type)
// 			assert.NotEmpty(t, check.Timestamp)
// 			assert.Empty(t, check.Message)
// 		}

// 	}
// }

// func TestIntegration_PluginHealth(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	health, err := client.PluginHealth()
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, health.Status)
// 	assert.NotEmpty(t, health.Updated)
// 	assert.Equal(t, 1, len(health.Healthy))
// 	assert.Equal(t, 0, len(health.Unhealthy))
// 	assert.Equal(t, 1, health.Active)
// 	assert.Equal(t, 0, health.Inactive)
// }

// func TestIntegration_Scan(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	opts := scheme.ScanOptions{}
// 	devices, err := client.Scan(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 22, len(devices))

// 	for _, device := range devices {
// 		assert.NotEmpty(t, device.ID)
// 		assert.NotEmpty(t, device.Info)
// 		assert.NotEmpty(t, device.Type)
// 		assert.NotEmpty(t, device.Plugin)
// 		assert.NotEmpty(t, device.Tags)

// 		// only led devices have alias
// 		if device.Type == "led" {
// 			assert.NotEmpty(t, device.Alias)
// 		} else {
// 			assert.Empty(t, device.Alias)
// 		}
// 	}
// }

// func TestIntegration_Tags(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	opts := scheme.TagsOptions{}
// 	tags, err := client.Tags(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 10, len(tags))

// 	for _, tag := range tags {
// 		assert.NotEmpty(t, tag)
// 	}
// }

// func TestIntegration_Info(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	opts := scheme.ScanOptions{}
// 	devices, err := client.Scan(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 22, len(devices))

// 	for _, device := range devices {
// 		assert.NotEmpty(t, device)
// 	}

// 	for _, device := range devices {
// 		assert.NotEmpty(t, device)

// 		info, err := client.Info(device.ID)
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, info.Timestamp)
// 		assert.NotEmpty(t, info.ID)
// 		assert.NotEmpty(t, info.Type)
// 		assert.NotEmpty(t, info.Plugin)
// 		assert.NotEmpty(t, info.Capabilities.Mode)
// 		assert.NotEmpty(t, info.Tags)
// 		assert.NotEmpty(t, info.Metadata)
// 		assert.Empty(t, info.Capabilities.Write.Actions)

// 		// TODO - add sort_index to the scheme

// 		// only led devices have alias
// 		if info.Type == "led" {
// 			assert.NotEmpty(t, info.Alias)
// 		} else {
// 			assert.Empty(t, info.Alias)
// 		}

// 		for _, output := range info.Outputs {
// 			assert.NotEmpty(t, output.Name)
// 			// assert.NotEmpty(t, output.Type) // FIXME - only airflow types are empty?

// 			// led and lock devices don't produce unit and precision output
// 			if output.Name == "state" || output.Name == "color" || output.Name == "status" {
// 				assert.Empty(t, output.Unit.Name)
// 				assert.Empty(t, output.Unit.Symbol)
// 				assert.Empty(t, output.Precision)
// 			} else {
// 				assert.NotEmpty(t, output.Unit.Name)
// 				assert.NotEmpty(t, output.Unit.Symbol)
// 				assert.NotEmpty(t, output.Precision)
// 			}

// 			assert.Empty(t, output.ScalingFactor)
// 		}
// 	}
// }

// func TestIntegration_Read(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	opts := scheme.ReadOptions{}
// 	devices, err := client.Read(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 25, len(devices))

// 	for _, device := range devices {
// 		assert.NotEmpty(t, device.Device)
// 		assert.NotEmpty(t, device.Timestamp)
// 		// assert.NotEmpty(t, device.Type) // FIXME - fan and airflow types are empty?
// 		assert.NotEmpty(t, device.DeviceType)

// 		// led and lock devices don't produce unit output
// 		if device.DeviceType == "led" || device.DeviceType == "lock" {
// 			assert.Empty(t, device.Unit.Name)
// 			assert.Empty(t, device.Unit.Symbol)
// 		} else {
// 			assert.NotEmpty(t, device.Unit.Name)
// 			assert.NotEmpty(t, device.Unit.Symbol)
// 		}

// 		// NOTE - device.Value could be 0 so no need to check that
// 	}
// }

// func TestIntegration_ReadDevice(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewHTTPClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	opts := scheme.ScanOptions{}
// 	devices, err := client.Scan(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 22, len(devices))

// 	for _, device := range devices {
// 		opts := scheme.ReadOptions{}
// 		reads, err := client.ReadDevice(device.ID, opts)
// 		assert.NoError(t, err)

// 		for _, read := range reads {
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
// 		}
// 	}
// }

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
