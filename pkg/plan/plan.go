package plan

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Resource is Plan resource.
type Resource interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Group retrurns group.
	Group() string
	// Version returns version.
	Version() string
	// Kind returns kind.
	Kind() string
}

// Plan is scrape resource plan.
// TODO: change this interface so we can have a DAG implementation.
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
