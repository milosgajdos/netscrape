package memory

import (
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Store is in-memory store.
type Store interface {
	store.Store
}

// BulkStore is in-memory bulk store
type BulkStore interface {
	store.BulkStore
}
