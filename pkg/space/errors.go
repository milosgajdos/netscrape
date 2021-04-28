package space

import "errors"

var (
	// ErrNotImplemented is returned when requesting a feature that has not been implemented yet.
	ErrNotImplemented = errors.New("ErrNotImplemented")
)
