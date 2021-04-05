package netscrape

import (
	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Options configure netscraping.
type Options struct {
	Store  store.Store
	Broker broker.Broker
}

// Option is functional netscrape option.
type Option func(*Options)

// WithStore sets Store options.
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithBroker sets Broker options.
func WithBroker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}
