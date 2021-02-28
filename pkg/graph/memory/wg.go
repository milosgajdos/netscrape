package memory

import (
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/uuid"
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
	// dot are graph options
	dot graph.DOTOptions
}

// NewWG creates a new weighted graph and returns it
func New(wg WeightedGraphBuilder, opts ...graph.Option) (*WG, error) {
	return nil, graph.ErrNotImplemented
}
