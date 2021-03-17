package memory

import (
	"github.com/milosgajdos/netscrape/pkg/graph"
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
