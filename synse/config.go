package synse

// config.go defines config options for the client.

import (
	"time"
)

// Options is the root config options.
type Options struct {
	// Address specifies the URL of Synse Server in the format `host[:port]`.
	Address string `default:"-"`

	// HTTP specifies the options for http protocol, used by a http client.
	HTTP HTTPOptions

	// WebSocket specifies the options for websocket protocol, used by a
	// websocket client.
	WebSocket WebSocketOptions

	// TLS specifies the options for TLS/SSL communication.
	TLS TLSOptions
}

// HTTPOptions is the config options for http protocol,
type HTTPOptions struct {
	// Timeout specifies a time limit for a http request.
	Timeout time.Duration `default:"2s"`

	// Retry specifies the options for retry mechanism.
	Retry RetryOptions

	// Redirects specifies the max number of allowed http redirects
	Redirects int `default:"5"`
}

// WebSocketOptions is the config options for websocket protocol.
type WebSocketOptions struct {
	// HandshakeTimeout specifies the duration for the handshake to complete.
	// FIXME - note that the DefaultDialer in the gorilla/websocket pkg that
	// we're using that specifies the default handshake timeout to be 45s. I
	// don't have a sense on what is a good value either so I just use what
	// they have there.
	HandshakeTimeout time.Duration `default:"45s"`
}

// RetryOptions is the config options for backoff retry mechanism. Its strategy
// is to increase retry intervals after each failed attempt, until some maximum
// value.
type RetryOptions struct {
	// Count specifies the number of retry attempts. Zero value means no retry.
	Count uint `default:"3"`

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

	// Enabled specifies whether tls is enabled.
	Enabled bool `default:"false"`

	// SkipVerify specifies whether the client can skip certificate check. If
	// it is set to true, TLS will accept any certificate presented. However,
	// due to security concern, this should only be used for testing.
	SkipVerify bool `default:"false"`
}
