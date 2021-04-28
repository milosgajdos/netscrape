package ingester

import "github.com/milosgajdos/netscrape/pkg/broker"

// Options configure ingester.
type Options struct {
	Marshaler broker.Marshaler
}

// Option is functional ingester option.
type Option func(*Options)

// WithMarshaler sets Marshaler option.
func WithMarshaler(m broker.Marshaler) Option {
	return func(o *Options) {
		o.Marshaler = m
	}
}
