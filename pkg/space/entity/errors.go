package entity

import "errors"

var (
	// ErrMissingResource is returned when Entity has no Resource
	ErrMissingResource = errors.New("ErrMissingResource")
)
