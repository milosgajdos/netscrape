package query

// PredKind is Predicate kind
type PredKind int

const (
	PUID PredKind = iota
	PName
	PGroup
	PVersion
	PKind
	PNamespace
	PWeight
	PEntity
	PAttrs
	PMetadata
)

func (p PredKind) String() string {
	switch p {
	case PUID:
		return "uid"
	case PName:
		return "name"
	case PGroup:
		return "group"
	case PVersion:
		return "version"
	case PKind:
		return "kind"
	case PNamespace:
		return "namesapce"
	case PWeight:
		return "weight"
	case PEntity:
		return "entity"
	case PAttrs:
		return "attrs"
	case PMetadata:
		return "metadata"
	}

	return "unknown predicate"
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

	return "unknown entity value"
}

// Match defines query match wildcards
type Match int

const (
	Any Match = iota
)

func (v Match) String() string {
	switch v {
	case Any:
		return "any"
	}

	return "unknown value"
}
