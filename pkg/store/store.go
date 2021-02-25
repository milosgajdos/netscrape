package store

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is store entity.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Name returns name
	Name() string
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}

// Store stores entities.
type Store interface {
	// Add Entity to store.
	Add(context.Context, Entity, ...Option) error
	// Delete Entity from store.
	Delete(context.Context, uuid.UID, ...Option) error
	// Link two entities in store.
	Link(ctx context.Context, from, to uuid.UID, opts ...Option) error
	// Unlink two entities in store.
	Unlink(ctx context.Context, from, to uuid.UID, opts ...Option) error
}
