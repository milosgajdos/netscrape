package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"gonum.org/v1/gonum/graph/encoding"
)

// Attrs are graph attributes
type Attrs map[string]string

// New creates new attributes and returns it.
func New() *Attrs {
	attrs := make(Attrs)

	return &attrs
}

// NewCopyFrom copies attributes from a and returns it.
func NewCopyFrom(ctx context.Context, a attrs.Attrs) (*Attrs, error) {
	attrs := make(Attrs)

	keys, err := a.Keys(ctx)
	if err != nil {
		return nil, err
	}

	if a != nil {
		for _, k := range keys {
			attrs[k], err = a.Get(ctx, k)
			if err != nil {
				return nil, err
			}
		}
	}

	return &attrs, nil
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
func (a Attrs) Keys(ctx context.Context) ([]string, error) {
	keys := make([]string, len(a))

	i := 0
	for key := range a {
		keys[i] = key
		i++
	}

	return keys, nil
}

// Get reads an attribute value for the given key and returns it.
// It returns an empty string if the attribute was not found.
func (a Attrs) Get(ctx context.Context, key string) (string, error) {
	return a[key], nil
}

// Set sets an attribute to the given value
func (a *Attrs) Set(ctx context.Context, key, val string) error {
	(*a)[key] = val
	return nil
}

// Attributes returns all attributes in a slice encoded
// as per gonum.graph.encoding requirements
func (a Attrs) Attributes() []encoding.Attribute {
	return DOTAttrs(&a)
}
