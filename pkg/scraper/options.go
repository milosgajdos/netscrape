package scraper

import (
	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Options configure scraper.
type Options struct {
	Store  store.Store
	Broker broker.Broker
}

// Option is functional scraper option.
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
