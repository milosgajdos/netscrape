package query

// Type is predicate type
type Type int

const (
	UID Type = iota
	Name
	Group
	Version
	Kind
	Namespace
	Weight
	Entity
	Attrs
)

func (p Type) String() string {
	switch p {
	case UID:
		return "uid"
	case Name:
		return "name"
	case Group:
		return "group"
	case Version:
		return "version"
	case Kind:
		return "kind"
	case Namespace:
		return "namesapce"
	case Weight:
		return "weight"
	case Entity:
		return "entity"
	case Attrs:
		return "attrs"
	}

	return "unknown type"
}

// EntityVal is query entity value
type EntityVal int

const (
	Node EntityVal = iota
	Edge
)

func (v EntityVal) String() string {
	switch v {
	case Node:
		return "node"
	case Edge:
		return "edge"
	}

	return "unknown entity"
}

// WildCard defines query wildcards
type WildCard int

const (
	Any WildCard = iota
)

func (v WildCard) String() string {
	switch v {
	case Any:
		return "any"
	}

	return "unknown wildcard"
}
