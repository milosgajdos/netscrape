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
// NOTE: this interface is a [temporary] hack!
// Store must provide query capabilities by default.
// Ideally, I would like to figure out how to parse
// a generic GraphQL query into query.Query interface.
// Or maybe store.Query should simply accept string,
// which would then have to be parsed into query.Query.
type Querier interface {
	// Query store and return the results.
	Query(context.Context, query.Query) ([]Entity, error)
}

// Store stores entities.
type Store interface {
	// Graph returns store graph handle.
	Graph(context.Context) (graph.Graph, error)
	// Add entity to the store.
	Add(context.Context, Entity, ...Option) error
	// Delete entity from the store.
	Delete(context.Context, Entity, ...Option) error
}
