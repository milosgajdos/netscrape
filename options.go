package netscrape

import (
	"github.com/milosgajdos/netscrape/pkg/store"
)

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
