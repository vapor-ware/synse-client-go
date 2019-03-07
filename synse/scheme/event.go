package scheme

// EventMeta describes core elements in a event request/response scheme.
type EventMeta struct {
	ID    uint64 `json:"id"`
	Event string `json:"event"`
}

// RequestVersion describes a scheme for request/version event.
type RequestVersion struct {
	EventMeta
}

// RequestConfig describes a scheme for request/config event.
type RequestConfig struct {
	EventMeta
}

// RequestPlugin describes a scheme for request/plugin event.
type RequestPlugin struct {
	EventMeta
	Data WriteData `json:"data"`
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
	Data WriteData `json:"data"`
}

// RequestRead describes a scheme for request/read event.
type RequestRead struct {
	EventMeta
	Data ReadOptions `json:"data"`
}

// RequestReadCache describes a scheme for request/read_cache event.
type RequestReadCache struct {
	EventMeta
	Data ReadCacheOptions `json:"data"`
}

// RequestWrite describes a scheme for request/write event.
type RequestWrite struct {
	EventMeta
	Data WriteData `json:"data"`
}

// RequestTransaction describes a scheme for request/transaction event.
type RequestTransaction struct {
	EventMeta
	Data WriteData `json:"data"`
}

// ResponseVersion describes a scheme for response/version event.
type ResponseVersion struct {
	EventMeta
	Data Version
}

// ResponseConfig describes a scheme for response/config event.
type ResponseConfig struct {
	EventMeta
	Data Config
}

// ResponsePlugin describes a scheme for response/plugin event.
type ResponsePlugin struct {
	EventMeta
	Data Plugin
}

// ResponsePluginHealth describes a scheme for response/plugin_health event.
type ResponsePluginHealth struct {
	EventMeta
	Data PluginHealth
}

// ResponseTags describes a scheme for response/tags event.
type ResponseTags struct {
	EventMeta
	Data TagsData
}

// TagsData describes the data for response/data event.
type TagsData struct {
	Tags []string `json:"tags"`
}

// ResponseDevice describes a scheme for response/device event.
type ResponseDevice struct {
	EventMeta
	Data Info
}

// ResponseDeviceSummary describes a scheme for response/device_summary event.
type ResponseDeviceSummary struct {
	EventMeta
	Data []Scan
}

// ResponseReading describes a scheme for response/reading event.
type ResponseReading struct {
	EventMeta
	Data []Read
}

// ResponseWriteState describes a scheme for response/write_state event.
type ResponseWriteState struct {
	EventMeta
	Data Transaction
}

// ResponseError describes a scheme for response/error event.
type ResponseError struct {
	EventMeta
	Data Error
}
