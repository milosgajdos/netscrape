package graph

import "errors"

var (
	// ErrInvalidNode is returned when attempting to use an invalid node
	ErrInvalidNode = errors.New("ErrInvalidNode")
	// ErrNodeNotFound is returned when a node could not be found
	ErrNodeNotFound = errors.New("ErrNodeNotFound")
	// ErrEdgeNotFound is returned when an edge could not be found
	ErrEdgeNotFound = errors.New("ErrEdgeNotFound")
	// ErrEdgeNotExist is returned when an edge does not exist
	ErrEdgeNotExist = errors.New("ErrEdgeNotExist")
	// ErrDuplicateNode is returned by store when duplicate nodes are found
	ErrDuplicateNode = errors.New("ErrDuplicateNode")
	// ErrUnknownEntity is returned when requesting an unknown entity
	ErrUnknownEntity = errors.New("ErrUnknownEntity")
	// ErrInvalidEntity is returned when requesting an invalid entity
	ErrInvalidEntity = errors.New("ErrInvalidEntity")
	// ErrMissingEntity is returned when a graph query is missing entity
	ErrMissingEntity = errors.New("ErrMissingEntity")
	// ErrMissingResource is returned by when scrape.Object.Resource() is nil
	ErrMissingResource = errors.New("ErrMissingResource")
	// ErrNotImplemented is returned when requesting functionality that has not been implemented
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrUnsupported is returned when requesting unsupported functionality
	ErrUnsupported = errors.New("ErrUnsupported")
)
