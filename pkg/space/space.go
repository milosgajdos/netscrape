package space

import (
	"context"
	"net/url"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is space Entity.
type Entity interface {
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

// Link links space objects.
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

// Object is an instance of Resource.
type Object interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Name returns Object name.
	Name() string
	// Namespace returns Object namespace.
	Namespace() string
	// Resource returns Object Resource.
	Resource() Resource
	// Link links two Objects.
	Link(uuid.UID, ...Option) error
	// Links returns all Object links.
	Links() []Link
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Origin identifies the origin of resources.
type Origin interface {
	// URL returns source URL
	URL() *url.URL
}

// Plan is space resource plan.
type Plan interface {
	// Origin returns origin.
	Origin(context.Context) (Origin, error)
	// Add adds resource to Plan.
	Add(context.Context, Resource, ...Option) error
	// Resources returns all Plan resources.
	Resources(context.Context) ([]Resource, error)
	// Get returns all Resources matching query.
	Get(context.Context, query.Query) ([]Resource, error)
}

// Top is space topology i.e. map of space Objects.
type Top interface {
	// Plan returns topology Plan.
	Plan(context.Context) (Plan, error)
	// Add adds Object to topology.
	Add(context.Context, Object, ...Option) error
	// Objects returns all topology Objects.
	Objects(context.Context) ([]Object, error)
	// Get returns all Objects matching query.
	Get(context.Context, query.Query) ([]Object, error)
}

// Planner builds space plans.
type Planner interface {
	// Plan builds plan for given Origin and returns it.
	Plan(context.Context, Origin) (Plan, error)
}

// Mapper maps space topology using Plan.
type Mapper interface {
	// Map returns Space tpology.
	Map(context.Context, Plan) (Top, error)
}

// Scraper builds Space plan and maps its Topology.
type Scraper interface {
	Planner
	Mapper
}
