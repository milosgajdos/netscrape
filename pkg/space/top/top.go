package top

import (
	"context"
	"fmt"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Top is space topology
type Top struct {
	// space is the source of topology
	space space.Plan
	// entities stores all entities by their UID
	entities map[string]space.Entity
	// index is topology "search index" (ns/kind/name)
	index map[string]map[string]map[string]space.Entity
	// links indexes all links for the given entity
	// NOTE: index key in this map is the UID of the *link*
	links map[string]space.Link
	// elinks indexes links from this entity to
	// other entities for faster lookups.
	// NOTE: index key in first level map is the UID of the from *entity*
	// and index key in the second level map is the UID of the to *entity*
	elinks map[string]map[string]space.Link
	// mu synchronizes access to Top
	mu *sync.RWMutex
}

// New creates a new topology and returns it.
func New(a space.Plan) (*Top, error) {
	return &Top{
		space:    a,
		entities: make(map[string]space.Entity),
		index:    make(map[string]map[string]map[string]space.Entity),
		links:    make(map[string]space.Link),
		elinks:   make(map[string]map[string]space.Link),
		mu:       &sync.RWMutex{},
	}, nil
}

// Plan returns topology Plan.
func (t Top) Plan(ctx context.Context) (space.Plan, error) {
	return t.space, nil
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
func (t *Top) add(o space.Entity) error {
	t.entities[o.UID().Value()] = o

	ns := o.Namespace()

	if t.index[ns] == nil {
		t.index[ns] = make(map[string]map[string]space.Entity)
	}

	kind := o.Resource().Kind()

	if t.index[ns][kind] == nil {
		t.index[ns][kind] = make(map[string]space.Entity)
	}

	name := o.Name()

	t.index[ns][kind][name] = o

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

// Link links entities with given UIDs.
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

// Links returns all links with origin in the given entity.
// If the entity is not found in Top it returns error.
func (t Top) Links(ctx context.Context, uid uuid.UID) ([]space.Link, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.elinks[uid.Value()]; !ok {
		return nil, fmt.Errorf("link entity %s: %w", uid, space.ErrEntityNotFound)
	}

	links := make([]space.Link, len(t.elinks[uid.Value()]))

	i := 0
	for _, link := range t.elinks[uid.Value()] {
		links[i] = link
		i++
	}

	return links, nil
}

// getNamespaceKindEntities returns all entities in given namespace with given kind matching query q.
func (t Top) getNamespaceKindEntities(ns, kind string, q query.Query) ([]space.Entity, error) {
	if m := q.Matcher(query.Name); m != nil {
		name, ok := m.Predicate().Value().(string)
		if ok && len(name) > 0 {
			o, ok := t.index[ns][kind][name]
			if !ok {
				return []space.Entity{}, nil
			}
			return []space.Entity{o}, nil
		}
	}

	entities := make([]space.Entity, len(t.index[ns][kind]))

	i := 0
	for _, o := range t.index[ns][kind] {
		entities[i] = o
		i++
	}

	return entities, nil
}

// getNamespaceEntities returns all entities in namespaces ns matching given query
func (t Top) getNamespaceEntities(ns string, q query.Query) ([]space.Entity, error) {
	if m := q.Matcher(query.Kind); m != nil {
		kind, ok := m.Predicate().Value().(string)
		if ok && len(kind) > 0 {
			return t.getNamespaceKindEntities(ns, kind, q)
		}
	}

	// nolint:prealloc
	var entities []space.Entity
	for kind := range t.index[ns] {
		ents, err := t.getNamespaceKindEntities(ns, kind, q)
		if err != nil {
			return nil, err
		}
		entities = append(entities, ents...)
	}

	return entities, nil
}

// getAllNamespacedEntities returns all entities from all namespaces
func (t Top) getAllNamespacedEntities(q query.Query) ([]space.Entity, error) {
	// nolint:prealloc
	var entities []space.Entity

	for ns := range t.index {
		ents, err := t.getNamespaceEntities(ns, q)
		if err != nil {
			return nil, err
		}
		entities = append(entities, ents...)
	}

	return entities, nil
}

// Get queries the mapped entities and returns the results
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

	if m := q.Matcher(query.Namespace); m != nil {
		ns, ok := m.Predicate().Value().(string)
		if ok && len(ns) > 0 {
			return t.getNamespaceEntities(ns, q)
		}
	}

	return t.getAllNamespacedEntities(q)
}
