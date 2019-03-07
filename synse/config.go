package synse

// config.go defines config options for the client.

import (
	"time"
)

// Options is the root config options.
type Options struct {
	// Address specifies the URL of Synse Server in the format `host[:port]`.
	Address string `default:"-"`

	// HTTP speciifies the options for http protocol, used by a http client.
	HTTP HTTPOptions

	// WebSocket speciifies the options for websocket protocol, used by a
	// websocket client.
	WebSocket WebSocketOptions

	// TLS sepcifies the options for TLS/SSL communication.
	TLS TLSOptions
}

// HTTPOptions is the config options for http protocol,
type HTTPOptions struct {
	// Timeout specifies a time limit for a http request.
	Timeout time.Duration `default:"2s"`

	// Retry specifies the options for retry mechanism.
	Retry RetryOptions
}

// WebSocketOptions is the config options for websocket protocol.
type WebSocketOptions struct {
	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration `default:"45s"`
}

// RetryOptions is the config options for backoff retry mechanism. Its strategy
// is to increase retry intervals after each failed attempt, until some maximum
// value.
type RetryOptions struct {
	// Count specifies the number of retry attempts. Zero value means no retry.
	Count int `default:"3"`

	// WaitTime specifies the wait time before retrying request. It is
	// increased after each attempt.
	WaitTime time.Duration `default:"100ms"`

	// MaxWaitTime specifies the maximum wait time, the cap, of all retry
	// requests that are made.
	MaxWaitTime time.Duration `default:"2s"`
}

// TLSOptions is the config options for TLS/SSL communication.
type TLSOptions struct {
	// CertFile and KeyFile are public/private key pair from a pair of files to
	// use when communicating with Synse Server.
	CertFile string `default:"-"`
	KeyFile  string `default:"-"`

	// Enabled specifies whethere tls is enabled.
	Enabled bool `default:"false"`

	// SkipVerify specifies whether the client can skip certificate checks.
	SkipVerify bool `default:"false"`
}
