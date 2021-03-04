package entity

import "errors"

var (
	// ErrUnknownType is returns when a string fails to decode to Type
	ErrUnknownType = errors.New("ErrUnknownType")
)
