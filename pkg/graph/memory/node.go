package memory

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/space"
	"gonum.org/v1/gonum/graph/encoding"
)

// Node is a graph node.
type Node struct {
	space.Entity
	id    int64
	dotid string
	attrs attrs.Attrs
}

// NewNode creates new Node and returns it.
// If WithAttrs is passed in, its values override e attributes.
func NewNode(id int64, e space.Entity, opts ...graph.Option) (*Node, error) {
	nopts := graph.Options{}
	for _, apply := range opts {
		apply(&nopts)
	}

	dotid := nopts.DOTID
	if dotid == "" {
		dotid = e.UID().Value()
		if dotEnt, ok := e.(graph.DOTEntity); ok {
			dotid = dotEnt.DOTID()
		}
	}

	attrs := attrs.NewCopyFrom(e.Attrs())
	if nopts.Attrs != nil {
		for _, k := range nopts.Attrs.Keys() {
			attrs.Set(k, nopts.Attrs.Get(k))
		}
	}
	attrs.Set("dotid", dotid)
	attrs.Set("name", dotid)

	return &Node{
		Entity: e,
		id:     id,
		dotid:  dotid,
		attrs:  attrs,
	}, nil
}

// ID returns node ID.
func (n Node) ID() int64 {
	return n.id
}

// DOTID returns GraphViz DOT ID.
func (n Node) DOTID() string {
	return n.dotid
}

// SetDOTID sets GraphViz DOT ID.
// It sets both dotid and name attributes to id.
func (n *Node) SetDOTID(id string) {
	n.attrs.Set("dotid", id)
	n.attrs.Set("name", id)
	n.dotid = id
}

// Attrs returns node attributes.
func (n Node) Attrs() attrs.Attrs {
	return n.attrs
}

// Attributes implements attrs.DOT.
func (n Node) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, len(n.attrs.Keys()))

	i := 0
	for _, k := range n.attrs.Keys() {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: n.attrs.Get(k),
		}
		i++
	}

	return attrs
}
