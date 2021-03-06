package store

import "errors"

var (
	// ErrNotImplemented is returned when requesting unimplemented functionality.
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrUnsupported is returned when requesting unsupported functionality.
	ErrUnsupported = errors.New("ErrUnsupported")
	// ErrEntityNotFound is returned when Entity could not be found in store.
	ErrEntityNotFound = errors.New("ErrEntityNotFound")
)
