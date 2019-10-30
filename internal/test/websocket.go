package test

// websocket.go provides testing functionalities against a mock websocket server.

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httptest"

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

// NewWebSocketServerV3 returns an instance of a mock websocket server for v3 API.
func NewWebSocketServerV3() WebSocketServer {
	m := http.NewServeMux()
	s := httptest.NewServer(m)

	return WebSocketServer{
		URL:        s.URL[7:], // remove `http://` prefix
		server:     s,
		mux:        m,
		version:    "v3",
		entryRoute: "connect",
	}
}

// NewWebSocketTLSServerV3 returns an instance of a mock websocket tls server for v3 API.
func NewWebSocketTLSServerV3() WebSocketServer {
	m := http.NewServeMux()
	s := httptest.NewTLSServer(m)

	return WebSocketServer{
		URL:        s.URL[8:], // remove `https://` prefix
		server:     s,
		mux:        m,
		version:    "v3",
		entryRoute: "connect",
	}
}

// Serve reads a request event and writes back a given response.
func (s WebSocketServer) Serve(resp string) {
	s.mux.HandleFunc(
		fmt.Sprintf("/%s/%s", s.version, s.entryRoute),
		func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}

			defer func() {
				err := c.Close()
				if err != nil {
					return
				}
			}()

			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					return
				}

				err = c.WriteMessage(websocket.TextMessage, []byte(resp))
				if err != nil {
					return
				}
			}
		},
	)
}

// Stream issues a requests and writes back a stream of response data.
func (s WebSocketServer) Stream(responses []string) {
	s.mux.HandleFunc(
		fmt.Sprintf("/%s/%s", s.version, s.entryRoute),
		func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer func() {
				if err := c.Close(); err != nil {
					fmt.Println(err)
					return
				}
			}()

			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, resp := range responses {
					err = c.WriteMessage(websocket.TextMessage, []byte(resp))
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		},
	)
}

// SetTLS starts TLS using the configured options.
func (s WebSocketServer) SetTLS(cfg *tls.Config) {
	s.tls = cfg
}

// GetCertificates returns the certificate used by the server.
func (s WebSocketServer) GetCertificates() *x509.Certificate {
	return s.server.Certificate()
}

// Close closes the connection.
func (s WebSocketServer) Close() {
	s.server.Close()
}
