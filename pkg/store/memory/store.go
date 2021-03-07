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
		return nil, store.ErrEntityNotFound
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

// BulkAdd adds entities to store.
func (m *Memory) BulkAdd(ctx context.Context, ents []store.Entity, opts ...store.Option) error {
	for _, ent := range ents {
		if err := m.Add(ctx, ent, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkGet gets entities from store.
func (m *Memory) BulkGet(ctx context.Context, uids []uuid.UID, opts ...store.Option) ([]store.Entity, error) {
	// TODO query graph for a list of UIDs
	return nil, store.ErrNotImplemented
}

// BulkDelete deletes entities from store.
func (m *Memory) BulkDelete(ctx context.Context, uids []uuid.UID, opts ...store.Option) error {
	for _, uid := range uids {
		if err := m.Delete(ctx, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkLink links from two given entities in store.
func (m *Memory) BulkLink(ctx context.Context, from uuid.UID, uids []uuid.UID, opts ...store.Option) error {
	for _, uid := range uids {
		if err := m.Link(ctx, from, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkUnlink unlinks entity from given entities in store.
func (m *Memory) BulkUnlink(ctx context.Context, from uuid.UID, uids []uuid.UID, opts ...store.Option) error {
	for _, uid := range uids {
		if err := m.Unlink(ctx, from, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}
