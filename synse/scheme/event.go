package scheme

// EventMeta describes core elements in a event request/response scheme.
type EventMeta struct {
	ID    uint64 `json:"id" yaml:"id"`
	Event string `json:"event" yaml:"event"`
}

// RequestStatus describes a scheme for request/status event.
type RequestStatus struct {
	EventMeta
}

// RequestVersion describes a scheme for request/version event.
type RequestVersion struct {
	EventMeta
}

// RequestConfig describes a scheme for request/config event.
type RequestConfig struct {
	EventMeta
}

// RequestPlugins describes a scheme for request/plugin event, with no
// plugin id being provided.
type RequestPlugins struct {
	EventMeta
}

// RequestPlugin describes a scheme for request/plugin event.
type RequestPlugin struct {
	EventMeta
	Data PluginData `json:"data"`
}

// PluginData describes the data for request/plugin event.
type PluginData struct {
	Plugin string `json:"plugin"`
}

// RequestPluginHealth describes a scheme for request/plugin_health event.
type RequestPluginHealth struct {
	EventMeta
}

// RequestScan describes a scheme for request/scan event.
type RequestScan struct {
	EventMeta
	Data ScanOptions `json:"data"`
}

// RequestTags describes a scheme for request/tags event.
type RequestTags struct {
	EventMeta
	Data TagsOptions `json:"data"`
}

// RequestInfo describes a scheme for request/info event.
type RequestInfo struct {
	EventMeta
	Data DeviceData `json:"data"`
}

// DeviceData describes the data for response/info event.
type DeviceData struct {
	Device string `json:"device"`
}

// RequestRead describes a scheme for request/read event.
type RequestRead struct {
	EventMeta
	Data ReadOptions `json:"data"`
}

// RequestReadDevice describes a scheme for request/read_device event.
type RequestReadDevice struct {
	EventMeta
	Data ReadDeviceData `json:"data"`
}

// ReadDeviceData describes the data for request/read_device event.
type ReadDeviceData struct {
	ID string `json:"id"`
	ReadOptions
}

// RequestReadCache describes a scheme for request/read_cache event.
type RequestReadCache struct {
	EventMeta
	Data ReadCacheOptions `json:"data"`
}

// RequestWrite describes a scheme for request/write_async and
// request/write_sync event.
type RequestWrite struct {
	EventMeta
	Data RequestWriteData `json:"data"`
}

// RequestWriteData describes the data for request/write_async and
// request/write_sync event.
type RequestWriteData struct {
	ID      string      `json:"id"`
	Payload []WriteData `json:"payload"`
}

// RequestTransactions describes a scheme for request/transaction event with no
// transaction id being provided.
type RequestTransactions struct {
	EventMeta
}

// RequestTransaction describes a scheme for request/transaction event.
type RequestTransaction struct {
	EventMeta
	Data WriteData `json:"data"`
}

// Response describes a generic response scheme.
type Response struct {
	EventMeta
	Data interface{} `json:"data"`
}
