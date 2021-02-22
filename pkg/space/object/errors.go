package object

import "errors"

var (
	// ErrMissingResource is returns when trying to create an Object with nil Resource
	ErrMissingResource = errors.New("ErrMissingResource")
)
