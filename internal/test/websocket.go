package test

// websocket.go provides testing functionalities against a mock websocket server.

import (
	"crypto/tls"
	// "crypto/x509"
	"fmt"
	// "io"
	"net/http"
	"net/http/httptest"
	// "testing"

	// "github.com/vapor-ware/synse-client-go/synse/scheme"

	"github.com/gorilla/websocket"
)

// upgrader holds the gorilla/websocket Upgrader for upgrading an HTTP
// connection to a WebSocket one.
var upgrader = websocket.Upgrader{}

// WebSocketServer describes a mock websocket server.
type WebSocketServer struct {
	// URL has the `host:port` format.
	URL string

	// tls holds the TLS configuration.
	tls *tls.Config

	// server is the mock websocket server.
	server *httptest.Server

	// mux is the http request multiplexer.
	mux *http.ServeMux

	// version is the current api version of Synse Server that we are
	// communicating with.
	version string

	// entryRoute is the entry route to start the websocket connection.
	entryRoute string
}

// request describes a request event scheme.
type request struct {
	ID    uint64                 `json:"id"`
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// NewWebSocketServerV3 returns an instance of a mock http server for v3 API.
func NewWebSocketServerV3() WebSocketServer {
	m := http.NewServeMux()
	s := httptest.NewServer(m)

	return WebSocketServer{
		URL:        s.URL[7:],
		server:     s,
		mux:        m,
		version:    "v3",
		entryRoute: "connect",
	}
}

func (s WebSocketServer) Serve(resp interface{}) {
	s.mux.HandleFunc(
		fmt.Sprintf("/%s/%s", s.version, s.entryRoute),
		func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			defer c.Close()

			for {
				in := new(request)
				err := c.ReadJSON(in)
				if err != nil {
					return
				}

				out := map[string]interface{}{
					"id":    in.ID,
					"event": matchEvent(in.Event),
					"data":  resp,
				}
				err = c.WriteJSON(out)
				if err != nil {
					return
				}
			}
		},
	)
}

func (s WebSocketServer) Close() {
	s.server.Close()
}

// TODO - clean this up
const (
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
	responseError = "response/error"
)

// matchEvent returns a corresponding response event for a given request event.
func matchEvent(reqEvent string) string {
	var respEvent string

	switch reqEvent {
	case requestVersion:
		respEvent = responseVersion
	case requestPlugin:
		respEvent = responsePlugin
	case requestPluginHealth:
		respEvent = responsePluginHealth
	case requestScan:
		respEvent = responseDeviceSummary
	default:
		respEvent = ""
	}

	return respEvent
}
