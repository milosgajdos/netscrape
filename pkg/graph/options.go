package graph

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	// DefaultWeight is default edge weight
	DefaultWeight = 1.0
)

// DOTOptions are DOT graph options.
type DOTOptions struct {
	GraphAttrs attrs.DOT
	NodeAttrs  attrs.DOT
	EdgeAttrs  attrs.DOT
}

// DOTOption configures DOT graph.
type DOTOption func(*Options)

// Options are graph options.
type Options struct {
	UID        uuid.UID
	DOTID      string
	Attrs      attrs.Attrs
	Weight     float64
	DOTOptions DOTOptions
}

// Option configures Options.
type Option func(*Options)

// WithUID sets UID Options.
func WithUID(u uuid.UID) Option {
	return func(o *Options) {
		o.UID = u
	}
}

// WithDOTID sets DOTID Options.
func WithDOTID(dotid string) Option {
	return func(o *Options) {
		o.DOTID = dotid
	}
}

// WithAttrs sets Attrs options
func WithAttrs(a attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}

// WithWeight sets Weight options.
func WithWeight(w float64) Option {
	return func(o *Options) {
		o.Weight = w
	}
}

// WithDOTOptions sets DOTOptions
func WithDOTOptions(do DOTOptions) Option {
	return func(o *Options) {
		o.DOTOptions = do
	}
}
