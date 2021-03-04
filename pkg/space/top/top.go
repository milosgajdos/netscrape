package top

import (
	"context"
	"fmt"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Top is space topology
type Top struct {
	// entities stores all entities by their UID
	entities map[string]space.Entity
	// index is topology "search index" for entity types.
	// NOTE: first index key is entity type, the second is its UID
	index map[string]map[string]space.Entity
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
		entities: make(map[string]space.Entity),
		index:    make(map[string]map[string]space.Entity),
		links:    make(map[string]space.Link),
		elinks:   make(map[string]map[string]space.Link),
		mu:       &sync.RWMutex{},
	}, nil
}

// Entities returns all entities stored in topology.
func (t Top) Entities(ctx context.Context) ([]space.Entity, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	entities := make([]space.Entity, len(t.entities))

	i := 0

	for _, entity := range t.entities {
		entities[i] = entity
		i++
	}

	return entities, nil
}

// add adds a new entity to topology
func (t *Top) add(e space.Entity) error {
	t.entities[e.UID().Value()] = e

	ent := e.Type().String()

	if t.index[ent][e.UID().Value()] == nil {
		t.index[ent] = make(map[string]space.Entity)
	}

	t.index[ent][e.UID().Value()] = e

	return nil
}

// Add adds o to topology with the given options.
func (t *Top) Add(ctx context.Context, e space.Entity, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.entities[e.UID().Value()]; !ok {
		if err := t.add(e); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes Entity with given uid from topology.
func (t *Top) Remove(ctx context.Context, uid uuid.UID, opts ...space.Option) error {
	oopts := space.Options{}
	for _, apply := range opts {
		apply(&oopts)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.entities, uid.Value())
	delete(t.elinks, uid.Value())

	for name, ents := range t.index {
		for u := range ents {
			if u == uid.Value() {
				delete(t.index[name], u)
				if len(t.index[name]) == 0 {
					delete(t.index, name)
				}
			}
		}
	}

	for luid, l := range t.links {
		from := l.From().Value()
		to := l.To().Value()

		if from == uid.Value() || to == uid.Value() {
			delete(t.links, luid)
		}
	}

	return nil
}

// link links from and to entity
func (t *Top) link(from, to uuid.UID, opts ...space.Option) error {
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

// Link links from and to entities setting the given options on the link.
// If either from or to is not found in topology it returns error.
func (t *Top) Link(ctx context.Context, from, to uuid.UID, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	l, ok := t.elinks[from.Value()][to.Value()]
	if !ok {
		return t.link(from, to, opts...)
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

// Links returns all links with origin in the entity with the given UID.
// It returns error if the entity is not found.
func (t Top) Links(ctx context.Context, uid uuid.UID) ([]space.Link, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.elinks[uid.Value()]; !ok {
		return nil, fmt.Errorf("top links entity %s: %w", uid, space.ErrEntityNotFound)
	}

	links := make([]space.Link, len(t.elinks[uid.Value()]))

	i := 0
	for _, link := range t.elinks[uid.Value()] {
		links[i] = link
		i++
	}

	return links, nil
}

// getEntities returns all entities with given entity type.
func (t Top) getEntities(ent entity.Type) ([]space.Entity, error) {
	// nolint:prealloc
	var entities []space.Entity

	if ents, ok := t.index[ent.String()]; ok {
		for _, e := range ents {
			entities = append(entities, e)
		}
	}

	return entities, nil
}

// Get queries topology entities and returns the results
func (t Top) Get(ctx context.Context, q query.Query) ([]space.Entity, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if m := q.Matcher(query.UID); m != nil {
		uid, ok := m.Predicate().Value().(uuid.UID)
		if ok && len(uid.Value()) > 0 {
			o, ok := t.entities[uid.Value()]
			if !ok {
				return []space.Entity{}, nil
			}
			return []space.Entity{o}, nil
		}
	}

	if m := q.Matcher(query.Entity); m != nil {
		ent, ok := m.Predicate().Value().(entity.Type)
		if ok {
			if _, err := entity.TypeFromString(ent.String()); err != nil {
				return nil, err
			}
			return t.getEntities(ent)
		}
		return []space.Entity{}, nil
	}

	return []space.Entity{}, nil
}
