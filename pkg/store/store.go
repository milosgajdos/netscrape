package store

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Object is store object
type Object interface {
	graph.Object
}

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
	// Graph returns graph handle.
	Graph(context.Context) (graph.Graph, error)
	// Add Entity to store.
	Add(context.Context, Entity, ...Option) error
	// Link two entities in store.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Delete Entity from store.
	Delete(context.Context, Entity, ...Option) error
	// Unlink two entities in store.
	Unlink(ctx context.Context, from, to uuid.UID, opts ...Option) error
}
