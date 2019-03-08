package synse

// event.go stores request/response events supported by the Synse WebSocket API.

const (
	// requestStatus describes the request/status event.
	requestStatus = "request/status"

	// requestVersion describes the request/version event.
	requestVersion = "request/version"

	// requestConfig describes the request/config event.
	requestConfig = "request/config"

	// requestPlugin describes the request/plugin event.
	requestPlugin = "request/plugin"

	// requestPluginHealth describes the request/plugin_health event.
	requestPluginHealth = "request/plugin_health"

	// requestScan describes the request/scan event.
	requestScan = "request/scan"

	// requestTags describes the request/tags event.
	requestTags = "request/tags"

	// requestInfo describes the request/info event.
	requestInfo = "request/info"

	// requestRead describes the request/read event.
	requestRead = "request/read"

	// requestReadCache describes the request/read_cache event.
	requestReadCache = "request/read_cache"

	// requestWrite describes the request/write event.
	requestWrite = "request/write"

	// requestTransaction describes the request/transaction event.
	requestTransaction = "request/transaction"

	// responseStatus describes the response/status event.
	responseStatus = "response/status"

	// responseVersion describes the response/version event.
	responseVersion = "response/version"

	// responseConfig describes the response/config event.
	responseConfig = "response/config"

	// responsePlugin describes the response/plugin event.
	responsePlugin = "response/plugin"

	// responsePluginHealth describes the response/plugin_health event.
	responsePluginHealth = "response/plugin_health"

	// responseTags describes the response/tags event.
	responseTags = "response/tags"

	// responseDevice describes the response/device event.
	responseDevice = "response/device"

	// responseDeviceSummary describes the response/device_summary event.
	responseDeviceSummary = "response/device_summary"

	// responseReading describes the response/reading event.
	responseReading = "response/reading"

	// responseWriteState describes the response/write_state event.
	responseWriteState = "response/write_state"

	// responseError describes the response/error event.
	// responseError = "response/error"
)
