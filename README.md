[![Build Status](https://github.com/vapor-ware/synse-client-go/workflows/build/badge.svg)](https://github.com/vapor-ware/synse-client-go/actions)
[![Godoc](https://godoc.org/github.com/vapor-ware/synse-client-go/synse?status.svg)](https://godoc.org/github.com/vapor-ware/synse-client-go/synse)
[![Go Report Card](https://goreportcard.com/badge/github.com/vapor-ware/synse-client-go)](https://goreportcard.com/report/github.com/vapor-ware/synse-client-go)

# Synse Client (Golang)

The official HTTP and WebSocket client for interacting with the Synse Server API,
written in Go.

For details on the Synse Server API, see the
[API Documentation](https://synse.readthedocs.io/en/latest/server/api.v3/).

## Installing

This package can be installed via `go get`

```
go get github.com/vapor-ware/synse-client-go
```

## Using

### Initializing

Both the HTTP and WebSocket client must be initialized with
[configuration options](https://godoc.org/github.com/vapor-ware/synse-client-go/synse#Options).
These options identify the Synse Server instance to communicate with as well as
set other connection-related capabilities (timeout, TLS, etc).

```go
import "github.com/vapor-ware/synse-client-go/synse"

func main() {
	opts := &synse.Options{
		Address: "localhost",
	}

	client, err := NewHTTPClientV3(opts)
	// or client, err := NewWebSocketClientV3(opts)
}
```

### API

The table below describes which API endpoint/event correspond with each client method.

| Method | HTTP endpoint | WebSocket request |
| ------ | ------------- | ----------------- |
| `Status()` | `/test` | `request/status` |
| `Version()` | `/version` | `request/version` |
| `Config()` | `/v3/config` | `request/config` |
| `Plugins()` | `/v3/plugin` | `request/plugins` |
| `Plugin(string)` | `/v3/plugin/{plugin_id}` | `request/plugin` |
| `PluginHealth()` | `/v3/plugin/health` | `request/plugin_health` |
| `Scan(scheme.ScanOptions)` | `/v3/scan` | `request/scan` |
| `Tags(scheme.TagsOptions)` | `/v3/tags` | `request/tags` |
| `Info(string)` | `/v3/info/{device_id}` | `request/info` |
| `Read(scheme.ReadOptions)` | `/v3/read` | `request/read` |
| `ReadDevice(string, scheme.ReadOptions)` | `/v3/read/{device_id}` | `request/read_device` |
| `ReadCache(scheme.ReadCacheOptions)` | `/v3/readcache` | `request/read_cache` |
| `WriteAsync(string, []scheme.WriteData)` | `/v3/write/{device_id}` | `request/write_async` |
| `WriteSync(string, []scheme.WriteData)` | `/v3/write/wait/{device_id}` | `request/write_sync` |
| `Transactions()` | `/v3/transaction` | `request/transactions` |
| `Transaction(string)` | `/v3/transaction/{transaction_id}` | `request/transaction` |

Additionally, there are a few client methods which do not correspond to an API endpoint:

| Method | Description |
| ------ | ----------- |
| `GetOptions()` | Return the current config options of the client. |
| `Open()` | Open the WebSocket connection between the client and Synse Server. *WebSocket client only.* |
| `Close()` | Close the WebSocket connection between the client and Synse Server. *WebSocket client only.* |

For more information about the response scheme, please refer to the
[documentation](https://godoc.org/github.com/vapor-ware/synse-client-go/synse#Client).

## Developing

To provide a simple and uniform development flow, Makefile targets should be used for
basic development actions. To lint and format the project source code:

```
make lint
make fmt
```

To run all unit tests:

```
make test
```

For a full accounting of available targets, see the Makefile or run `make help`.
