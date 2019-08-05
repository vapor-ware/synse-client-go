[![Build Status](https://build.vio.sh/buildStatus/icon?job=vapor-ware/synse-client-go/master)](https://build.vio.sh/blue/organizations/jenkins/vapor-ware%2Fsynse-client-go/activity)
[![Godoc](https://godoc.org/github.com/vapor-ware/synse-client-go/synse?status.svg)](https://godoc.org/github.com/vapor-ware/synse-client-go/synse)
[![Go Report Card](https://goreportcard.com/badge/github.com/vapor-ware/synse-client-go)](https://goreportcard.com/report/github.com/vapor-ware/synse-client-go)

# synse-client-go

This repo contains a HTTP and WebSocket client for interacting with Synse
Server API, written in Go. For more information about the API, please visit its
specification at [HTTP
API](https://github.com/vapor-ware/synse-server/blob/master/proposals/v3/api.md),
and [WebSocket
API](https://github.com/vapor-ware/synse-server/blob/master/proposals/v3/api-websocket.md).

## Installing

```
git clone https://github.com/vapor-ware/synse-client-go.git
```

## Using

### Initializing

For both HTTP client and WebSocket client, in order to initialize their
instances, we need to pass in the [configuration
options](https://godoc.org/github.com/vapor-ware/synse-client-go/synse#Options),
namely the Synse Server's address that we need to interface with, other
associated HTTP or WebSocket configs, and TLS communication. For example:
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

Below is the map of client methods with their corresponding HTTP and WebSocket API.

| Method | HTTP endpoint | WebSocket request |
| ------ | ------------- | ----------------- |
| `Status()` | `/test` | `request/status` |
| `Version()` | `/version` | `request/version` |
| `Config()` | `/v3/config` | `request/config` |
| `Plugins()` | `/v3/plugin` | `request/plugin` |
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
| `Transactions()` | `/v3/transaction` | `request/transaction` |
| `Transaction(string)` | `/v3/transaction/{transaction_id}` | `request/transaction` |

Other than these:

| Method | Description |
| ------ | ----------- |
| `GetOptions()` | Return the current config options of the client. |
| `Open()` | Open the WebSocket connection between the client and Synse Server. *Not applicable for a HTTP client.* |
| `Close()` | Close the WebSocket connection between the client and Synse Server. *Not applicable for a HTTP client.* |

For more information about the response scheme, please refer to the
[documentation](https://godoc.org/github.com/vapor-ware/synse-client-go/synse#Client).


## Developing

To lint and format the project:
```
make lint
make fmt
```

To run all unit tests:
```
make test
```
