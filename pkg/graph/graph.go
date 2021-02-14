package graph

import (
	"context"

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
	// Edges returns all graph edges.
	Edges(context.Context) ([]Edge, error)
}
