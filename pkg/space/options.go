package space

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/matcher"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Options are space options.
type Options struct {
	// UID options
	UID uuid.UID
	// Attrs options
	Attrs attrs.Attrs
	// Merge options
	Merge bool
	// Matcher options.
	Filter matcher.Matcher
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

// WithMerge sets Merge Options.
func WithMerge() Option {
	return func(o *Options) {
		o.Merge = true
	}
}

// WithFilter sets Filter Options.
func WithFilter(m matcher.Matcher) Option {
	return func(o *Options) {
		o.Filter = m
	}
}
