package graph

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	"gonum.org/v1/gonum/graph/encoding"
)

// Object is graph object.
type Object interface {
	space.Object
}

// Entity is graph entity.
type Entity interface {
	space.Entity
}

// DOTer implements GraphViz DOT properties.
type DOTer interface {
	// DOTID returns Graphviz DOT ID.
	DOTID() string
	// SetDOTID sets Graphviz DOT ID.
	SetDOTID(string)
}

// DOTNode is a GraphViz DOT Node.
type DOTNode interface {
	DOTer
	Node
}

// DOTEdge is a GraphViz DOT Edge.
type DOTEdge interface {
	DOTer
	Edge
}

// Node is a Graph node.
type Node interface {
	Object
	Entity
}

// Edge is an edge between two Graph nodes.
type Edge interface {
	Object
	// FromNode returns the from node of the edge.
	FromNode(context.Context) (Node, error)
	// ToNode returns the to node of the edge.
	ToNode(context.Context) (Node, error)
	// Weight returns edge weight.
	Weight() float64
}

// DOTGraph returns GraphViz DOT graph.
type DOTGraph interface {
	Graph
	// DOTID returns grapph DOT ID.
	DOTID() string
	// DOTAttributers returns graph DOT attributes.
	DOTAttributers() (graph, node, edge encoding.Attributer)
	// DOT returns Graphviz DOT graph.
	DOT() (string, error)
}

// SubGrapher returns the maximum subgraph of a graph.
type SubGrapher interface {
	// SubGraph returns the maximum subgraph of a graph
	// starting at node with given uid up to given depth.
	SubGraph(ctx context.Context, uid uuid.UID, depth int, opts ...Option) (Graph, error)
}

// Querier queries graph.
type Querier interface {
	// Query the graph and return the results.
	Query(context.Context, query.Query) ([]Object, error)
}

// NodeAdder adds nodes to graph.
type NodeAdder interface {
	// NewNode returns a new node.
	NewNode(context.Context, space.Entity, ...Option) (Node, error)
	// AddNode adds a new node to graph.
	AddNode(context.Context, Node) error
}

// NodeRemover removes node from graph
type NodeRemover interface {
	// RemoveNode removes node from graph.
	RemoveNode(context.Context, uuid.UID) error
}

// Linker links two nodes in graph.
type Linker interface {
	// Link links two nodes and returns the new edge.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) (Edge, error)
}

// Unlinker removes link between two Nodes.
type Unlinker interface {
	// Unlink removes link from graph.
	Unlink(ctx context.Context, from, to uuid.UID) error
}

// Graph is a graph of entities.
type Graph interface {
	// UID returns graph uid.
	UID() uuid.UID
	// Node returns node with given uid.
	Node(context.Context, uuid.UID) (Node, error)
	// Nodes returns all graph nodes.
	Nodes(context.Context) ([]Node, error)
	// Edge returns the edge between two nodes.
	Edge(ctx context.Context, from, to uuid.UID) (Edge, error)
	// Edges returns all graph edges.
	Edges(context.Context) ([]Edge, error)
}
