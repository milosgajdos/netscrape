package scraper

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Scraper scrapes data.
type Scraper interface {
	// Scrape scrapes data following the provided plan.
	Scrape(context.Context, Plan, ...Option) error
}

// Resource is plan resource.
type Resource interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Type of resource.
	Type() string
	// Group retrurns group.
	Group() string
	// Version returns version.
	Version() string
	// Kind returns kind.
	Kind() string
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Plan is scrape resource plan.
// TODO: make this a DAG.
type Plan interface {
	// Add adds resource to plan.
	Add(context.Context, Resource, ...Option) error
	// Get returns entity with the given uid.
	Get(context.Context, uuid.UID, ...Option) (Resource, error)
	// GetAll returns all resources in plan.
	GetAll(context.Context, ...Option) ([]Resource, error)
	// Delete removes Resource with the given uid from plan.
	Delete(context.Context, uuid.UID, ...Option) error
}
