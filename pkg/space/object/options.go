package object

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Options configure Object
type Options struct {
	// UID options
	UID uuid.UID
	// Attrs options
	Attrs attrs.Attrs
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
