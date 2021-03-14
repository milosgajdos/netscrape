package memory

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

// Node is a memory Graph node.
type Node struct {
	graph.Entity
	id    int64
	dotid string
	attrs attrs.Attrs
}

// New returns a new Node.
func NewNode(id int64, e graph.Entity, opts ...graph.Option) (*Node, error) {
	nopts := graph.Options{}
	for _, apply := range opts {
		apply(&nopts)
	}

	dotid := nopts.DOTID
	if dotid == "" {
		if dotEnt, ok := e.(graph.DOTEntity); ok {
			dotid = dotEnt.DOTID()
		} else {
			dotid = e.UID().String()
			if dotEnt, ok := e.(graph.DOTer); ok {
				dotid = dotEnt.DOTID()
			}
		}
	}

	a := attrs.NewCopyFrom(e.Attrs())
	if nopts.Attrs != nil {
		for _, k := range nopts.Attrs.Keys() {
			a.Set(k, nopts.Attrs.Get(k))
		}
	}
	a.Set(attrs.DOTID, dotid)
	a.Set(attrs.Name, dotid)

	return &Node{
		Entity: e,
		id:     id,
		dotid:  dotid,
		attrs:  a,
	}, nil
}

// ID returns node ID.
func (n Node) ID() int64 {
	return n.id
}

// DOTID returns Graphviz DOT ID.
func (n Node) DOTID() string {
	return n.dotid
}

// SetDOTID sets Graphviz DOT ID.
func (n *Node) SetDOTID(id string) {
	n.attrs.Set(attrs.DOTID, id)
	n.attrs.Set(attrs.Name, id)
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
