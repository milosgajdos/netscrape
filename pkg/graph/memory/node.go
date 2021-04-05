package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"gonum.org/v1/gonum/graph/encoding"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
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

	a := nopts.Attrs
	if a == nil {
		a = memattrs.New()
	}

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
	n.dotid = id
}

// Attrs returns node attributes.
func (n Node) Attrs() attrs.Attrs {
	return n.attrs
}

// Attributes implements attrs.DOT.
func (n Node) Attributes() []encoding.Attribute {
	keys, err := n.attrs.Keys(context.Background())
	if err != nil {
		return nil
	}

	attrs := make([]encoding.Attribute, len(keys))

	i := 0
	for _, k := range keys {
		val, err := n.attrs.Get(context.Background(), k)
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
