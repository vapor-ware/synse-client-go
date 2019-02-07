package test

// server.go provides testing functionalities against a mock http server.

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// UnversionedHTTPServer describes an unversioned mock http server.
type UnversionedHTTPServer struct {
	URL    string
	server *httptest.Server
	mux    *http.ServeMux
}

// VersionedHTTPServer describes a versioned mock http server.
type VersionedHTTPServer struct {
	UnversionedHTTPServer
	versionURI string
}

// NewUnversionedHTTPServer returns an instance of an unversioned mock http server.
func NewUnversionedHTTPServer() UnversionedHTTPServer {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	return UnversionedHTTPServer{
		URL:    s.URL[7:],
		server: s,
		mux:    m,
	}
}

// NewVersionedHTTPServer returns an instance of a versioned mock http server.
func NewVersionedHTTPServer() VersionedHTTPServer {
	s := NewUnversionedHTTPServer()
	s.mux.HandleFunc(
		"/version",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"version": "3.x.x", "api_version": "v3"}`) // nolint
		})

	return VersionedHTTPServer{
		UnversionedHTTPServer: s,
		versionURI:            "/v3",
	}
}

// Serve serves an unversioned endpoint.
func (s UnversionedHTTPServer) Serve(t *testing.T, uri string, statusCode int, response interface{}) {
	serve(s.mux, t, uri, statusCode, response)
}

// Serve serves a versioned endpoint.
func (s VersionedHTTPServer) Serve(t *testing.T, uri string, statusCode int, response interface{}) {
	serve(s.mux, t, fmt.Sprintf("%v%v", s.versionURI, uri), statusCode, response)
}

// Close closes the unversioned server connection.
func (s UnversionedHTTPServer) Close() {
	s.server.Close()
}

// Close closes the versioned server connection.
func (s VersionedHTTPServer) Close() {
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
