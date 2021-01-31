package store

import "errors"

var (
	// ErrUnknownEntity is returned when trying to use unknown entity
	ErrUnknownEntity = errors.New("ErrUnknownEntity")
	// ErrNotImplemented is returned when requesting unimplemented functionality
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrUnsupported is returned when requesting unsupported functionality
	ErrUnsupported = errors.New("ErrUnsupported")
	// ErrNodeNotFound is returned when Node could not be found in store
	ErrNodeNotFound = errors.New("ErrNodeNotFound")
)
