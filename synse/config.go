package synse

// config.go defines config options for the client.

import (
	"time"
)

// Options is the root config options.
type Options struct {
	// Address specifies the URL of Synse Server in the format `host[:port]`.
	Address string `default:"-"`

	// Timeout specifies a time limit for a request.
	Timeout time.Duration `default:"2s"`

	// Retry specifies the options for retry mechanism.
	Retry RetryOptions
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
