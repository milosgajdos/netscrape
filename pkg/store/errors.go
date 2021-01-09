package store

import "errors"

var (
	// ErrUnknownEntity is returned when requesting an unknown entity
	ErrUnknownEntity = errors.New("unknown entity")
	// ErrNotImplemented is returned when requesting a feature that has not been implemented yet
	ErrNotImplemented = errors.New("not implemented")
	// ErrUnsupported is returned when requesting unsupported functionality
	ErrUnsupported = errors.New("not supported")
	// ErrNodeNotFound is returned when Node entity could not be found in store
	ErrNodeNotFound = errors.New("node not found")
)
