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

In order to install this package, you can clone the repo, `cd` into the repo
root, and install via `make`:
```
$ git clone https://github.com/vapor-ware/synse-client-go.git
$ cd synse-client-go
$ make setup
```

## Using

### Initializing

For both HTTP client and WebSocket client, in order to initialize their
instances, we need to pass in the [configuration
options](https://godoc.org/github.com/vapor-ware/synse-client-go/synse#Options),
namely the Synse Server's address that we need to interface wish, other
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

Below are the map of HTTP and WebSocket API with their corresponding methods.

| Name | HTTP endpoint | WebSocket request | Method |
| ---- | ------------- | ----------------- | ------ |
| Status | `/test` | `request/status` | `Status()` |
| Version | `/version` | `request/version` | `Version()` |
| Config | `/v3/config` | `request/config` | `Config()` |
| Plugins | `/v3/plugin[/{plugin_id}]` | `request/plugin` | `Plugins()` or `Plugin(string)` |
| Plugin Health | `/v3/plugin/health` | `request/plugin_health` | `PluginHealth()` |
| Scan | `/v3/scan` | `request/scan` | `Scan(scheme.ScanOptions)` |
| Tags | `/v3/tags` | `request/tags` | `Tags(scheme.TagsOptions)` |
| Info | `/v3/info/{device_id}` | `request/info` | `Info(string)` |
| Read | `/v3/read` | `request/read` | `Read(scheme.ReadOptions)` |
| Read Device | `/v3/read/{device_id}` | `request/read_device` | `ReadDevice(string, scheme.ReadOptions)` |
| Read Cache | `/v3/readcache` | `request/read_cache` | `ReadCache(scheme.ReadCacheOptions)` |
| Write (Asynchronous) | `/v3/write/{device_id}` | `request/write_async` | `WriteAsync(string, []scheme.WriteData)` |
| Write (Synchronous) | `/v3/write/wait/{device_id}` | `request/write_sync` | `WriteSync(string, []scheme.WriteData)` |
| Transaction | `/v3/transaction[/{transaction_id}]` | `request/transaction` | `Transactions()` or `Transaction(string)` |

Other than these,

| Method | Description |
| ------ | ----------- |
| `GetOptions()` | Return the current config options of the client |
| `Open()` | Open the WebSocket connection between the client and Synse Server. This is not applicable for a HTTP client |
| `Close()` | Close the WebSocket connection between the client and Synse Server. This is not applicable for a HTTP client |
