package memory

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"gonum.org/v1/gonum/graph/encoding"
)

// Attrs are graph attributes
type Attrs map[string]string

// New creates new attributes and returns it.
func New() (*Attrs, error) {
	attrs := make(Attrs)

	return &attrs, nil
}

// NewCopyFrom copies attributes from a and returns it.
func NewCopyFrom(a attrs.Attrs) *Attrs {
	attrs := make(Attrs)

	if a != nil {
		for _, k := range a.Keys() {
			attrs[k] = a.Get(k)
		}
	}

	return &attrs
}

// NewFromMap creates new attributes from a and returns it.
func NewFromMap(m map[string]string) *Attrs {
	a := make(Attrs)

	for k, v := range m {
		a[k] = v
	}

	return &a
}

// Keys returns all attribute keys
func (a Attrs) Keys() []string {
	keys := make([]string, len(a))

	i := 0
	for key := range a {
		keys[i] = key
		i++
	}

	return keys
}

// Get reads an attribute value for the given key and returns it.
// It returns an empty string if the attribute was not found.
func (a Attrs) Get(key string) string {
	return a[key]
}

// Set sets an attribute to the given value
func (a *Attrs) Set(key, val string) {
	(*a)[key] = val
}

// Attributes returns all attributes in a slice encoded
// as per gonum.graph.encoding requirements
func (a Attrs) Attributes() []encoding.Attribute {
	return DOTAttrs(&a)
}
