package test

// server.go provides testing functionalities against a mock http server.

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// HTTPServer describes a mock http server.
type HTTPServer struct {
	// URL is the base URL form http://ipaddr:port with no trailing slash.
	URL string

	// tls holds the TLS configuration.
	tls *tls.Config

	// server is the mock http server.
	server *httptest.Server

	// mux is the http request multiplexer.
	mux *http.ServeMux

	// version is the current api version of Synse Server that we are
	// communicating with.
	version string
}

// NewHTTPServerV3 returns an instance of a mock http server for v3 API.
func NewHTTPServerV3() HTTPServer {
	m := http.NewServeMux()
	s := httptest.NewServer(m)

	return HTTPServer{
		URL:     s.URL[7:],
		server:  s,
		mux:     m,
		version: "v3",
	}
}

// ServeUnversioned serves an unversioned endpoint.
func (s HTTPServer) ServeUnversioned(t *testing.T, uri string, statusCode int, response interface{}) {
	serve(s.mux, t, uri, statusCode, response)
}

// ServeVersioned serves a versioned endpoint.
func (s HTTPServer) ServeVersioned(t *testing.T, uri string, statusCode int, response interface{}) {
	// FIXME - need a better way to handle this. This might relate to #6 with
	// the use of https://golang.org/pkg/net/url/.
	serve(s.mux, t, fmt.Sprintf("/%v%v", s.version, uri), statusCode, response)
}

// SetTLS starts TLS using the configured options.
func (s HTTPServer) SetTLS(cfg *tls.Config) {
	s.tls = cfg
}

// Close closes the unversioned server connection.
func (s HTTPServer) Close() {
	s.server.Close()
}

// serve registers a path handler and writes to its responses.
func serve(m *http.ServeMux, t *testing.T, uri string, statusCode int, response interface{}) {
	m.HandleFunc(
		uri,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			fprint(t, w, response)
		},
	)
}

// fprint calls fmt.Fprint and validates its returned error.
func fprint(t *testing.T, w io.Writer, a interface{}) {
	_, err := fmt.Fprint(w, a)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
}
