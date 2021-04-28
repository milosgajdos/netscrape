package plan

import "errors"

var (
	// ErrResourceNotFound is returned when a Resource could not be found in Plan.
	ErrResourceNotFound = errors.New("ErrResourceNotFound")
	// ErrNotImplemented is returned when requesting a feature that has not been implemented yet.
	ErrNotImplemented = errors.New("ErrNotImplemented")
)
