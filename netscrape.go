package netscrape

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Filter filters scrapes Objects.
type Filter func(space.Object) bool

// NetScraper builds a network of scraped Objects.
type NetScraper interface {
	// Run runs netscaping and returns error if it fails
	Run(context.Context, space.Scraper, space.Origin, ...Filter) error
	// Store returns store that stores Object graph.
	Store() store.Store
}

// Options are kraph options.
type Options struct {
	Store store.Store
}

// Option is functional kraph option.
type Option func(*Options)

// WithStore sets Store Options
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}
