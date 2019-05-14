package scheme

// EventMeta describes core elements in a event request/response scheme.
type EventMeta struct {
	ID    uint64 `json:"id" yaml:"id" mapstructure:"id"`
	Event string `json:"event" yaml:"event" mapstructure:"event"`
}

// RequestStatus describes a scheme for request/status event.
type RequestStatus struct {
	EventMeta `mapstructure:",squash"`
}

// RequestVersion describes a scheme for request/version event.
type RequestVersion struct {
	EventMeta `mapstructure:",squash"`
}

// RequestConfig describes a scheme for request/config event.
type RequestConfig struct {
	EventMeta `mapstructure:",squash"`
}

// RequestPlugins describes a scheme for request/plugin event, with no
// plugin id being provided.
type RequestPlugins struct {
	EventMeta `mapstructure:",squash"`
}

// RequestPlugin describes a scheme for request/plugin event.
type RequestPlugin struct {
	EventMeta `mapstructure:",squash"`
	Data      PluginData `json:"data" yaml:"data" mapstructure:"data"`
}

// PluginData describes the data for request/plugin event.
type PluginData struct {
	Plugin string `json:"plugin" yaml:"plugin" mapstructure:"plugin"`
}

// RequestPluginHealth describes a scheme for request/plugin_health event.
type RequestPluginHealth struct {
	EventMeta `mapstructure:",squash"`
}

// RequestScan describes a scheme for request/scan event.
type RequestScan struct {
	EventMeta `mapstructure:",squash"`
	Data      ScanOptions `json:"data" yaml:"data" mapstructure:"data"`
}

// RequestTags describes a scheme for request/tags event.
type RequestTags struct {
	EventMeta `mapstructure:",squash"`
	Data      TagsOptions `json:"data" yaml:"data" mapstructure:"data"`
}

// RequestInfo describes a scheme for request/info event.
type RequestInfo struct {
	EventMeta `mapstructure:",squash"`
	Data      DeviceData `json:"data" yaml:"data" mapstructure:"data"`
}

// DeviceData describes the data for response/info event.
type DeviceData struct {
	Device string `json:"device" yaml:"device" mapstructure:"device"`
}

// RequestRead describes a scheme for request/read event.
type RequestRead struct {
	EventMeta `mapstructure:",squash"`
	Data      ReadOptions `json:"data" yaml:"data" mapstructure:"data"`
}

// RequestReadDevice describes a scheme for request/read_device event.
type RequestReadDevice struct {
	EventMeta `mapstructure:",squash"`
	Data      ReadDeviceData `json:"data" yaml:"data" mapstructure:"data"`
}

// ReadDeviceData describes the data for request/read_device event.
type ReadDeviceData struct {
	ID          string `json:"id" yaml:"id" mapstructure:"id"`
	ReadOptions `mapstructure:",squash"`
}

// RequestReadCache describes a scheme for request/read_cache event.
type RequestReadCache struct {
	EventMeta `mapstructure:",squash"`
	Data      ReadCacheOptions `json:"data" yaml:"data" mapstructure:"data"`
}

// RequestWrite describes a scheme for request/write_async and
// request/write_sync event.
type RequestWrite struct {
	EventMeta `mapstructure:",squash"`
	Data      RequestWriteData `json:"data" yaml:"data" mapstructure:"data"`
}

// RequestWriteData describes the data for request/write_async and
// request/write_sync event.
type RequestWriteData struct {
	ID      string      `json:"id" yaml:"id" mapstructure:"id"`
	Payload []WriteData `json:"payload" yaml:"payload" mapstructure:"payload"`
}

// RequestTransactions describes a scheme for request/transaction event with no
// transaction id being provided.
type RequestTransactions struct {
	EventMeta `mapstructure:",squash"`
}

// RequestTransaction describes a scheme for request/transaction event.
type RequestTransaction struct {
	EventMeta `mapstructure:",squash"`
	Data      WriteData `json:"data" yaml:"data" mapstructure:"data"`
}

// Response describes a generic response scheme.
type Response struct {
	EventMeta `mapstructure:",squash"`
	Data      interface{} `json:"data" yaml:"data" mapstructure:"data"`
}
