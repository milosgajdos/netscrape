package store

import (
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
)

// Entity is store entity.
type Entity interface {
	graph.Entity
}

// Querier queries store.
// NOTE: this is a temporary hack
// Store must provide query capability by default.
// Ideally, I would like to figure out how to parse
// a generic GraphQL query into query.Query interface.
// query package must be properly reworked.
type Querier interface {
	// Query store and return the results.
	Query(query.Query) ([]Entity, error)
}

// Store stores entities.
type Store interface {
	// ID returns store id.
	ID() string
	// Graph returns store graph handle.
	Graph() graph.Graph
	// Add entity to the store.
	Add(Entity, AddOptions) error
	// Delete entity from the store.
	Delete(Entity, DelOptions) error
}
