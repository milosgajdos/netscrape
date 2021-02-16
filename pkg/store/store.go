package store

import (
	"context"
	"text/template"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is store entity.
type Entity interface {
	graph.Entity
}

// Querier queries store.
type Querier interface {
	// Query store and return the results.
	Query(context.Context, template.Template, map[string]string) ([]Entity, error)
}

// Store stores entities.
type Store interface {
	// Graph returns graph handle.
	Graph(context.Context) (graph.Graph, error)
	// Add Entity to store.
	Add(context.Context, Entity, ...Option) error
	// Delete Entity from store.
	Delete(context.Context, uuid.UID, ...Option) error
	// Link two entities in store.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Unlink two entities in store.
	Unlink(ctx context.Context, from, to uuid.UID, opts ...Option) error
}
