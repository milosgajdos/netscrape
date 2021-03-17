package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Memory is in-memory store.
type Memory struct {
	// uid is the UID of the graph
	uid uuid.UID
	// g is the store graph
	g memory.Graph
	// mu synchronizes access to store
	mu *sync.RWMutex
}

// NewStore creates a new in-memory store backed by graph g and returns it.
// By default store uses memory.WUG unless overridden by WithGraph options.
func NewStore(opts ...Option) (*Memory, error) {
	gopts := Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	uid := gopts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
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
		uid: uid,
		g:   g,
		mu:  &sync.RWMutex{},
	}, nil
}

func (m Memory) UID() uuid.UID {
	return m.uid
}

// Graph returns graph handle.
func (m *Memory) Graph() (graph.Graph, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	g := m.g

	return g, nil
}

func (m *Memory) add(ctx context.Context, e store.Entity, opts ...store.Option) error {
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

// Add stores e in memory store.
func (m *Memory) Add(ctx context.Context, e store.Entity, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.add(ctx, e, opts...)
}

func (m *Memory) get(ctx context.Context, uid uuid.UID, opts ...store.Option) (store.Entity, error) {
	gopts := store.Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	e, err := m.g.(memory.Graph).Node(ctx, uid)
	if err != nil {
		if errors.Is(err, graph.ErrNodeNotFound) {
			return nil, store.ErrEntityNotFound
		}
		return nil, err
	}

	return e.(*memory.Node), nil
}

// Get Entity from store.
func (m *Memory) Get(ctx context.Context, uid uuid.UID, opts ...store.Option) (store.Entity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.get(ctx, uid, opts...)
}

func (m *Memory) delete(ctx context.Context, uid uuid.UID, opts ...store.Option) error {
	dopts := store.Options{}
	for _, apply := range opts {
		apply(&dopts)
	}

	return m.g.RemoveNode(ctx, uid)
}

// Delete deletes e from memory store.
func (m *Memory) Delete(ctx context.Context, uid uuid.UID, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.delete(ctx, uid, opts...)
}

func (m *Memory) link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	lopts := store.Options{}
	for _, apply := range opts {
		apply(&lopts)
	}

	if _, err := m.g.Link(ctx, from, to, graph.WithAttrs(lopts.Attrs)); err != nil {
		return err
	}

	return nil
}

// Link links entities with given UIDs in store.
func (m *Memory) Link(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.link(ctx, from, to, opts...)
}

func (m *Memory) unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	ulopts := store.Options{}
	for _, apply := range opts {
		apply(&ulopts)
	}

	if err := m.g.Unlink(ctx, from, to); err != nil {
		return err
	}

	return nil
}

// Unlink two entities with given UIDs in store.
func (m *Memory) Unlink(ctx context.Context, from, to uuid.UID, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.unlink(ctx, from, to, opts...)
}

// BulkAdd adds entities to store.
func (m *Memory) BulkAdd(ctx context.Context, ents []store.Entity, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ent := range ents {
		if err := m.add(ctx, ent, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkGet gets entities from store.
func (m *Memory) BulkGet(ctx context.Context, uids []uuid.UID, opts ...store.Option) ([]store.Entity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ents := make([]store.Entity, len(uids))

	i := 0
	for _, uid := range uids {
		e, err := m.get(ctx, uid, opts...)
		if err != nil {
			return nil, err
		}
		ents[i] = e
	}
	return ents, nil
}

// BulkDelete deletes entities from store.
func (m *Memory) BulkDelete(ctx context.Context, uids []uuid.UID, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, uid := range uids {
		if err := m.delete(ctx, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkLink links from two given entities in store.
func (m *Memory) BulkLink(ctx context.Context, from uuid.UID, uids []uuid.UID, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, uid := range uids {
		if err := m.link(ctx, from, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkUnlink unlinks entity from given entities in store.
func (m *Memory) BulkUnlink(ctx context.Context, from uuid.UID, uids []uuid.UID, opts ...store.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, uid := range uids {
		if err := m.unlink(ctx, from, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}
