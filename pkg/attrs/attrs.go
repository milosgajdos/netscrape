package attrs

import (
	"context"

	"gonum.org/v1/gonum/graph/encoding"
)

const (
	// Name defines name attribute key.
	Name = "name"
	// DOTID defined DOT ID attribute key.
	DOTID = "dotid"
	// Weight defines weight attribute key.
	Weight = "weight"
	// Relation defines relation attribute key.
	Relation = "relation"
	// DOTLabel defines GraphViz DOT label attribute key.
	DOTLabel = "label"
)

// Attrs provide a simple key-value store
// for storing arbitrary entity attributes.
type Attrs interface {
	// Keys returns all attribute keys.
	Keys(context.Context) ([]string, error)
	// Get returns the attribute value for the given key.
	Get(context.Context, string) (string, error)
	// Set sets the value of the attribute for the given key.
	Set(ctx context.Context, key, val string) error
}

// DOT are Attrs which implement graph.DOTAttributes interface.
type DOT interface {
	// Attributes returns attributes as a slice of encoding.Attribute.
	Attributes() []encoding.Attribute
}

// ToMap returns map of attributes.
func ToMap(ctx context.Context, a Attrs) (map[string]string, error) {
	m := make(map[string]string)

	keys, err := a.Keys(ctx)
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		m[k], err = a.Get(ctx, k)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
