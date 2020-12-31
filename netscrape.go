package netscrape

import (
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Filter filters scrapes Objects.
type Filter func(space.Object) bool

// NetScraper builds a network of scraped Objects.
type NetScraper interface {
	// Run runs netscaping and returns error if it fails
	Run(space.Scraper, space.Origin, ...Filter) error
	// Store returns store that stores Object graph.
	Store() store.Store
}

// Options are kraph options.
type Options struct{}

// Option is functional kraph option.
type Option func(*Options)
