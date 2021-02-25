package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Memory is in-memory store.
type Memory struct {
	// g is the store graph
	g memory.Graph
}

// NewStore creates a new in-memory store backed by graph g and returns it.
// By default store uses memory.WUG unless overridden by WithGraph options.
func NewStore(opts ...Option) (*Memory, error) {
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

	return &Memory{
		g: g,
	}, nil
}

// Graph returns graph handle.
func (m *Memory) Graph() (graph.Graph, error) {
	return m.g, nil
}

// Add stores e in memory store.
func (m *Memory) Add(ctx context.Context, e store.Entity, opts ...store.Option) error {
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

// Get Entity from store.
func (m *Memory) Get(ctx context.Context, uid uuid.UID, opts ...store.Option) (store.Entity, error) {
	gopts := store.Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	g, ok := m.g.(memory.Querier)
	if !ok {
		return nil, store.ErrNotImplemented
	}

	q := base.Build().Add(predicate.UID(uid))

	entities, err := g.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}

	switch {
	case len(entities) > 1:
		panic("duplicate nodes found")
	case len(entities) == 0:
		return nil, store.ErrNodeNotFound
	}

	return entities[0], nil
}

// Link links entities with given UIDs in store.
func (m *Memory) Link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
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
func (m *Memory) Delete(ctx context.Context, uid uuid.UID, opts ...store.Option) error {
	dopts := store.Options{}
	for _, apply := range opts {
		apply(&dopts)
	}

	return m.g.RemoveNode(ctx, uid)
}

// Unlink two entities with given UIDs in store.
func (m *Memory) Unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	ulopts := store.Options{}
	for _, apply := range opts {
		apply(&ulopts)
	}

	if err := m.g.Unlink(ctx, from, to); err != nil {
		return err
	}

	return nil
}
