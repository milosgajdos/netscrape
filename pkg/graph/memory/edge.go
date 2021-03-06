package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	"gonum.org/v1/gonum/graph/encoding"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
	gngraph "gonum.org/v1/gonum/graph"
)

// Edge implements graph.WeightedEdge
type Edge struct {
	uid    uuid.UID
	attrs  attrs.Attrs
	from   *Node
	to     *Node
	dotid  string
	weight float64
}

// NewEdge creates new Edge and returns it.
// If no DOTID is given, DOTID returns UID.
// If both WithUID and WithDOTID are provided
// DOTID overrides the UID value.
func NewEdge(from, to *Node, opts ...graph.Option) (*Edge, error) {
	eopts := graph.Options{}
	for _, apply := range opts {
		apply(&eopts)
	}

	uid := eopts.UID
	if uid == nil {
		uid = memuid.New()
	}

	dotid := eopts.DOTID
	if dotid == "" {
		dotid = uid.String()
	} else {
		uid = memuid.NewFromString(dotid)
	}

	a := eopts.Attrs
	if a == nil {
		a = memattrs.New()
	}

	return &Edge{
		uid:    uid,
		attrs:  a,
		from:   from,
		to:     to,
		dotid:  dotid,
		weight: eopts.Weight,
	}, nil
}

// UID returns Edge UID
func (e Edge) UID() uuid.UID {
	return e.uid
}

// Attrs returns edge attributes
func (e Edge) Attrs() attrs.Attrs {
	return e.attrs
}

// From returns the from node of the first non-nil edge, or nil.
func (e *Edge) From() gngraph.Node {
	return e.from
}

// To returns the to node of the first non-nil edge, or nil.
func (e *Edge) To() gngraph.Node {
	return e.to
}

// ReversedEdge returns a new edge with end points of the pair swapped.
func (e *Edge) ReversedEdge() gngraph.Edge {
	return &Edge{
		uid:    e.uid,
		attrs:  e.attrs,
		from:   e.to,
		to:     e.from,
		dotid:  e.dotid,
		weight: e.weight,
	}
}

// FromNode returns the from node of the first non-nil edge, or nil.
func (e Edge) FromNode() (graph.Node, error) {
	return e.from, nil
}

// ToNode returns the to node of the first non-nil edge, or nil.
func (e Edge) ToNode() (graph.Node, error) {
	return e.to, nil
}

// Weight returns edge weight
func (e Edge) Weight() float64 {
	return e.weight
}

// DOTID returns the edge's DOT ID.
func (e Edge) DOTID() string {
	return e.dotid
}

// SetDOTID sets the edge's DOT ID.
func (e *Edge) SetDOTID(id string) {
	e.dotid = id
}

// Attributes implements store.DOTAttrs
func (e Edge) Attributes() []encoding.Attribute {
	keys, err := e.attrs.Keys(context.Background())
	if err != nil {
		return nil
	}

	attrs := make([]encoding.Attribute, len(keys))

	i := 0
	for _, k := range keys {
		val, err := e.attrs.Get(context.Background(), k)
		if err != nil {
			return nil
		}
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: val,
		}
		i++
	}

	return attrs
}
