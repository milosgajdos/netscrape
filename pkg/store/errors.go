package store

import "errors"

var (
	// ErrNotImplemented is returned when requesting unimplemented functionality.
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrUnsupported is returned when requesting unsupported functionality.
	ErrUnsupported = errors.New("ErrUnsupported")
	// ErrEntityNotFound is returned when Entity could not be found in store.
	ErrEntityNotFound = errors.New("ErrEntityNotFound")
	// ErrAlreadyExists is returned when either Entity or Link already exist in the store.
	ErrAlreadyExists = errors.New("ErrAlreadyExists")
	// ErrNotExist is returned when either Entity or Link do not exist in the store.
	ErrNotExist = errors.New("ErrNotExist")
)
