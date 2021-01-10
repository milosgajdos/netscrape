package store

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
)

// Entity is store entity.
type Entity interface {
	graph.Entity
}

// Querier queries store.
// NOTE: this interface is a temporary hack/placeholder!
// Store must provide some query capability by default.
// Ideally, I would like to figure out how to parse
// a generic GraphQL query into query.Query interface.
type Querier interface {
	// Query store and return the results.
	Query(context.Context, query.Query) ([]Entity, error)
}

// Store stores entities.
type Store interface {
	// ID returns store id.
	ID() string
	// Graph returns store graph handle.
	Graph(context.Context) (graph.Graph, error)
	// Add entity to the store.
	Add(context.Context, Entity, AddOptions) error
	// Delete entity from the store.
	Delete(context.Context, Entity, DelOptions) error
}
