package memory

import (
	"github.com/milosgajdos/netscrape/pkg/graph"
	"gonum.org/v1/gonum/graph/simple"
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
