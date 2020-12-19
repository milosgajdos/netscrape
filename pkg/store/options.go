package store

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
)

// Options are store options.
type Options struct {
	Graph graph.Graph
}

// AddOptions are store add options.
type AddOptions struct {
	Attrs attrs.Attrs
}

// AddOption sets add options.
type AddOption func(*AddOptions)

// DelOptions are store delete options.
type DelOptions struct {
	Attrs attrs.Attrs
}

// DelOption sets delete options.
type DelOption func(*DelOptions)
