package resource

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Options are Space options.
type Options struct {
	// UID options
	UID uuid.UID
	// Attrs options
	Attrs attrs.Attrs
	// DOTID options
	DOTID string
}

// Option configures Options.
type Option func(*Options)

// WithAttrs sets Attrs options
func WithAttrs(a attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}

// WithUID sets UID Options.
func WithUID(u uuid.UID) Option {
	return func(o *Options) {
		o.UID = u
	}
}

// WithDOTID sets Attrs options
func WithDOTID(d string) Option {
	return func(o *Options) {
		o.DOTID = d
	}
}
