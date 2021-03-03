package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	gngraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// WDG is a weighted directed graph.
type WDG struct {
	*WG
}

// NewWDG creates a new weighted directed graph and returns it.
// If DOTID is not provided via options, it's set to graph UID.
func NewWDG(opts ...graph.Option) (*WDG, error) {
	gopts := graph.Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	wg, err := NewWG(simple.NewWeightedDirectedGraph(gopts.Weight, gopts.Weight), opts...)
	if err != nil {
		return nil, err
	}

	return &WDG{
		WG: wg,
	}, nil
}

// Query queries the in-memory graph and returns the matched results.
// Query does a linear search over all of the graph nodes.
// TODO: make this more efficient i.e. faster
func (g WDG) Query(ctx context.Context, q query.Query) ([]graph.Entity, error) {
	if m := q.Matcher(query.UID); m != nil {
		if uid, ok := m.Predicate().Value().(uuid.UID); ok && len(uid.Value()) > 0 {
			if n, ok := g.nodes[uid.Value()]; ok {
				return []graph.Entity{n.(graph.Entity)}, nil
			}
			return []graph.Entity{}, nil
		}
	}

	graphNodes := gngraph.NodesOf(g.WG.WeightedGraphBuilder.Nodes())

	results := make([]graph.Entity, len(graphNodes))

	for i, n := range graphNodes {
		node := n.(*Node)

		if m := q.Matcher(query.Attrs); m != nil {
			if !m.Match(node.Attrs()) {
				continue
			}
		}
		results[i] = node
	}

	return results, nil
}
