package property

// Property is property type
type Property int

const (
	Type Property = iota
	Name
	Group
	Version
	Kind
	Namespace
	Weight
	Attrs
)

// String implements fmt.Stringer.
func (p Property) String() string {
	switch p {
	case Type:
		return "type"
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
