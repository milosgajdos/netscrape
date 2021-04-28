package broker

// Type defines a message payload type.
// It is used to identify what type is
// encoded in the broker message payload.
type Type int

const (
	// Entity is space.Entity.
	Entity Type = iota
	// Object is space.Object
	Object
	// Resource is space.Resource.
	Resource
	// Link is space.Link.
	Link
	// Unknown type.
	Unknown
)

const (
	entityString   = "Entity"
	resourceString = "Resource"
	objectString   = "Object"
	linkString     = "Link"
	unknownString  = "Unknown"
)

// String implements fmt.Stringer
func (t Type) String() string {
	switch t {
	case Entity:
		return entityString
	case Resource:
		return resourceString
	case Link:
		return linkString
	case Object:
		return objectString
	default:
		return unknownString
	}
}
