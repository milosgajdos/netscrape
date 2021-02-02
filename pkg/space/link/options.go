package link

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Options configure Links
type Options struct {
	// UID options
	UID uuid.UID
	// Attrs options
	Attrs attrs.Attrs
	// Merge links
	Merge bool
}

// Option configures Options.
type Option func(*Options)

// WithUID sets UID Options.
func WithUID(u uuid.UID) Option {
	return func(o *Options) {
		o.UID = u
	}
}

// WithAttrs sets Attrs options
func WithAttrs(a attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}

// WithMerge set merge option
func WithMerge(m bool) Option {
	return func(o *Options) {
		o.Merge = m
	}
}
