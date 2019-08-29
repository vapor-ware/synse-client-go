package synse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
)

func TestIntegrationWebSocket_Status(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

	status, err := client.Status()
	assert.NoError(t, err)
	assert.Equal(t, "ok", status.Status)
	assert.NotEmpty(t, status.Timestamp)
}

func TestIntegrationWebSocket_Version(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

	version, err := client.Version()
	assert.NoError(t, err)
	assert.NotEmpty(t, version.Version)
	assert.Equal(t, "v3", version.APIVersion)
}

func TestIntegrationWebSocket_Config(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

	config, err := client.Config()
	assert.NoError(t, err)
	assert.Equal(t, "debug", config.Logging)
	assert.True(t, config.PrettyJSON)
	assert.Equal(t, "en_US", config.Locale)
	assert.Equal(t, 1, len(config.Plugin.TCP))
	assert.Equal(t, "emulator:5001", config.Plugin.TCP[0])
	assert.Equal(t, 0, len(config.Plugin.Unix))
	assert.Equal(t, 180, config.Cache.Device.RebuildEvery)
	assert.Equal(t, 300, config.Cache.Transaction.TTL)
	assert.Equal(t, 3, config.GRPC.Timeout)
	assert.False(t, config.Metrics.Enabled)
}

func TestIntegrationWebSocket_Plugins(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

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

func TestIntegrationWebSocket_PluginInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

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
	assert.Empty(t, readCheck.Message)

	writeCheck := plugin.Health.Checks[1]
	assert.Equal(t, "write queue health", writeCheck.Name)
	assert.Equal(t, "OK", writeCheck.Status)
	assert.Equal(t, "periodic", writeCheck.Type)
	assert.Empty(t, writeCheck.Message)

	// NOTE - health check timestamp is not populated after at least 30s of
	// deployment. that's a pretty long time so we won't check that for now.
	// assert.NotEmpty(t, readCheck.Timestamp)
	// assert.NotEmpty(t, writeCheck.Timestamp)
}

func TestIntegrationWebSocket_PluginHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

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

// FIXME - got a 500 error response from synse server at 2019-08-29T15:55:52Z,
// saying An unexpected error occurred., with context: object of type
// 'NoneType' has no len()
// func TestIntegrationWebSocket_Scan(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewWebSocketClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	err = client.Open()
// 	assert.NoError(t, err)

// 	defer func() {
// 		err = client.Close()
// 		assert.NoError(t, err)
// 	}()

// 	opts := scheme.ScanOptions{}
// 	devices, err := client.Scan(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 4, len(devices))

// 	// since scan responses are sorted by default, the devices order should be
// 	// consistent.
// 	tempDevice1 := devices[0]
// 	assert.Equal(t, "89fd576d-462c-53be-bcb6-7870e70c304a", tempDevice1.ID)
// 	assert.Empty(t, tempDevice1.Alias)
// 	assert.Equal(t, "Synse Temperature Sensor 2", tempDevice1.Info)
// 	assert.Equal(t, "temperature", tempDevice1.Type)
// 	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", tempDevice1.Plugin)
// 	assert.Equal(t, 3, len(tempDevice1.Tags))
// 	assert.Equal(t, "foo/bar", tempDevice1.Tags[0])
// 	assert.Equal(t, "system/id:89fd576d-462c-53be-bcb6-7870e70c304a", tempDevice1.Tags[1])
// 	assert.Equal(t, "system/type:temperature", tempDevice1.Tags[2])

// 	tempDevice2 := devices[1]
// 	assert.Equal(t, "9907bdfa-75e1-5af5-8385-87184f356b22", tempDevice2.ID)
// 	assert.Empty(t, tempDevice2.Alias)
// 	assert.Equal(t, "Synse Temperature Sensor 1", tempDevice2.Info)
// 	assert.Equal(t, "temperature", tempDevice2.Type)
// 	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", tempDevice2.Plugin)
// 	assert.Equal(t, 3, len(tempDevice2.Tags))
// 	assert.Equal(t, "foo/bar", tempDevice2.Tags[0])
// 	assert.Equal(t, "system/id:9907bdfa-75e1-5af5-8385-87184f356b22", tempDevice2.Tags[1])
// 	assert.Equal(t, "system/type:temperature", tempDevice2.Tags[2])

// 	tempDevice3 := devices[2]
// 	assert.Equal(t, "b9324904-385b-581d-b790-5e53eaabfd20", tempDevice3.ID)
// 	assert.Equal(t, "emulator-temp", tempDevice3.Alias)
// 	assert.Equal(t, "Synse Temperature Sensor 3", tempDevice3.Info)
// 	assert.Equal(t, "temperature", tempDevice3.Type)
// 	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", tempDevice3.Plugin)
// 	assert.Equal(t, 2, len(tempDevice3.Tags))
// 	assert.Equal(t, "system/id:b9324904-385b-581d-b790-5e53eaabfd20", tempDevice3.Tags[0])
// 	assert.Equal(t, "system/type:temperature", tempDevice3.Tags[1])

// 	ledDevice := devices[3]
// 	assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", ledDevice.ID)
// 	assert.Equal(t, "emulator-led", ledDevice.Alias)
// 	assert.Equal(t, "Synse LED", ledDevice.Info)
// 	assert.Equal(t, "led", ledDevice.Type)
// 	assert.Equal(t, "4032ffbe-80db-5aa5-b794-f35c88dff85c", ledDevice.Plugin)
// 	assert.Equal(t, 3, len(ledDevice.Tags))
// 	assert.Equal(t, "foo/bar", ledDevice.Tags[0])
// 	assert.Equal(t, "system/id:f041883c-cf87-55d7-a978-3d3103836412", ledDevice.Tags[1])
// 	assert.Equal(t, "system/type:led", ledDevice.Tags[2])
// }

// FIXME - got a 500 error response from synse server at 2019-08-29T15:59:15Z,
// saying An unexpected error occurred., with context: tags() argument after *
// must be an iterable, not NoneType
// func TestIntegrationWebSocket_Tags(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewWebSocketClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	err = client.Open()
// 	assert.NoError(t, err)

// 	defer func() {
// 		err = client.Close()
// 		assert.NoError(t, err)
// 	}()

// 	opts := scheme.TagsOptions{}
// 	tags, err := client.Tags(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 3, len(tags))
// 	assert.Equal(t, "foo/bar", tags[0])
// 	assert.Equal(t, "system/type:led", tags[1])
// 	assert.Equal(t, "system/type:temperature", tags[2])
// }

func TestIntegrationWebSocket_Info(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

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
	assert.Equal(t, 3, len(device.Tags))
	assert.Equal(t, "foo/bar", device.Tags[0])
	assert.Equal(t, "system/id:f041883c-cf87-55d7-a978-3d3103836412", device.Tags[1])
	assert.Equal(t, "system/type:led", device.Tags[2])
	assert.Equal(t, 2, len(device.Outputs))
	assert.Equal(t, 0, device.SortIndex)

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
}

// FIXME - got a 500 error response from synse server at 2019-08-29T16:05:38Z,
// saying An unexpected error occurred., with context: object of type
// 'NoneType' has no len()
// func TestIntegrationWebSocket_Read(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	client, err := NewWebSocketClientV3(&Options{
// 		Address: "localhost:5000",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	err = client.Open()
// 	assert.NoError(t, err)

// 	defer func() {
// 		err = client.Close()
// 		assert.NoError(t, err)
// 	}()

// 	opts := scheme.ReadOptions{}
// 	readings, err := client.Read(opts)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 5, len(readings))

// 	counts := countDeviceType(readings)
// 	assert.Equal(t, 2, counts["led"])
// 	assert.Equal(t, 3, counts["temperature"])

// 	for _, read := range readings {
// 		if read.DeviceType == "led" {
// 			assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", read.Device)
// 			assert.NotEmpty(t, read.Timestamp)
// 			assert.Contains(t, []string{"state", "color"}, read.Type)
// 			assert.Empty(t, read.Unit)
// 			assert.Contains(t, []string{"off", "000000"}, read.Value)
// 			assert.Empty(t, read.Context)
// 		} else if read.DeviceType == "temperature" {
// 			assert.Contains(t, []string{"89fd576d-462c-53be-bcb6-7870e70c304a", "9907bdfa-75e1-5af5-8385-87184f356b22", "b9324904-385b-581d-b790-5e53eaabfd20"}, read.Device)
// 			assert.NotEmpty(t, read.Timestamp)
// 			assert.Equal(t, "temperature", read.Type)
// 			assert.Equal(t, "celsius", read.Unit.Name)
// 			assert.Equal(t, "C", read.Unit.Symbol)
// 			assert.NotEmpty(t, read.Value)
// 			assert.Empty(t, read.Context)
// 		} else {
// 			t.Error("unexpected reading device type in response")
// 		}
// 	}
// }

func TestIntegrationWebSocket_ReadDevice(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

	readings, err := client.ReadDevice("f041883c-cf87-55d7-a978-3d3103836412")
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

func TestIntegrationWebSocket_ReadCache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	defer func() {
		err = client.Close()
		assert.NoError(t, err)
	}()

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

			if read.DeviceType == "led" {
				assert.Equal(t, "f041883c-cf87-55d7-a978-3d3103836412", read.Device)
				assert.NotEmpty(t, read.Timestamp)
				assert.Contains(t, []string{"state", "color"}, read.Type)
				assert.Empty(t, read.Unit)
				assert.Contains(t, []string{"off", "000000"}, read.Value)
				assert.Empty(t, read.Context)
			} else if read.DeviceType == "temperature" {
				assert.Contains(t, []string{"89fd576d-462c-53be-bcb6-7870e70c304a", "9907bdfa-75e1-5af5-8385-87184f356b22", "b9324904-385b-581d-b790-5e53eaabfd20"}, read.Device)
				assert.NotEmpty(t, read.Timestamp)
				assert.Equal(t, "temperature", read.Type)
				assert.Equal(t, "celsius", read.Unit.Name)
				assert.Equal(t, "C", read.Unit.Symbol)
				assert.NotEmpty(t, read.Value)
				assert.Empty(t, read.Context)
			} else {
				t.Error("unexpected reading device type in response")
			}

		case <-time.After(2 * time.Second):
			// if the test does not complete after 2s, error.
			t.Fatal("timeout: failed getting readcache data from channel")
		}

		if done {
			break
		}
	}
}
