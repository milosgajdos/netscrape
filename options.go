package netscrape

import (
	"github.com/milosgajdos/netscrape/pkg/cache"
	"github.com/milosgajdos/netscrape/pkg/plan"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Options are kraph options.
type Options struct {
	Store store.Store
	Plan  plan.Plan
	Links cache.Links
}

// Option is functional kraph option.
type Option func(*Options)

// WithStore sets Store Options
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithPlan configures Plan Options.
func WithPlan(p plan.Plan) Option {
	return func(o *Options) {
		o.Plan = p
	}
}

// WithLinksCache sets Cache options
func WithLinksCache(c cache.Links) Option {
	return func(o *Options) {
		o.Links = c
	}
}
