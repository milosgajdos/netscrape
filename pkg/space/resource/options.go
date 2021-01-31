package resource

import "github.com/milosgajdos/netscrape/pkg/attrs"

// Options are Space options.
type Options struct {
	// Attrs options
	Attrs attrs.Attrs
}

// Option configures Options.
type Option func(*Options)

// WithAttrs sets Attrs options
func WithAttrs(a attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}
