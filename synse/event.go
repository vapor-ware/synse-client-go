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

	// requestPlugins describes the request/plugins event.
	requestPlugins = "request/plugins"

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

	// requestReadDevice describes the request/read_device event.
	requestReadDevice = "request/read_device"

	// requestReadCache describes the request/read_cache event.
	requestReadCache = "request/read_cache"

	// requestWriteAsync describes the request/write_async event.
	requestWriteAsync = "request/write_async"

	// requestWriteSync describes the request/write_sync event.
	requestWriteSync = "request/write_sync"

	// requestTransactions describes the request/transactions event.
	requestTransactions = "request/transactions"

	// requestTransaction describes the request/transaction event.
	requestTransaction = "request/transaction"

	// responseStatus describes the response/status event.
	responseStatus = "response/status"

	// responseVersion describes the response/version event.
	responseVersion = "response/version"

	// responseConfig describes the response/config event.
	responseConfig = "response/config"

	// responsePlugin describes the response/plugin_info event.
	responsePlugin = "response/plugin_info"

	// responsePluginSummary describes the response/plugin_summary event.
	responsePluginSummary = "response/plugin_summary"

	// responsePluginHealth describes the response/plugin_health event.
	responsePluginHealth = "response/plugin_health"

	// responseTags describes the response/tags event.
	responseTags = "response/tags"

	// responseDevice describes the response/device event.
	responseDevice = "response/device_info"

	// responseDeviceSummary describes the response/device_summary event.
	responseDeviceSummary = "response/device_summary"

	// responseReading describes the response/reading event.
	responseReading = "response/reading"

	// responseWriteAsync describes the response/write_async event.
	responseWriteAsync = "response/write_async"

	// responseWriteSync describes the response/write_sync event.
	responseWriteSync = "response/write_sync"

	// responseTransactionList describes the response/transaction_list event.
	responseTransactionList = "response/transaction_list"

	// responseTransactionInfo describes the response/transaction_info event.
	responseTransactionInfo = "response/transaction_info"

	// responseError describes the response/error event.
	responseError = "response/error"
)
