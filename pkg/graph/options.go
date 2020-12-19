package graph

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
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
	DOTOptions DOTOptions
	Weight     float64
}

// NodeOptions are graph node options.
type NodeOptions struct {
	Attrs attrs.Attrs
}

// NodeOption sets Node options.
type NodeOption func(*NodeOptions)

// LinkOptions are graph Link options.
type LinkOptions struct {
	Attrs  attrs.Attrs
	Weight float64
}

// LinkOption sets Link options.
type LinkOption func(*LinkOptions)
