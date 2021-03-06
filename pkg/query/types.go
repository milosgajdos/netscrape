package query

// Type is predicate type
type Type int

const (
	UID Type = iota
	Entity
	Name
	Group
	Version
	Kind
	Namespace
	Weight
	Attrs
)

// String implements fmt.Stringer
func (p Type) String() string {
	switch p {
	case UID:
		return "uid"
	case Entity:
		return "entity"
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
	case Attrs:
		return "attrs"
	}

	return "unknown"
}

// WildCard defines query wildcards.
type WildCard int

const (
	// Any means any value is acceptable
	Any WildCard = iota
)

// String implements fmt.Stringer.
func (v WildCard) String() string {
	switch v {
	case Any:
		return "any"
	}

	return "unknown"
}
