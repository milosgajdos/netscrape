package scraper

import (
	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/cache"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Options configure scraper.
type Options struct {
	Store  store.Store
	Broker broker.Broker
	Links  cache.Links
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

// WithLinksCache sets Links options.
func WithLinksCache(l cache.Links) Option {
	return func(o *Options) {
		o.Links = l
	}
}
