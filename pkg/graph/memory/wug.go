package memory

import (
	"context"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	gngraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

// WUG is weighted undirected graph.
type WUG struct {
	*simple.WeightedUndirectedGraph
	// id is ID of the graph
	id string
	// nodes maps graph nodes
	nodes map[string]graph.Node
	// opts are graph options
	opts graph.Options
}

// NewWUG creates a new weighted undirected graph and returns it.
func NewWUG(id string, opts graph.Options) (*WUG, error) {
	return &WUG{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(opts.Weight, opts.Weight),
		id:                      id,
		nodes:                   make(map[string]graph.Node),
		opts:                    opts,
	}, nil
}

// NewNode creates a new graph node and returns it.
// NOTE: this is a convenience method which creates
// the new *memory.Node such that it makes sure its
// ID does not already exist in the graph relieving you
// from the necessity to make sure your new Node.ID()
// returns unique ID in the underlying graph.
func (g *WUG) NewNode(ctx context.Context, obj space.Object, opts graph.NodeOptions) (graph.Node, error) {
	gnode := g.WeightedUndirectedGraph.NewNode()

	node, err := NewNode(gnode.ID(), obj, opts)
	if err != nil {
		return nil, err
	}

	if n, ok := g.nodes[node.UID().Value()]; ok {
		return n, nil
	}

	return node, nil
}

// AddNode adds node to the graph.
func (g *WUG) AddNode(ctx context.Context, n graph.Node) error {
	if _, ok := g.nodes[n.UID().Value()]; ok {
		return nil
	}

	gnode, ok := n.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	if node := g.WeightedUndirectedGraph.Node(gnode.ID()); node != nil {
		g.nodes[n.UID().Value()] = n
		return nil
	}

	g.nodes[n.UID().Value()] = n

	g.WeightedUndirectedGraph.AddNode(gnode)

	return nil
}

// Node returns the node with the given ID if it exists
// in the graph, and error if it could not be retrieved.
func (g *WUG) Node(ctx context.Context, uid uuid.UID) (graph.Node, error) {
	if node, ok := g.nodes[uid.Value()]; ok {
		return node, nil
	}

	return nil, graph.ErrNodeNotFound
}

// Nodes returns all the nodes in the graph.
func (g *WUG) Nodes(ctx context.Context) ([]graph.Node, error) {
	graphNodes := gngraph.NodesOf(g.WeightedUndirectedGraph.Nodes())

	nodes := make([]graph.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// RemoveNode removes the node with the given id from graph.
func (g *WUG) RemoveNode(ctx context.Context, uid uuid.UID) error {
	node, ok := g.nodes[uid.Value()]
	if !ok {
		return nil
	}

	gnode, ok := node.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	g.WeightedUndirectedGraph.RemoveNode(gnode.ID())

	delete(g.nodes, uid.Value())

	return nil
}

// Link creates a new edge between from and to and returns it or it returns the existing edge.
// It returns error if either of the nodes does not exist in the graph.
func (g *WUG) Link(ctx context.Context, from, to uuid.UID, opts graph.LinkOptions) (graph.Edge, error) {
	e, err := g.Edge(ctx, from, to)
	if err != nil && err != graph.ErrEdgeNotExist {
		return nil, err
	}

	if e != nil {
		return e, nil
	}

	f, ok := g.nodes[from.Value()]
	if !ok {
		return nil, fmt.Errorf("node %s link error: %w", from, graph.ErrNodeNotFound)
	}

	t, ok := g.nodes[to.Value()]
	if !ok {
		return nil, fmt.Errorf("node %s link error: %w", to, graph.ErrNodeNotFound)
	}

	var entOpts []entity.Option

	attrs := attrs.NewCopyFrom(opts.Attrs)
	entOpts = append(entOpts, entity.Attrs(attrs))

	w := opts.Weight
	if opts.Weight < 0 {
		w = graph.DefaultWeight
	}

	edge, err := NewEdge(f.(*Node), t.(*Node), w, entOpts...)
	if err != nil {
		return nil, err
	}

	g.SetWeightedEdge(edge)

	return edge, nil
}

// Edge returns edge between the two nodes
func (g *WUG) Edge(ctx context.Context, uid, vid uuid.UID) (graph.Edge, error) {
	from, ok := g.nodes[uid.Value()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", uid, graph.ErrNodeNotFound)
	}

	to, ok := g.nodes[vid.Value()]
	if !ok {
		return nil, fmt.Errorf("%s: %w", vid, graph.ErrNodeNotFound)
	}

	// NOTE: it's safe to typecast without checking as
	// the nodes in g.nodes have *Node type since the only
	// way to add the nodes in is via AddNode which does type assertion
	if e := g.WeightedEdge(from.(*Node).ID(), to.(*Node).ID()); e != nil {
		return e.(*Edge), nil
	}

	return nil, graph.ErrEdgeNotExist
}

// Edges returns all the edges (lines) from u to v.
func (g *WUG) Edges(ctx context.Context) ([]graph.Edge, error) {
	wedges := g.WeightedUndirectedGraph.Edges()

	graphEdges := gngraph.EdgesOf(wedges)

	edges := make([]graph.Edge, len(graphEdges))

	for i, e := range graphEdges {
		edges[i] = e.(*Edge)
	}

	return edges, nil
}

// RemoveLink removes link between two nodes.
func (g *WUG) RemoveLink(ctx context.Context, from, to uuid.UID) error {
	f, ok := g.nodes[from.Value()]
	if !ok {
		return nil
	}

	t, ok := g.nodes[to.Value()]
	if !ok {
		return nil
	}

	// NOTE: it's safe to typecast without checking as
	// the nodes in g.nodes have *Node type since the only
	// way to add the nodes in is via AddNode which does type assertion
	g.WeightedUndirectedGraph.RemoveEdge(f.(*Node).ID(), t.(*Node).ID())

	return nil
}

// SubGraph returns the subgraph of the node up to given depth.
func (g *WUG) SubGraph(ctx context.Context, uid uuid.UID, depth int) (graph.Graph, error) {
	root, ok := g.nodes[uid.Value()]
	if !ok {
		return nil, graph.ErrNodeNotFound
	}

	sub, err := NewWUG(g.id+"subgraph", g.opts)
	if err != nil {
		return nil, err
	}

	var sgErr error

	subnodes := make(map[int64]graph.Node)

	visit := func(n gngraph.Node) {
		vnode := n.(*Node)

		if err := sub.AddNode(ctx, vnode); err != nil {
			sgErr = err
			return
		}
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	_ = bfs.Walk(g.WeightedUndirectedGraph, root.(*Node), func(n gngraph.Node, d int) bool {
		return d == depth
	})

	if sgErr != nil {
		return nil, sgErr
	}

	for id, node := range subnodes {
		nodes := sub.From(id)
		for nodes.Next() {
			pnode := nodes.Node()
			peer := pnode.(*Node)
			if to, ok := subnodes[peer.ID()]; ok {
				if edges := g.WeightedEdges(); edges != nil {
					for edges.Next() {
						we := edges.WeightedEdge()
						e := we.(*Edge)

						a := attrs.NewCopyFrom(e.Attrs())

						opts := graph.LinkOptions{
							Attrs:  a,
							Weight: e.Weight(),
						}

						if _, err := sub.Link(ctx, node.UID(), to.UID(), opts); err != nil {
							return nil, fmt.Errorf("subgraph %s link error: %v", sub.id, err)
						}
					}
				}
			}
		}
	}

	return sub, nil
}

// queryEdge returns all the edges that match given query
func (g WUG) queryEdge(q query.Query) ([]graph.Edge, error) {
	traversed := make(map[string]bool)

	var results []graph.Edge

	trav := func(e gngraph.Edge) bool {
		edge := e.(*Edge)

		if traversed[edge.UID().Value()] {
			return false
		}

		traversed[edge.UID().Value()] = true

		if m := q.Matcher(query.PWeight); m != nil {
			if !m.Match(edge.Weight()) {
				return false
			}
		}

		if m := q.Matcher(query.PAttrs); m != nil {
			if !m.Match(edge.Attrs()) {
				return false
			}
		}

		results = append(results, edge)

		return true
	}

	dfs := traverse.DepthFirst{
		Traverse: trav,
	}

	dfs.WalkAll(g.WeightedUndirectedGraph, nil, nil, func(gngraph.Node) {})

	return results, nil
}

// queryNode returns all the nodes that match the given query.
func (g WUG) queryNode(q query.Query) ([]graph.Node, error) {
	if m := q.Matcher(query.PUID); m != nil {
		if uid, ok := m.Predicate().Value().(uuid.UID); ok && len(uid.Value()) > 0 {
			if n, ok := g.nodes[uid.Value()]; ok {
				return []graph.Node{n}, nil
			}
		}
	}

	var results []graph.Node

	visit := func(n gngraph.Node) {
		node := n.(*Node)

		if m := q.Matcher(query.PNamespace); m != nil {
			if !m.Match(node.Namespace()) {
				return
			}
		}

		if m := q.Matcher(query.PKind); m != nil {
			if !m.Match(node.Resource().Kind()) {
				return
			}
		}

		if m := q.Matcher(query.PName); m != nil {
			if !m.Match(node.Name()) {
				return
			}
		}

		if m := q.Matcher(query.PAttrs); m != nil {
			if !m.Match(node.Attrs()) {
				return
			}
		}

		results = append(results, node)
	}

	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	dfs.WalkAll(g.WeightedUndirectedGraph, nil, nil, func(gngraph.Node) {})

	return results, nil
}

// Query queries the in-memory graph and returns the matched results.
func (g WUG) Query(ctx context.Context, q query.Query) ([]graph.Entity, error) {
	var e query.EntityVal

	if m := q.Matcher(query.PEntity); m != nil {
		var ok bool
		e, ok = m.Predicate().Value().(query.EntityVal)
		if !ok {
			return nil, graph.ErrMissingEntity
		}
	}

	var entities []graph.Entity

	switch e {
	case query.Node:
		nodes, err := g.queryNode(q)
		if err != nil {
			return nil, fmt.Errorf("node query: %w", err)
		}

		for _, node := range nodes {
			entities = append(entities, node)
		}
	case query.Edge:
		edges, err := g.queryEdge(q)
		if err != nil {
			return nil, fmt.Errorf("edge query: %w", err)
		}

		for _, edge := range edges {
			entities = append(entities, edge)
		}
	default:
		return nil, graph.ErrUnknownEntity
	}

	return entities, nil
}

// DOTID returns the store DOT ID.
func (g WUG) DOTID() string {
	return g.id
}

// DOTAttributers implements encoding.Attributer.
func (g *WUG) DOTAttributers() (graph, node, edge encoding.Attributer) {
	graph = g.opts.DOTOptions.GraphAttrs
	node = g.opts.DOTOptions.NodeAttrs
	edge = g.opts.DOTOptions.EdgeAttrs

	return graph, node, edge
}

// DOT returns the GrapViz dot representation of netscrape.
func (g *WUG) DOT() (string, error) {
	b, err := dot.Marshal(g.WeightedUndirectedGraph, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("DOT marshal error: %w", err)
	}

	return string(b), nil
}
