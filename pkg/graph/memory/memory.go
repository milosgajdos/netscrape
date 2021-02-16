package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Querier queries graph entities.
// NOTE: this interface is a [temporary] hack!
type Querier interface {
	// Query the graph and return the results.
	Query(context.Context, query.Query) ([]graph.Entity, error)
}

// NodeAdder adds nodes to graph.
type NodeAdder interface {
	// NewNode returns a new node.
	NewNode(context.Context, space.Entity, ...graph.Option) (graph.Node, error)
	// AddNode adds a new node to graph.
	AddNode(context.Context, graph.Node) error
}

// NodeRemover removes node from graph
type NodeRemover interface {
	// RemoveNode removes node from graph.
	RemoveNode(context.Context, uuid.UID) error
}

// Linker links two nodes in graph.
type Linker interface {
	// Link links two nodes and returns the new edge.
	Link(ctx context.Context, from, to uuid.UID, opts ...graph.Option) (graph.Edge, error)
}

// Unlinker removes link between two Nodes.
type Unlinker interface {
	// Unlink removes the link between the nodes with given UIDs from graph.
	Unlink(ctx context.Context, from, to uuid.UID) error
}

// Graph is in-memory graph.
type Graph interface {
	graph.Graph
	graph.SubGrapher
	NodeAdder
	NodeRemover
	Linker
	Unlinker
}
