package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
)

// Querier queries graph entities.
// NOTE: this interface is a [temporary] hack!
type Querier interface {
	// Query the graph and return the results.
	Query(context.Context, query.Query) ([]graph.Entity, error)
}

// Graph is in-memory graph.
type Graph interface {
	graph.Graph
	graph.SubGrapher
	graph.NodeAdder
	graph.NodeRemover
	graph.Edger
	graph.Linker
	graph.Unlinker
}
