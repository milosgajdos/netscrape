package attrs

import (
	"gonum.org/v1/gonum/graph/encoding"
)

const (
	// Name is name attribute key
	Name = "name"
	// DOTID is DOT ID attribute key
	DOTID = "dotid"
	// Weight is weight attribute key
	Weight = "weight"
	// Rel is relation attribute
	Relation = "relation"
	// DOTLabel is GraphViz DOT label attribute
	DOTLabel = "label"
)

// Attrs provide a simple key-value store
// for storing arbitrary entity attributes.
type Attrs interface {
	// Keys returns all attribute keys.
	Keys() []string
	// Get returns the attribute value for the given key.
	Get(string) string
	// Set sets the value of the attribute for the given key.
	Set(string, string)
}

// DOT are Attrs which implement graph.DOTAttributes interface.
type DOT interface {
	// Attributes returns attributes as a slice of encoding.Attribute.
	Attributes() []encoding.Attribute
}
