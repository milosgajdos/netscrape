package cache

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
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

// Links is a simple in-memory key-value store for links.
type Links interface {
	// Put stores link in the cache.
	Put(context.Context, Link, ...Option) error
	// GetFrom returns all links from the given UID.
	GetFrom(context.Context, uuid.UID, ...Option) ([]Link, error)
	// GetTo returns all the links to the given UID>
	GetTo(context.Context, uuid.UID, ...Option) ([]Link, error)
	// Delete removes all links to and from the given UID.
	Delete(context.Context, uuid.UID, ...Option) error
	// Clear clears the whole cache.
	Clear(context.Context, ...Option) error
}

// BulkLinks provides bulk operations on links.
type BulkLinks interface {
	Links
	// BulkPut puts all links key-ed by UID into cache.
	BulkPut(context.Context, []Link, ...Option) error
	// BulkGetFrom returns all links from the given UIDs.
	BulkGetFrom(context.Context, []uuid.UID, ...Option) (map[uuid.UID][]Link, error)
	// BulkGetTo returns all links to the given UIDs.
	BulkGetTo(context.Context, []uuid.UID, ...Option) (map[uuid.UID][]Link, error)
	// BulkDelete removes all links to and from the given UIDs.
	BulkDelete(context.Context, []uuid.UID, ...Option) error
}
