package graph

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	"gonum.org/v1/gonum/graph/encoding"
)

// Entity is graph entity.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Name returns name
	Name() string
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// DOTEntity
type DOTEntity interface {
	Entity
	// DOTID returns Graphviz DOT ID.
	DOTID() string
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
	Entity
}

// Edge is an edge between two Graph nodes.
type Edge interface {
	// UID returns unique ID.
	UID() uuid.UID
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

// SubGrapher returns subgraph of a graph.
type SubGrapher interface {
	// SubGraph returns the maximum subgraph of a graph
	// starting at node with given uid up to given depth.
	SubGraph(ctx context.Context, uid uuid.UID, depth int, opts ...Option) (Graph, error)
}

// Graph is a graph of entities.
type Graph interface {
	// UID returns graph uid.
	UID() uuid.UID
	// Node returns the node with given uid.
	Node(context.Context, uuid.UID) (Node, error)
	// Nodes returns all graph nodes.
	Nodes(context.Context) ([]Node, error)
	// Edge returns the edge between two nodes.
	Edge(ctx context.Context, from, to uuid.UID) (Edge, error)
	// From returns all directly reachable nodes from node with the given UID.
	From(context.Context, uuid.UID) ([]Node, error)
}
