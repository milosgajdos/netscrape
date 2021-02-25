package memory

import (
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Grapher provides access to store graph.
type Grapher interface {
	// Graph returns graph handle
	Graph() (graph.Graph, error)
}

// Store is memory store
type Store interface {
	Grapher
	store.Store
}
