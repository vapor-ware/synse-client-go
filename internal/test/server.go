package test

// server.go provides testing functionalities against a test http server.

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockHTTPServer provides a test http server.
type MockHTTPServer struct {
	URL    string
	Server *httptest.Server
	mux    *http.ServeMux
}

// NewMockHTTPServer returns a new instance of a test http server.
func NewMockHTTPServer() MockHTTPServer {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	return MockHTTPServer{
		mux:    m,
		Server: s,
		URL:    s.URL[7:],
	}
}

// ServeUnversionedSuccess registers an unversioned success route.
func (s MockHTTPServer) ServeUnversionedSuccess(t *testing.T, path string, response interface{}) {
	s.serveUnversion(t, path, 200, response)
}

// ServeUnversionedFailure registers an unversioned failure route.
func (s MockHTTPServer) ServeUnversionedFailure(t *testing.T, path string, response interface{}) {
	s.serveUnversion(t, path, 500, response)
}

// ServeVersionedSuccess registers a versioned success route.
func (s MockHTTPServer) ServeVersionedSuccess(t *testing.T, path string, response interface{}) {
	s.serveVersion(t, path, 200, response)
}

// ServeVersionedFailure registers an unversioned failure route.
func (s MockHTTPServer) ServeVersionedFailure(t *testing.T, path string, response interface{}) {
	s.serveVersion(t, path, 500, response)
}

// serveUnversioned registers the unversioned route.
func (s MockHTTPServer) serveUnversion(t *testing.T, path string, statusCode int, response interface{}) {
	s.serve(t, "/synse", path, statusCode, response)
}

// serveVersioned registers the versioned route.
func (s MockHTTPServer) serveVersion(t *testing.T, path string, statusCode int, response interface{}) {
	s.serve(t, "/synse/v3", path, statusCode, response)

	s.mux.HandleFunc(
		"/synse/version",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fprint(t, w, `{"version": "3.x.x", "api_version": "v3"}`)
		})
}

// serve registers a path handler and writes to its response.
func (s MockHTTPServer) serve(t *testing.T, version string, path string, statusCode int, response interface{}) {
	s.mux.HandleFunc(
		fmt.Sprintf("%v%v", version, path),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			fprint(t, w, response)
		},
	)
}

// Close closes the server connection.
func (s MockHTTPServer) Close() {
	s.Server.Close()
}

// fprint calls fmt.Fprint and validates its returned error.
func fprint(t *testing.T, w io.Writer, a interface{}) {
	_, err := fmt.Fprint(w, a)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
}
