package top

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Top is space topology
type Top struct {
	// index stores all entities by their UID
	index map[string]space.Entity
	// links indexes all links for the given entity
	// NOTE: index key in this map is the UID of the *link*
	links map[string]space.Link
	// elinks indexes links from each entity to other entities for faster lookups.
	// NOTE: index key in first level map is the UID of the from *entity*
	// and index key in the second level map is the UID of the to *entity*
	elinks map[string]map[string]space.Link
	// mu synchronizes access to Top
	mu *sync.RWMutex
}

// New creates a new topology and returns it.
func New() (*Top, error) {
	return &Top{
		index:  make(map[string]space.Entity),
		links:  make(map[string]space.Link),
		elinks: make(map[string]map[string]space.Link),
		mu:     &sync.RWMutex{},
	}, nil
}

func (t *Top) add(ctx context.Context, e space.Entity, opts ...space.Option) error {
	t.index[e.UID().Value()] = e

	return nil
}

// Add adds o to topology with the given options.
func (t *Top) Add(ctx context.Context, e space.Entity, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.add(ctx, e, opts...)
}

func (t *Top) getAll(ctx context.Context, opts ...space.Option) ([]space.Entity, error) {
	ents := make([]space.Entity, len(t.index))

	i := 0
	for _, ent := range t.index {
		ents[i] = ent
		i++
	}
	return ents, nil
}

// GetAll returns all entities and returns them.
func (t *Top) GetAll(ctx context.Context, opts ...space.Option) ([]space.Entity, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.getAll(ctx, opts...)
}

func (t Top) get(ctx context.Context, uid uuid.UID, opts ...space.Option) (space.Entity, error) {
	e, ok := t.index[uid.Value()]
	if !ok {
		return nil, space.ErrEntityNotFound
	}

	return e, nil
}

// Get returns entity with the given uid.
func (t Top) Get(ctx context.Context, uid uuid.UID, opts ...space.Option) (space.Entity, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.get(ctx, uid, opts...)
}

func (t *Top) delete(ctx context.Context, uid uuid.UID, opts ...space.Option) error {
	delete(t.index, uid.Value())
	delete(t.elinks, uid.Value())

	for luid, l := range t.links {
		from := l.From().Value()
		to := l.To().Value()

		if from == uid.Value() || to == uid.Value() {
			delete(t.links, luid)
		}
	}

	return nil
}

// Delete removes Entity with given uid from topology.
func (t *Top) Delete(ctx context.Context, uid uuid.UID, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.delete(ctx, uid, opts...)
}

// newLink links from and to entity
func (t *Top) newLink(from, to uuid.UID, opts ...space.Option) error {
	sopts := space.Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	lopts := []link.Option{
		link.WithUID(sopts.UID),
		link.WithAttrs(sopts.Attrs),
		link.WithMerge(sopts.Merge),
	}

	link, err := link.New(from, to, lopts...)
	if err != nil {
		return err
	}

	t.links[link.UID().Value()] = link

	if t.elinks[from.Value()] == nil {
		t.elinks[from.Value()] = make(map[string]space.Link)
	}

	t.elinks[from.Value()][to.Value()] = link

	return nil
}

func (t *Top) link(ctx context.Context, from, to uuid.UID, opts ...space.Option) error {
	if t.elinks[from.Value()] == nil {
		return t.newLink(from, to, opts...)
	}

	l, ok := t.elinks[from.Value()][to.Value()]
	if !ok {
		return t.newLink(from, to, opts...)
	}

	lopts := space.Options{}
	for _, apply := range opts {
		apply(&lopts)
	}

	if lopts.Merge {
		if lopts.Attrs != nil {
			for _, k := range lopts.Attrs.Keys() {
				l.Attrs().Set(k, lopts.Attrs.Get(k))
			}
		}
	}

	return nil
}

// Link links from and to entities setting the given options on the link.
// If either from or to is not found in topology it returns error.
func (t *Top) Link(ctx context.Context, from, to uuid.UID, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.link(ctx, from, to, opts...)
}

func (t *Top) unlink(ctx context.Context, from, to uuid.UID, opts ...space.Option) error {
	l, ok := t.elinks[from.Value()][to.Value()]
	if !ok {
		return nil
	}

	delete(t.links, l.UID().Value())
	delete(t.elinks[from.Value()], to.Value())

	return nil
}

// Unlink unlinks entities with the given uids.
func (t *Top) Unlink(ctx context.Context, from, to uuid.UID, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.unlink(ctx, from, to, opts...)
}

func (t Top) getLinks(ctx context.Context, uid uuid.UID, opts ...space.Option) ([]space.Link, error) {
	if _, ok := t.elinks[uid.Value()]; !ok {
		return []space.Link{}, nil
	}

	links := make([]space.Link, len(t.elinks[uid.Value()]))

	i := 0
	for _, link := range t.elinks[uid.Value()] {
		links[i] = link
		i++
	}

	return links, nil
}

// Links returns all links with origin in the entity with the given UID.
// It returns error if the entity is not found.
func (t Top) Links(ctx context.Context, uid uuid.UID, opts ...space.Option) ([]space.Link, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.getLinks(ctx, uid, opts...)
}

// BulkAdd adds Entites to topology.
func (t *Top) BulkAdd(ctx context.Context, ents []space.Entity, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, e := range ents {
		if err := t.add(ctx, e, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkDelete removes Entities with given uid from topology.
func (t *Top) BulkDelete(ctx context.Context, uids []uuid.UID, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, uid := range uids {
		if err := t.delete(ctx, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkGet returns entities with the given UIDs.
func (t *Top) BulkGet(ctx context.Context, uids []uuid.UID, opts ...space.Option) ([]space.Entity, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	ents := make([]space.Entity, len(uids))
	for i, uid := range uids {
		e, err := t.get(ctx, uid)
		if err != nil {
			return nil, err
		}
		ents[i] = e
	}
	return ents, nil
}

// BulkLink links from entity to entities with given UIDs.
func (t *Top) BulkLink(ctx context.Context, from uuid.UID, to []uuid.UID, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, uid := range to {
		if err := t.link(ctx, from, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkUnlink unlinks from entity from entities with given UIDs
func (t *Top) BulkUnlink(ctx context.Context, from uuid.UID, to []uuid.UID, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, uid := range to {
		if err := t.unlink(ctx, from, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkLinks returns all links with origin in the given entity.
func (t *Top) BulkLinks(ctx context.Context, uids []uuid.UID) (map[string][]space.Link, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	links := make(map[string][]space.Link)
	for _, uid := range uids {
		lx, err := t.getLinks(ctx, uid)
		if err != nil {
			return nil, err
		}
		links[uid.Value()] = lx
	}
	return links, nil
}
