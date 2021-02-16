package space

import (
	"context"
	"net/url"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is graph entity.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Name returns name
	Name() string
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Resource is space resource.
type Resource interface {
	// UID returns unique ID.
	UID() uuid.UID
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

// Object is an instance of resource.
type Object interface {
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

// Origin identifies the origin of space.
type Origin interface {
	// URL returns source URL
	URL() *url.URL
}

// Plan is resource plan.
type Plan interface {
	// Origin returns origin.
	Origin(context.Context) (Origin, error)
	// Add adds resource to plan.
	Add(context.Context, Resource, ...Option) error
	// Remove removes Resource from plan.
	Remove(context.Context, Resource, ...Option) error
	// Resources returns all plan resources.
	Resources(context.Context) ([]Resource, error)
	// Get returns all Resources matching query.
	Get(context.Context, query.Query) ([]Resource, error)
}

// Top is topology i.e. a map of Entities.
type Top interface {
	// Add adds Entity to topology.
	Add(context.Context, Entity, ...Option) error
	// Remove removes Entity with given uid from topology.
	Remove(context.Context, uuid.UID, ...Option) error
	// Entities returns all topology Entities.
	Entities(context.Context) ([]Entity, error)
	// Get returns Entity with the given UID
	Get(context.Context, query.Query) ([]Entity, error)
	// Link links entities with given UIDs.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Links returns all links with origin in the given entity.
	Links(context.Context, uuid.UID) ([]Link, error)
}

// Planner builds plan.
type Planner interface {
	// Plan builds plan for given origin and returns it.
	Plan(context.Context, Origin) (Plan, error)
}

// Mapper maps topology.
type Mapper interface {
	// Map returns Top built from Plan.
	Map(context.Context, Plan) (Top, error)
}

// Scraper builds plan and maps topology.
type Scraper interface {
	Planner
	Mapper
}
