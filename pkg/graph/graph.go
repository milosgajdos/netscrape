package graph

import (
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	"gonum.org/v1/gonum/graph/encoding"
)

// Entity is graph entity
type Entity interface {
	entity.Entity
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

// Node is a graph node.
type Node interface {
	Entity
	space.Object
}

// Edge is an edge between two graph nodes.
type Edge interface {
	Entity
	// FromNode returns the from node of the edge.
	FromNode() Node
	// ToNode returns the to node of the edge.
	ToNode() Node
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

// SubGrapher returns the maximal reachable subgraph of graph.
type SubGrapher interface {
	// SubGraph returns the max subgraph of graph
	// starting at node with given uid up to given depth.
	SubGraph(uid uuid.UID, depth int) (Graph, error)
}

// Querier queries graph.
type Querier interface {
	// Query the graph and return the results.
	Query(query.Query) ([]Entity, error)
}

// NodeAdder adds new Nodes to graph.
type NodeAdder interface {
	// NewNode returns a new Node.
	NewNode(space.Object, NodeOptions) (Node, error)
	// AddNode adds a new node to the graph.
	AddNode(Node) error
}

// NodeRemover removes node from Grapg
type NodeRemover interface {
	// RemoveNode removes node from the graph.
	RemoveNode(uid uuid.UID) error
}

// NodeLinker links arbitrary nodes in graph.
type NodeLinker interface {
	// Link links two nodes and returns the new edge.
	Link(from, to uuid.UID, opts LinkOptions) (Edge, error)
}

// LinkRemover removes link between two Nodes.
type LinkRemover interface {
	// RemoveEdge removes edge(s) from the graph.
	RemoveLink(from, to uuid.UID) error
}

// Graph is a graph of Space objects.
type Graph interface {
	// Node returns node with given uid.
	Node(uuid.UID) (Node, error)
	// Nodes returns all graph nodes.
	Nodes() ([]Node, error)
	// Edge returns the edge between the two nodes.
	Edge(from, to uuid.UID) (Edge, error)
	// Edges returns all graph edges.
	Edges() ([]Edge, error)
}
