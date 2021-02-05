package top

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Top is generic Space topology
type Top struct {
	// space is the source of topology
	space space.Plan
	// entities stores all entities by their UID
	entities map[string]space.Entity
	// index is topology "search index" (ns/kind/name)
	index map[string]map[string]map[string]space.Entity
	// mu synchronizes access to Top
	mu *sync.RWMutex
}

// New creates a new topology and returns it.
func New(a space.Plan) (*Top, error) {
	return &Top{
		space:    a,
		entities: make(map[string]space.Entity),
		index:    make(map[string]map[string]map[string]space.Entity),
		mu:       &sync.RWMutex{},
	}, nil
}

// Plan returns topology Plan.
func (t Top) Plan(ctx context.Context) (space.Plan, error) {
	return t.space, nil
}

// Entities returns all entities in space topology.
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
func (t *Top) Add(ctx context.Context, o space.Entity, opts ...space.Option) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.entities[o.UID().Value()]; !ok {
		if err := t.add(o); err != nil {
			return err
		}
	}

	topObj := t.entities[o.UID().Value()]

	// topObj and o have the same UID so we need to
	// update topObj links with all the o links
	for _, l := range o.Links() {
		if err := topObj.Link(l.To(), space.WithAttrs(l.Attrs()), space.WithMerge(true)); err != nil {
			return err
		}
	}

	return nil
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
		objs, err := t.getNamespaceKindEntities(ns, kind, q)
		if err != nil {
			return nil, err
		}
		entities = append(entities, objs...)
	}

	return entities, nil
}

// getAllNamespacedEntities returns all entities from all namespaces
func (t Top) getAllNamespacedEntities(q query.Query) ([]space.Entity, error) {
	// nolint:prealloc
	var entities []space.Entity

	for ns := range t.index {
		objs, err := t.getNamespaceEntities(ns, q)
		if err != nil {
			return nil, err
		}
		entities = append(entities, objs...)
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
