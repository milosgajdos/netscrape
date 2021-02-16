package memory

import (
	"context"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Store is in-memory store.
type Store struct {
	// g is the store graph
	g memory.Graph
}

// New creates a new in-memory store backed by graph g and returns it.
// By default store uses memory.WUG unless overridden by WithGraph options.
func New(opts ...Option) (*Store, error) {
	gopts := Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	g := gopts.Graph
	if g == nil {
		var err error
		g, err = memory.NewWUG()
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		g: g,
	}, nil
}

// Graph returns graph handle.
func (m *Store) Graph(ctx context.Context) (graph.Graph, error) {
	return m.g, nil
}

// Add stores e in memory store.
func (m *Store) Add(ctx context.Context, e store.Entity, opts ...store.Option) error {
	aopts := store.Options{}
	for _, apply := range opts {
		apply(&aopts)
	}

	n, err := m.g.NewNode(ctx, e, graph.WithAttrs(aopts.Attrs))
	if err != nil {
		return err
	}

	return m.g.AddNode(ctx, n)
}

// Link links entities with given UIDs in store.
func (m *Store) Link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	lopts := store.Options{}
	for _, apply := range opts {
		apply(&lopts)
	}

	if _, err := m.g.Link(ctx, from, to, graph.WithAttrs(lopts.Attrs)); err != nil {
		return err
	}

	return nil
}

// Delete deletes e from memory store.
func (m *Store) Delete(ctx context.Context, e store.Entity, opts ...store.Option) error {
	dopts := store.Options{}
	for _, apply := range opts {
		apply(&dopts)
	}

	return m.g.RemoveNode(ctx, e.UID())
}

// Unlink two entities with given UIDs in store.
func (m *Store) Unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	ulopts := store.Options{}
	for _, apply := range opts {
		apply(&ulopts)
	}

	if err := m.g.Unlink(ctx, from, to); err != nil {
		return err
	}

	return nil
}

// Query queries the store and returns the results.
func (m Store) Query(ctx context.Context, q query.Query) ([]store.Entity, error) {
	g, ok := m.g.(memory.Querier)
	if !ok {
		return nil, fmt.Errorf("query: %w", graph.ErrUnsupported)
	}

	qents, err := g.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	results := make([]store.Entity, len(qents))

	for i, e := range qents {
		results[i] = e.(store.Entity)
	}

	return results, nil
}
