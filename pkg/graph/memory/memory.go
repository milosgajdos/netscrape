package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"

	gngraph "gonum.org/v1/gonum/graph"
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
	graph.NodeAdder
	graph.NodeRemover
	graph.Edger
	graph.Linker
	graph.Unlinker
	graph.SubGrapher
}

// WeightEdger returns all of the graph weighted edges
type WeightEdger interface {
	WeightedEdges() gngraph.WeightedEdges
}

// WeightedGraphBuilder allows to build in-memory weighted graphs.
type WeightedGraphBuilder interface {
	WeightEdger
	gngraph.Weighted
	gngraph.WeightedBuilder
	gngraph.NodeRemover
	gngraph.EdgeRemover
}
