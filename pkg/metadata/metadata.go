package metadata

// Metadata provides a simple key-valule store
// for arbitrary data of arbitrary type.
type Metadata interface {
	// Keys returns all metadata keys.
	Keys() []string
	// Get returns metadata for the given key.
	Get(string) interface{}
	// Set sets metadata for the given key.
	Set(string, interface{})
}
