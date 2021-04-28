package space

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is space entity.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Type returns entity type.
	Type() string
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Resource is space resource.
type Resource interface {
	Entity
	// Name returns name.
	Name() string
	// Group returns group.
	Group() string
	// Version returns version.
	Version() string
	// Kind returns kind.
	Kind() string
	// Namespaced flag.
	Namespaced() bool
}

// Object is space object.
type Object interface {
	Entity
	// Name returns name.
	Name() string
	// Namespace returns namespace.
	Namespace() string
	// Resource returns Resource.
	Resource() Resource
}

// Link between two Entities.
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
