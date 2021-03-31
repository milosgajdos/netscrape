package marshal

import "errors"

var (
	// ErrUnsuportedType is returned when attempting to marshal unsupported type.
	ErrUnsuportedType = errors.New("ErrUnsuportedType")
	// ErrUnsupportedFormat is returned when attempting to marshal into unsupported data type.
	ErrUnsupportedFormat = errors.New("ErrUnsupportedFormat")
	// ErrNotImplemented is returned when requesting a feature that has not been implemented yet.
	ErrNotImplemented = errors.New("ErrNotImplemented")
)
