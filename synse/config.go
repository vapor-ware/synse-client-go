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

// ServerOptions
type ServerOptions struct {
	Address string
	Timeout time.Duration
}

// RetryOptions
type RetryOptions struct {
	Count       int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}
