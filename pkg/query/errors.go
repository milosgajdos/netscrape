package query

import "errors"

var (
	// ErrInvalidUID is returned when uid could not be decoded from the query.
	ErrInvalidUID = errors.New("ErrInvalidUID")
	// ErrInvalidEntity is returned when entity could not be decoded from the query.
	ErrInvalidEntity = errors.New("ErrInvalidEntity")
	// ErrInvalidName is returned when name could not be decoded from the query.
	ErrInvalidName = errors.New("ErrInvalidName")
	// ErrInvalidGroup is returned when group could not be decoded from the query.
	ErrInvalidGroup = errors.New("ErrInvalidGroup")
	// ErrInvalidVersion is returned when version could not be decoded from query.
	ErrInvalidVersion = errors.New("ErrInvalidVersion")
	// ErrInvalidKind is returned when kind could not be decoded from query.
	ErrInvalidKind = errors.New("ErrInvalidKind")
	// ErrInvalidNamespace is returned when namespace could not be decoded from the query.
	ErrInvalidNamespace = errors.New("ErrInvalidNamespace")
	// ErrNotImplemented is returned when requesting a feature that has not been implemented yet
	ErrNotImplemented = errors.New("ErrNotImplemented")
)
