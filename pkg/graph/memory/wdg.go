package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	gngraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

// WDG is a weighted directed graph.
type WDG struct {
	*simple.WeightedDirectedGraph
	// uid is the UID of the graph
	uid uuid.UID
	// dotid is graph DOTID
	dotid string
	// nodes maps graph nodes
	nodes map[string]graph.Node
	// dot are graph options
	dot graph.DOTOptions
}

// NewWDG creates a new weighted directed graph and returns it.
// If DOTID is not provided via options, it's set to graph UID.
func NewWDG(opts ...graph.Option) (*WDG, error) {
	gopts := graph.Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	uid := gopts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	dotid := gopts.DOTID
	if dotid == "" {
		dotid = uid.Value()
	}

	return &WDG{
		WeightedDirectedGraph: simple.NewWeightedDirectedGraph(gopts.Weight, gopts.Weight),
		uid:                   uid,
		dotid:                 dotid,
		nodes:                 make(map[string]graph.Node),
		dot:                   gopts.DOTOptions,
	}, nil
}

// UID returns graph UID
func (g WDG) UID() uuid.UID {
	return g.uid
}

// NewNode creates a new graph node and returns it.
// NOTE: this is a convenience method that creates
// the new *memory.Node such that it makes sure its
// ID does not already exist in the graph relieving you
// from the necessity to make sure your new Node.ID()
// returns an unique ID in the underlying graph.
func (g *WDG) NewNode(ctx context.Context, ent graph.Entity, opts ...graph.Option) (graph.Node, error) {
	gnode := g.WeightedDirectedGraph.NewNode()

	node, err := NewNode(gnode.ID(), ent, opts...)
	if err != nil {
		return nil, err
	}

	if n, ok := g.nodes[node.UID().Value()]; ok {
		return n, nil
	}

	return node, nil
}

// AddNode adds node to the graph.
// If the node already exists in the graph it's left untouched.
// It returns error if n is not memory.Node.
func (g *WDG) AddNode(ctx context.Context, n graph.Node) error {
	gnode, ok := n.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	if node := g.WeightedDirectedGraph.Node(gnode.ID()); node != nil {
		if _, ok := g.nodes[n.UID().Value()]; ok {
			return nil
		}

		g.nodes[n.UID().Value()] = n

		return nil
	}

	g.WeightedDirectedGraph.AddNode(gnode)

	g.nodes[n.UID().Value()] = n

	return nil
}

// Node returns the node with the given ID if it exists
// in the graph, and error if it could not be retrieved.
func (g *WDG) Node(ctx context.Context, uid uuid.UID) (graph.Node, error) {
	if node, ok := g.nodes[uid.Value()]; ok {
		return node, nil
	}

	return nil, graph.ErrNodeNotFound
}

// Nodes returns all the nodes in the graph.
func (g *WDG) Nodes(ctx context.Context) ([]graph.Node, error) {
	graphNodes := gngraph.NodesOf(g.WeightedDirectedGraph.Nodes())

	nodes := make([]graph.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// RemoveNode removes the node with the given uid from graph.
func (g *WDG) RemoveNode(ctx context.Context, uid uuid.UID) error {
	node, ok := g.nodes[uid.Value()]
	if !ok {
		return nil
	}

	gnode, ok := node.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	g.WeightedDirectedGraph.RemoveNode(gnode.ID())

	delete(g.nodes, uid.Value())

	return nil
}

// Link creates a new edge between from and to and returns it or it returns the existing edge.
// It returns error if either of the nodes does not exist in the graph.
func (g *WDG) Link(ctx context.Context, from, to uuid.UID, opts ...graph.Option) (graph.Edge, error) {
	e, err := g.Edge(ctx, from, to)
	if err != nil && !errors.Is(err, graph.ErrEdgeNotExist) {
		return nil, err
	}

	if e != nil {
		return e, nil
	}

	f, ok := g.nodes[from.Value()]
	if !ok {
		return nil, fmt.Errorf("node link %s: %w", from, graph.ErrNodeNotFound)
	}

	t, ok := g.nodes[to.Value()]
	if !ok {
		return nil, fmt.Errorf("node link %s: %w", to, graph.ErrNodeNotFound)
	}

	edge, err := NewEdge(f.(*Node), t.(*Node), opts...)
	if err != nil {
		return nil, err
	}

	g.SetWeightedEdge(edge)

	return edge, nil
}

// Edge returns edge between the two nodes
func (g *WDG) Edge(ctx context.Context, uid, vid uuid.UID) (graph.Edge, error) {
	from, ok := g.nodes[uid.Value()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", uid, graph.ErrNodeNotFound)
	}

	to, ok := g.nodes[vid.Value()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", vid, graph.ErrNodeNotFound)
	}

	// NOTE: it's safe to switch type without checking if possible
	// the nodes in g.nodes are guaranteed to be *Nodes since the only
	// way to add nodes in is via AddNode which does type assertion
	if e := g.WeightedEdge(from.(*Node).ID(), to.(*Node).ID()); e != nil {
		return e.(*Edge), nil
	}

	return nil, fmt.Errorf("%w", graph.ErrEdgeNotExist)
}

// Edges returns all the edges (lines) from u to v.
func (g *WDG) Edges(ctx context.Context) ([]graph.Edge, error) {
	wedges := g.WeightedDirectedGraph.Edges()

	graphEdges := gngraph.EdgesOf(wedges)

	edges := make([]graph.Edge, len(graphEdges))

	for i, e := range graphEdges {
		edges[i] = e.(*Edge)
	}

	return edges, nil
}

// From returns all directly reachable nodes from the node with the given uid.
func (g *WDG) From(ctx context.Context, uid uuid.UID) ([]graph.Node, error) {
	node, ok := g.nodes[uid.Value()]
	if !ok {
		return nil, nil
	}

	gnodes := g.WeightedDirectedGraph.From(node.(*Node).ID())

	graphNodes := gngraph.NodesOf(gnodes)

	nodes := make([]graph.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// Unlink removes the link between from and to nodes.
func (g *WDG) Unlink(ctx context.Context, from, to uuid.UID) error {
	f, ok := g.nodes[from.Value()]
	if !ok {
		return nil
	}

	t, ok := g.nodes[to.Value()]
	if !ok {
		return nil
	}

	// NOTE: it's safe to switch type without checking if possible
	// the nodes in g.nodes are guaranteed to be *Nodes since the only
	// way to add nodes in is via AddNode which does type assertion
	g.WeightedDirectedGraph.RemoveEdge(f.(*Node).ID(), t.(*Node).ID())

	return nil
}

// SubGraph returns the subgraph of the graph rooted in the node with the given uid up to the given depth.
func (g *WDG) SubGraph(ctx context.Context, uid uuid.UID, depth int, opts ...graph.Option) (graph.Graph, error) {
	root, ok := g.nodes[uid.Value()]
	if !ok {
		return nil, graph.ErrNodeNotFound
	}

	sg, err := NewWDG(opts...)
	if err != nil {
		return nil, err
	}

	var sgErr error

	sgNodes := make(map[int64]graph.Node)

	visit := func(n gngraph.Node) {
		vnode := n.(*Node)

		if err := sg.AddNode(ctx, vnode); err != nil {
			sgErr = err
			return
		}

		sgNodes[vnode.ID()] = vnode
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	_ = bfs.Walk(g.WeightedDirectedGraph, root.(*Node), func(n gngraph.Node, d int) bool {
		return d == depth
	})

	if sgErr != nil {
		return nil, sgErr
	}

	for id, node := range sgNodes {
		nodes := g.WeightedDirectedGraph.From(id)
		for nodes.Next() {
			pnode := nodes.Node()
			peer := pnode.(*Node)
			if to, ok := sgNodes[peer.ID()]; ok {
				if edges := g.WeightedEdges(); edges != nil {
					for edges.Next() {
						we := edges.WeightedEdge()
						e := we.(*Edge)

						a := attrs.NewCopyFrom(e.Attrs())

						opts := []graph.Option{
							graph.WithAttrs(a),
							graph.WithWeight(e.Weight()),
						}

						if _, err := sg.Link(ctx, node.UID(), to.UID(), opts...); err != nil {
							return nil, fmt.Errorf("subgraph %s link error: %v", sg.UID(), err)
						}
					}
				}
			}
		}
	}

	return sg, nil
}

// Query queries the in-memory graph and returns the matched results.
// Query does a linear search over all of the graph nodes.
// TODO: make this more efficient
func (g WDG) Query(ctx context.Context, q query.Query) ([]graph.Entity, error) {
	if m := q.Matcher(query.UID); m != nil {
		if uid, ok := m.Predicate().Value().(uuid.UID); ok && len(uid.Value()) > 0 {
			if n, ok := g.nodes[uid.Value()]; ok {
				return []graph.Entity{n.(graph.Entity)}, nil
			}
			return []graph.Entity{}, nil
		}
	}

	var results []graph.Entity

	/*// Linear search over topologically sorted nodes
	sorted, err := topo.Sort(g.WeightedDirectedGraph)
	if err != nil {
		return nil, fmt.Errorf("unable to sort graph: %v", err)
	}

	for _, n := range sorted {
		node := n.(*Node)

		if m := q.Matcher(query.Attrs); m != nil {
			if !m.Match(node.Attrs()) {
				continue
			}
		}

		results = append(results, node)
	}
	*/

	graphNodes := gngraph.NodesOf(g.WeightedDirectedGraph.Nodes())

	for _, n := range graphNodes {
		node := n.(*Node)

		if m := q.Matcher(query.Attrs); m != nil {
			if !m.Match(node.Attrs()) {
				continue
			}
		}
		results = append(results, node)
	}

	return results, nil
}

// DOTID returns the graph DOT ID.
func (g WDG) DOTID() string {
	return g.dotid
}

// DOTAttributers implements encoding.Attributer.
func (g *WDG) DOTAttributers() (graph, node, edge encoding.Attributer) {
	graph = g.dot.GraphAttrs
	node = g.dot.NodeAttrs
	edge = g.dot.EdgeAttrs

	return graph, node, edge
}

// DOT returns the GrapViz dot representation of the graph.
func (g *WDG) DOT() (string, error) {
	b, err := dot.Marshal(g.WeightedDirectedGraph, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("DOT marshal error: %w", err)
	}

	return string(b), nil
}
