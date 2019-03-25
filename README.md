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
instances, we need to pass in the configuration
[options](https://godoc.org/github.com/vapor-ware/synse-client-go/synse#Options),
namely the Synse Server's address that we need to interface wish, other
associated HTTP or WebSocket configs, and TLS communication configs if enabled.
For example,
```go
import (
	"github.com/vapor-ware/synse-client-go/synse"
)

func main() {
	opts := &synse.Options{
		Address: "localhost",
	}

	client, err := NewHTTPClientV3(opts) // or NewWebSocketClientV3(opts)
}
```

### API

Below are the map of HTTP and WebSocket API with its corresponding methods.

| HTTP endpoint | WebSocket request | Client method |
| ------------- | ----------------- | ------ |
| `/test` | `request/status` | `Status()` |
| `/version` | `request/version` | `Version()` |
