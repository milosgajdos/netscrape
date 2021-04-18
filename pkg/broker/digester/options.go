package digester

import "github.com/milosgajdos/netscrape/pkg/broker"

// Options configure digester.
type Options struct {
	Handler     broker.Handler
	Unmarshaler broker.Unmarshaler
}

// Option is functional digester option.
type Option func(*Options)

// WithHandler sets handler option.
func WithHandler(h broker.Handler) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

// WithMarshaler sets Marshaler option.
func WithMarshaler(m broker.Unmarshaler) Option {
	return func(o *Options) {
		o.Unmarshaler = m
	}
}
