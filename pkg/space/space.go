package space

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Resource is space resource.
type Resource interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Type of resource.
	Type() string
	// Name returns name.
	Name() string
	// Group retrurns group.
	Group() string
	// Version returns version.
	Version() string
	// Kind returns kind.
	Kind() string
	// Namespaced flag.
	Namespaced() bool
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Entity is space entity.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Type of entity.
	Type() string
	// Name returns human readable name.
	Name() string
	// Namespace returns namespace.
	Namespace() string
	// Resource returns Resource.
	Resource() Resource
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Link links space entities.
type Link interface {
	// UID returns unique ID.
	UID() uuid.UID
	// From returns uid of the origin of link.
	From() uuid.UID
	// To returns uid of the end of link.
	To() uuid.UID
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}
