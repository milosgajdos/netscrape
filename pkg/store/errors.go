package store

import "errors"

var (
	// ErrUnknownObject is returned when trying to manage unknown store.Object
	ErrUnknownObject = errors.New("ErrUnknownObject")
	// ErrNotImplemented is returned when requesting unimplemented functionality
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrUnsupported is returned when requesting unsupported functionality
	ErrUnsupported = errors.New("ErrUnsupported")
	// ErrNodeNotFound is returned when Node could not be found in store
	ErrNodeNotFound = errors.New("ErrNodeNotFound")
)
