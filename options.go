package netscrape

import (
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Filter filters scrapes Objects.
type Filter func(space.Object) bool

// Options are kraph options.
type Options struct {
	Store   store.Store
	Filters []Filter
}

// Option is functional kraph option.
type Option func(*Options)

// WithStore sets Store Options
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithFilters set Filters options
func WithFilters(fx ...Filter) Option {
	return func(o *Options) {
		o.Filters = fx
	}
}
