package synse

// config.go defines config options for the client.

import (
	"time"
)

// Options is the root config options.
type Options struct {
	Server ServerOptions
	Retry  RetryOptions
}

// ServerOptions is the config options for server.
type ServerOptions struct {
	Address string
	Timeout time.Duration
}

// RetryOptions is the config options for backoff retry mechanism.
type RetryOptions struct {
	Count       int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}
