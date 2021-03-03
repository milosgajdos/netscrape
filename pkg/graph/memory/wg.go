package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/uuid"

	gngraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/traverse"
)

// WG is weighted graph.
type WG struct {
	WeightedGraphBuilder
	// uid is the UID of the graph
	uid uuid.UID
	// dotid is graph DOTID
	dotid string
	// nodes maps graph nodes
	nodes map[string]graph.Node
	// dot are graph DOT options
	dot graph.DOTOptions
}

// NewWG creates a new weighted graph and returns it
func NewWG(wg WeightedGraphBuilder, opts ...graph.Option) (*WG, error) {
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

	return &WG{
		WeightedGraphBuilder: wg,
		uid:                  uid,
		dotid:                dotid,
		nodes:                make(map[string]graph.Node),
		dot:                  gopts.DOTOptions,
	}, nil
}

// UID returns graph UID
func (g WG) UID() uuid.UID {
	return g.uid
}

// NewNode creates a new graph node and returns it.
// NOTE: this is a convenience method that creates
// the new *memory.Node such that it makes sure its
// ID does not already exist in the graph.
func (g *WG) NewNode(ctx context.Context, ent graph.Entity, opts ...graph.Option) (graph.Node, error) {
	gnode := g.WeightedGraphBuilder.NewNode()

	node, err := NewNode(gnode.ID(), ent, opts...)
	if err != nil {
		return nil, err
	}

	if n, ok := g.nodes[node.UID().Value()]; ok {
		return n, nil
	}

	return node, nil
}

// AddNode adds node n to the graph.
// If the node already exists in the graph it's left untouched.
// It returns error if n is not memory.Node.
func (g *WG) AddNode(ctx context.Context, n graph.Node) error {
	gnode, ok := n.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	if node := g.WeightedGraphBuilder.Node(gnode.ID()); node != nil {
		if _, ok := g.nodes[n.UID().Value()]; ok {
			return nil
		}

		g.nodes[n.UID().Value()] = n

		return nil
	}

	g.WeightedGraphBuilder.AddNode(gnode)

	g.nodes[n.UID().Value()] = n

	return nil
}

// Node returns the node with the given ID if it exists
// in the graph, and error if it could not be found.
func (g *WG) Node(ctx context.Context, uid uuid.UID) (graph.Node, error) {
	if node, ok := g.nodes[uid.Value()]; ok {
		return node, nil
	}

	return nil, graph.ErrNodeNotFound
}

// Nodes returns all the nodes in the graph.
func (g *WG) Nodes(ctx context.Context) ([]graph.Node, error) {
	graphNodes := gngraph.NodesOf(g.WeightedGraphBuilder.Nodes())

	nodes := make([]graph.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// RemoveNode removes the node with the given uid from graph.
func (g *WG) RemoveNode(ctx context.Context, uid uuid.UID) error {
	node, ok := g.nodes[uid.Value()]
	if !ok {
		return nil
	}

	gnode, ok := node.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	g.WeightedGraphBuilder.RemoveNode(gnode.ID())

	delete(g.nodes, uid.Value())

	return nil
}

// Link creates a new link (edge) between from and to and returns it or it returns the existing edge.
// It returns error if either of the nodes does not exist in the graph.
func (g *WG) Link(ctx context.Context, from, to uuid.UID, opts ...graph.Option) (graph.Edge, error) {
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

// Edge returns edge between nodes with the given UIDs.
func (g *WG) Edge(ctx context.Context, uid, vid uuid.UID) (graph.Edge, error) {
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
func (g *WG) Edges(ctx context.Context) ([]graph.Edge, error) {
	wedges := g.WeightedGraphBuilder.WeightedEdges()

	graphEdges := gngraph.WeightedEdgesOf(wedges)

	edges := make([]graph.Edge, len(graphEdges))

	for i, e := range graphEdges {
		edges[i] = e.(*Edge)
	}

	return edges, nil
}

// From returns all directly reachable nodes from the node with the given uid.
func (g *WG) From(ctx context.Context, uid uuid.UID) ([]graph.Node, error) {
	node, ok := g.nodes[uid.Value()]
	if !ok {
		return nil, nil
	}

	gnodes := g.WeightedGraphBuilder.From(node.(*Node).ID())

	graphNodes := gngraph.NodesOf(gnodes)

	nodes := make([]graph.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// Unlink removes the link between from and to nodes.
func (g *WG) Unlink(ctx context.Context, from, to uuid.UID) error {
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
	g.WeightedGraphBuilder.RemoveEdge(f.(*Node).ID(), t.(*Node).ID())

	return nil
}

// SubGraph returns the subgraph of the graph rooted in the node with the given uid up to the given depth.
func (g *WG) SubGraph(ctx context.Context, uid uuid.UID, depth int, opts ...graph.Option) (graph.Graph, error) {
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

	_ = bfs.Walk(g.WeightedGraphBuilder, root.(*Node), func(n gngraph.Node, d int) bool {
		return d == depth
	})

	if sgErr != nil {
		return nil, sgErr
	}

	for id, node := range sgNodes {
		nodes := g.WeightedGraphBuilder.From(id)
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

// DOTID returns the graph DOT ID.
func (g WG) DOTID() string {
	return g.dotid
}

// DOTAttributers implements encoding.Attributer.
func (g *WG) DOTAttributers() (graph, node, edge encoding.Attributer) {
	graph = g.dot.GraphAttrs
	node = g.dot.NodeAttrs
	edge = g.dot.EdgeAttrs

	return graph, node, edge
}

// DOT returns the GrapViz dot representation of the graph.
func (g *WG) DOT() (string, error) {
	b, err := dot.Marshal(g.WeightedGraphBuilder, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("DOT marshal error: %w", err)
	}

	return string(b), nil
}
