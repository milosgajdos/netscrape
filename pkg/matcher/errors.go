package matcher

import "errors"

var (
	// ErrFilterNotFound is returned when matching against a property filter that has not been defined
	ErrFilterNotFound = errors.New("ErrFilterNotFound")
)
