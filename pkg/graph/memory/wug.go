package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	gngraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

// WUG is a weighted undirected graph.
type WUG struct {
	*WG
}

// NewWUG creates a new weighted undirected graph and returns it.
// If DOTID is not provided via options, it's set to graph UID.
func NewWUG(opts ...graph.Option) (*WUG, error) {
	gopts := graph.Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	wg, err := NewWG(simple.NewWeightedUndirectedGraph(gopts.Weight, gopts.Weight), opts...)
	if err != nil {
		return nil, err
	}

	return &WUG{
		WG: wg,
	}, nil
}

// Query queries the in-memory graph and returns the matched results.
func (g WUG) Query(ctx context.Context, q query.Query) ([]graph.Entity, error) {
	if m := q.Matcher(query.UID); m != nil {
		if uid, ok := m.Predicate().Value().(uuid.UID); ok && len(uid.Value()) > 0 {
			if n, ok := g.nodes[uid.Value()]; ok {
				return []graph.Entity{n.(graph.Entity)}, nil
			}
			return []graph.Entity{}, nil
		}
	}

	var results []graph.Entity

	visit := func(n gngraph.Node) {
		node := n.(*Node)

		if m := q.Matcher(query.Attrs); m != nil {
			if !m.Match(node.Attrs()) {
				return
			}
		}

		results = append(results, node)
	}

	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	dfs.WalkAll(g.WG.WeightedGraphBuilder.(gngraph.Undirected), nil, nil, func(gngraph.Node) {})

	return results, nil
}
