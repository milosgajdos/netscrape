package memory

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/space"
	"gonum.org/v1/gonum/graph/encoding"
)

// Node is a graph node.
type Node struct {
	space.Object
	id    int64
	dotid string
	attrs attrs.Attrs
}

// NewNodeWithDOTID creates a new Node with the given DOTID and returns it.
func NewNodeWithDOTID(id int64, obj space.Object, dotid string, opts graph.NodeOptions) (*Node, error) {

	attrs := attrs.NewCopyFrom(opts.Attrs)

	return &Node{
		Object: obj,
		id:     id,
		dotid:  dotid,
		attrs:  attrs,
	}, nil
}

// NewNode creates a new Node and returns it.
func NewNode(id int64, obj space.Object, opts graph.NodeOptions) (*Node, error) {
	dotid, err := graph.DOTID(obj)
	if err != nil {
		return nil, err
	}

	attrs := attrs.NewCopyFrom(opts.Attrs)
	attrs.Set("dotid", dotid)
	attrs.Set("name", dotid)

	// copy string metadata to node attributes
	for _, k := range obj.Metadata().Keys() {
		if v, ok := obj.Metadata().Get(k).(string); ok {
			attrs.Set(k, v)
		}
	}

	return NewNodeWithDOTID(id, obj, dotid, graph.NodeOptions{Attrs: attrs})
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
func (n *Node) Attrs() attrs.Attrs {
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
