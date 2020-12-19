package uuid

// UID is a generic UID.
type UID interface {
	// Value returns string UID value
	Value() string
}
