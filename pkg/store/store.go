package store

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Link is a link between two entities.
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

// Entity is stored in Store.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Type of entity.
	Type() string
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Store stores entities.
type Store interface {
	// Add Entity to store.
	Add(context.Context, Entity, ...Option) error
	// Get Entity from store.
	Get(context.Context, uuid.UID, ...Option) (Entity, error)
	// Delete Entity from store.
	Delete(context.Context, uuid.UID, ...Option) error
	// Link two entities in store.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Unlink two entities in store.
	Unlink(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Graph returns graph handle.
	Graph(ctx context.Context, opts ...Option) (graph.Graph, error)
}

// BulkStore stores bulks of entities.
type BulkStore interface {
	Store
	// BulkAdd adds entities to store.
	BulkAdd(context.Context, []Entity, ...Option) error
	// BulkGet gets entities from store.
	BulkGet(context.Context, []uuid.UID, ...Option) ([]Entity, error)
	// BulkDelete deletes entities from store.
	BulkDelete(context.Context, []uuid.UID, ...Option) error
	// BulkLink links the given entity to the list of given entities in store.
	BulkLink(context.Context, uuid.UID, []uuid.UID, ...Option) error
	// BulkUnlink unlinks the given entity from the list of given entities in store.
	BulkUnlink(context.Context, uuid.UID, []uuid.UID, ...Option) error
}
