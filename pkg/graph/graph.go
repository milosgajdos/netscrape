package graph

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	"gonum.org/v1/gonum/graph/encoding"
)

// DOTer implements Graphviz DOT properties.
type DOTer interface {
	// DOTID returns Graphviz DOT ID.
	DOTID() string
	// SetDOTID sets Graphviz DOT ID.
	SetDOTID(string)
	// Attributes returns Graphviz DOT attributes
	Attributes() []encoding.Attribute
}

// Enitty is stored as a node in graph.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Type of entity.
	Type() entity.Type
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// DOTEntity is DOT Entity
type DOTEntity interface {
	Entity
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
	// UID returns unique ID.
	UID() uuid.UID
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Edge is an edge between two Graph nodes.
type Edge interface {
	// UID returns unique ID.
	UID() uuid.UID
	// FromNode returns the from node of the edge.
	FromNode() (Node, error)
	// ToNode returns the to node of the edge.
	ToNode() (Node, error)
	// Weight returns edge weight.
	Weight() float64
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// DOTGraph is a GraphViz DOT graph.
type DOTGraph interface {
	Graph
	// DOT returns Graphviz DOT graph.
	DOT() (string, error)
	// DOTID returns grapph DOT ID.
	DOTID() string
	// DOTAttributers returns graph DOT attributes.
	DOTAttributers() (graph, node, edge encoding.Attributer)
}

// SubGrapher returns subgraph of a graph.
type SubGrapher interface {
	// SubGraph returns the maximum subgraph of a graph
	// starting at node with given uid up to given depth.
	SubGraph(ctx context.Context, uid uuid.UID, depth int, opts ...Option) (Graph, error)
}

// NodeAdder adds nodes to graph.
type NodeAdder interface {
	// NewNode returns a new node.
	NewNode(context.Context, Entity, ...Option) (Node, error)
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
	// Unlink removes the link between the nodes with given UIDs from graph.
	Unlink(ctx context.Context, from, to uuid.UID) error
}

// Edger returns graph edges
type Edger interface {
	// Edges returns graph edges
	Edges(ctx context.Context) ([]Edge, error)
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
