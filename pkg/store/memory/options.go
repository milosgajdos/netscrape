package memory

import (
	"github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Options are graph options.
type Options struct {
	UID   uuid.UID
	Graph memory.Graph
}

// Option configures Options.
type Option func(*Options)

// WithUID sets UID Options.
func WithUID(u uuid.UID) Option {
	return func(o *Options) {
		o.UID = u
	}
}

// WithGraph sets Graph options.
func WithGraph(g memory.Graph) Option {
	return func(o *Options) {
		o.Graph = g
	}
}
