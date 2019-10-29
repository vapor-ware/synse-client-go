package synse

// event.go stores request/response events supported by the Synse WebSocket API.

const (

	// Synse Server WebSocket API request events.
	requestStatus       = "request/status"
	requestVersion      = "request/version"
	requestConfig       = "request/config"
	requestPlugin       = "request/plugin"
	requestPlugins      = "request/plugins"
	requestPluginHealth = "request/plugin_health"
	requestScan         = "request/scan"
	requestTags         = "request/tags"
	requestInfo         = "request/info"
	requestRead         = "request/read"
	requestReadDevice   = "request/read_device"
	requestReadCache    = "request/read_cache"
	requestReadStream   = "request/read_stream"
	requestWriteAsync   = "request/write_async"
	requestWriteSync    = "request/write_sync"
	requestTransaction  = "request/transaction"
	requestTransactions = "request/transactions"

	// Synse Server WebSocket API response events.
	responseStatus            = "response/status"
	responseVersion           = "response/version"
	responseConfig            = "response/config"
	responsePluginInfo        = "response/plugin_info"
	responsePluginSummary     = "response/plugin_summary"
	responsePluginHealth      = "response/plugin_health"
	responseDeviceSummary     = "response/device_summary"
	responseTags              = "response/tags"
	responseDeviceInfo        = "response/device_info"
	responseReading           = "response/reading"
	responseTransactionInfo   = "response/transaction_info"
	responseTransactionStatus = "response/transaction_status"
	responseTransactionList   = "response/transaction_list"
	responseError             = "response/error"
)
