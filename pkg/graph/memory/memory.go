package memory

import "github.com/milosgajdos/netscrape/pkg/graph"

// Graph is in-memory graph.
type Graph interface {
	graph.Graph
	graph.SubGrapher
	graph.NodeAdder
	graph.NodeRemover
	graph.Linker
	graph.Unlinker
}
