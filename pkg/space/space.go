package space

import (
	"context"
	"net/url"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Resource is space resource.
type Resource interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Type of entity
	Type() entity.Type
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
	// Type of entity
	Type() entity.Type
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
	// GetAll returns all resources in plan.
	GetAll(context.Context) ([]Resource, error)
	// Get returns entity with the given uid.
	Get(context.Context, uuid.UID, ...Option) (Resource, error)
	// Delete removes Resource with the given uid from plan.
	Delete(context.Context, uuid.UID, ...Option) error
}

// Top is space topology i.e. a map of Entities.
type Top interface {
	// Add adds Entity to topology.
	Add(context.Context, Entity, ...Option) error
	// GetAll returns all entities store in topology.
	GetAll(context.Context) ([]Entity, error)
	// Get returns entity with the given uid.
	Get(context.Context, uuid.UID, ...Option) (Entity, error)
	// Delete removes Entity with given uid from topology.
	Delete(context.Context, uuid.UID, ...Option) error
	// Link links entities with given UIDs.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Unlink unlinks entities with given UIDs.
	Unlink(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Links returns all links with origin in the given entity.
	Links(context.Context, uuid.UID, ...Option) ([]Link, error)
}

// BulkTop provides bulk operations on topology
type BulkTop interface {
	Top
	// BulkAdd adds Entites to topology.
	BulkAdd(context.Context, []Entity, ...Option) error
	// BulkDelete removes Entities with given uid from topology.
	BulkDelete(context.Context, []uuid.UID, ...Option) error
	// BulkGet returns entities with the given UIDs.
	BulkGet(context.Context, []uuid.UID, ...Option) ([]Entity, error)
	// BulkLink links from entity to entities with given UIDs.
	BulkLink(ctx context.Context, from uuid.UID, to []uuid.UID, opts ...Option) error
	// BulkUnlink unlinks from entity from entities with given UIDs.
	BulkUnlink(ctx context.Context, from uuid.UID, to []uuid.UID, opts ...Option) error
	// BulkLinks returns all links with origin in the given entity.
	BulkLinks(context.Context, []uuid.UID, ...Option) (map[string][]Link, error)
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
