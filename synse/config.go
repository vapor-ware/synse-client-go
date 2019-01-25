package synse

// config.go defines config options for the client.

import (
	"time"
)

// Options is the root config options.
type Options struct {
	// Address specifies the URL of Synse Server in the format `host[:port]`.
	Address string

	// Timeout specifies a time limit for a request.
	Timeout time.Duration

	// Retry specifies the options for retry mechanism.
	Retry RetryOptions
}

// RetryOptions is the config options for backoff retry mechanism. Its strategy
// is to increase retry intervals after each failed attempt, until some maximum
// value.
type RetryOptions struct {
	// Count specifies the number of retry attempts. Zero value means no retry.
	Count int

	// WaitTime specifies the wait time before retrying request. It is
	// increased after each attempt.
	// Default is 100 milliseconds.
	WaitTime time.Duration

	// MaxWaitTime specifies the maximum wait time, the cap, of all retry
	// requests that are made.
	// Default is 2 seconds
	MaxWaitTime time.Duration
}
