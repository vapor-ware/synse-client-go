package scheme

// EventMeta describes core elements in a event request/response scheme.
type EventMeta struct {
	ID    string `json:"id"`
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
	Data ScanData `json:"data"`
}

// ScanData describes the data for request/scan event.
type ScanData struct {
	NS    string   `json:"ns"`
	Tags  []string `json:"tags"`
	Force bool     `json:"force"`
}

// RequestTags describes a scheme for request/tags event.
type RequestTags struct {
	EventMeta
	Data TagsData `json:"data"`
}

// TagsData describes the data for request/tags event.
type TagsData struct {
	NS  []string `json:"ns"`
	IDs bool     `json:"ids"`
}

// RequestInfo describes a scheme for request/info event.
type RequestInfo struct {
	EventMeta
	Data InfoData `json:"data"`
}

// InfoData describes the data for request/info event.
type InfoData struct {
	Device string `json:"device"`
}

// RequestRead describes a scheme for request/read event.
type RequestRead struct {
	EventMeta
	Data ReadData `json:"data"`
}

// ReadData describes the data for request/read event.
type ReadData struct {
	NS   string   `json:"ns"`
	Tags []string `json:"tags"`
}

// RequestReadCache describes a scheme for request/read_cache event.
type RequestReadCache struct {
	EventMeta
}

// RequestWrite describes a scheme for request/write event.
type RequestWrite struct {
	EventMeta
	Data WriteData `json:"data"`
}

// RequestTransaction describes a scheme for request/transaction event.
type RequestTransaction struct {
	EventMeta
	Data ReadData `json:"data"`
}

// TransactionData describes the data for request/transaction event.
type TransactionData struct {
	Transaction string `json:"transaction"`
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
