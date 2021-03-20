package space

import "errors"

var (
	// ErrEntityNotFound is returned an Entity could not be found in Top.
	ErrEntityNotFound = errors.New("ErrEntityNotFound")
	// ErrNotImplemented is returned when requesting a feature that has not been implemented yet.
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrNoLinksFound is returned when querying for links and none are found
	ErrNoLinksFound = errors.New("ErrNoLinksFound")
)
