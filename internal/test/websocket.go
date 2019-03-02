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

type mock struct {
}

// NewServerV3 returns an instance of a mock http server for v3 API.
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

func (s WebSocketServer) Serve(event string, resp interface{}) {
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

				out := map[string]interface{}{}
				switch event {
				case "request/version":
					out = map[string]interface{}{
						"id":    in.ID,
						"event": in.Event,
						"data":  resp,
					}
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
