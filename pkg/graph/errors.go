package graph

import "errors"

var (
	// ErrInvalidNode is returned when attempting to use an invalid node
	ErrInvalidNode = errors.New("invalid node")
	// ErrNodeNotFound is returned when a node could not be found
	ErrNodeNotFound = errors.New("node not found")
	// ErrEdgeNotFound is returned when an edge could not be found
	ErrEdgeNotFound = errors.New("edge not found")
	// ErrEdgeNotExist is returned when an edge does not exist
	ErrEdgeNotExist = errors.New("edge does not exist")
	// ErrDuplicateNode is returned by store when duplicate nodes are found
	ErrDuplicateNode = errors.New("duplicate node")
	// ErrUnknownEntity is returned when requesting an unknown entity
	ErrUnknownEntity = errors.New("unknown entity")
	// ErrInvalidEntity is returned when requesting an invalid entity
	ErrInvalidEntity = errors.New("invalid entity")
	// ErrMissingEntity is returned when a graph query is missing entity
	ErrMissingEntity = errors.New("missing entity")
	// ErrMissingResource is returned by when scrape.Object.Resource() is nil
	ErrMissingResource = errors.New("missing resource")
	// ErrNotImplemented is returned when requesting functionality that has not been implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrUnsupported is returned when requesting unsupported functionality
	ErrUnsupported = errors.New("not supported")
)
