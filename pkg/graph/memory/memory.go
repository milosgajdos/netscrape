package memory

import (
	"github.com/milosgajdos/netscrape/pkg/graph"

	gngraph "gonum.org/v1/gonum/graph"
)

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

// WeightEdger returns all of the graph weighted edges.
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
