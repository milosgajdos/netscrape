package store

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
)

// Options are store options.
type Options struct {
	Graph graph.Graph
	Attrs attrs.Attrs
}

// Option configures Options.
type Option func(*Options)

// WithGraph sets Graph options
func WithGraph(g graph.Graph) Option {
	return func(o *Options) {
		o.Graph = g
	}
}

// WithAttrs sets Attrs options
func WithAttrs(a attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}
