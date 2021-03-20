package plan

import (
	"context"
	"net/url"

	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Origin identifies the origin of space.
type Origin interface {
	// URL returns source URL
	URL() *url.URL
}

// Plan is scrape resource plan.
type Plan interface {
	// Origin returns origin.
	Origin(context.Context) (Origin, error)
	// Add adds resource to plan.
	Add(context.Context, space.Resource, ...Option) error
	// GetAll returns all resources in plan.
	GetAll(context.Context, ...Option) ([]space.Resource, error)
	// Get returns entity with the given uid.
	Get(context.Context, uuid.UID, ...Option) (space.Resource, error)
	// Delete removes Resource with the given uid from plan.
	Delete(context.Context, uuid.UID, ...Option) error
}
