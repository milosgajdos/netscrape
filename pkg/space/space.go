package space

import (
	"net/url"

	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Resource is a space resource.
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
	// Metadata returns Resource metadata.
	Metadata() metadata.Metadata
}

// Link links space objects.
type Link interface {
	// UID returns unique ID.
	UID() uuid.UID
	// From returns uid of the origin of the link.
	From() uuid.UID
	// To returns uid of the end of the link.
	To() uuid.UID
	// Metadata returns Link metadata.
	Metadata() metadata.Metadata
}

// Object is an instance of a Resource.
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
	Link(uuid.UID, LinkOptions) error
	// Links returns all Object links.
	Links() []Link
	// Metadata returns Object metadata.
	Metadata() metadata.Metadata
}

// Origin identifies the origin of resources.
type Origin interface {
	// URL returns source URL
	URL() *url.URL
}

// Plan is space resource plan.
type Plan interface {
	// Origin returns origin.
	Origin() Origin
	// Add adds resource to Plan.
	Add(Resource, AddOptions) error
	// Resources returns all Plan resources.
	Resources() []Resource
	// Get returns all Resources matching query.
	Get(query.Query) ([]Resource, error)
}

// Top is space topology i.e. map of space Objects.
type Top interface {
	// Plan returns topology Plan.
	Plan() Plan
	// Add adds Object to topology.
	Add(Object, AddOptions) error
	// Objects returns all topology Objects.
	Objects() []Object
	// Get returns all Objects matching query.
	Get(query.Query) ([]Object, error)
}

// Planner builds space plans.
type Planner interface {
	// Plan builds plan for given Origin and returns it.
	Plan(Origin) (Plan, error)
}

// Mapper maps space topology using Plan.
type Mapper interface {
	// Map returns Space tpology.
	Map(Plan) (Top, error)
}

// Scraper discovers Space and maps its Topology.
type Scraper interface {
	Planner
	Mapper
}