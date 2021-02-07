package space

import (
	"context"
	"net/url"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Object is space Object.
type Object interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Resource is space resource.
type Resource interface {
	// Name returns Resource name.
	Name() string
	// Group retrurns Resource group.
	Group() string
	// Version returns Resource version.
	Version() string
	// Kind returns Resource kind.
	Kind() string
	// Namespaced returns true if Resource is namespaced.
	Namespaced() bool
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Link links space entities.
type Link interface {
	// UID returns unique ID.
	UID() uuid.UID
	// From returns uid of the origin of the link.
	From() uuid.UID
	// To returns uid of the end of the link.
	To() uuid.UID
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Entity is an instance of Resource.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Name returns human readable name.
	Name() string
	// Namespace returns namespace.
	Namespace() string
	// Resource returns Resource.
	Resource() Resource
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Origin identifies the origin of resources.
type Origin interface {
	// URL returns source URL
	URL() *url.URL
}

// Plan is resource plan.
type Plan interface {
	// Origin returns origin.
	Origin(context.Context) (Origin, error)
	// Add adds resource to Plan.
	Add(context.Context, Resource, ...Option) error
	// Remove removes Resource from Plan.
	Remove(context.Context, Resource, ...Option) error
	// Resources returns all Plan resources.
	Resources(context.Context) ([]Resource, error)
	// Get returns all Resources matching query.
	Get(context.Context, query.Query) ([]Resource, error)
}

// Top is topology i.e. a map of Entities.
type Top interface {
	// Plan returns topology Plan.
	Plan(context.Context) (Plan, error)
	// Entities returns all topology Entities.
	Entities(context.Context) ([]Entity, error)
	// Get returns all Entities matching query.
	Get(context.Context, query.Query) ([]Entity, error)
	// Add adds Entity to topology.
	Add(context.Context, Entity, ...Option) error
	// Remove removes Entity with given uid from topology.
	Remove(context.Context, uuid.UID, ...Option) error
	// Link links entities with given UIDs.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Links returns all links with origin in the given entity.
	Links(context.Context, uuid.UID) ([]Link, error)
}

// Planner builds Plan.
type Planner interface {
	// Plan builds plan for given Origin and returns it.
	Plan(context.Context, Origin) (Plan, error)
}

// Mapper maps topology.
type Mapper interface {
	// Map returns Top built from Plan.
	Map(context.Context, Plan) (Top, error)
}

// Scraper builds Plan and maps topology.
type Scraper interface {
	Planner
	Mapper
}
