package broker

// Type is entity type
type Type int

const (
	Entity Type = iota
	Resource
	Link
	Unknown
)

const (
	entityString   = "Entity"
	resourceString = "Resource"
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
	default:
		return unknownString
	}
}

// Message is broker message.
type Message struct {
	// UID is unique message ID.
	UID string
	// Type is a message type.
	Type Type
	// Data contains message payload.
	Data []byte
	// Attrs are message attributes.
	Attrs map[string]string
}
