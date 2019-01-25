package synse

// config.go defines config options for the client.

import (
	"time"
)

// Options is the root config options.
type Options struct {
	// Address specifies the URL of Synse Server in format host[:port]
	Address string

	// Timeout specifies a time limit for a request.
	Timeout time.Duration

	// Retry specifies the options for retry mechanism.
	Retry RetryOptions
}

// RetryOptions is the config options for backoff retry mechanism.
// TODO - add docs for all the fields.
type RetryOptions struct {
	Count       int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}
