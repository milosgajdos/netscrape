package top

import (
	"sync"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Top is generic Space topology
type Top struct {
	// space is the source of topology
	space space.Plan
	// objects stores all objects by their UID
	objects map[string]space.Object
	// index is topology "search index" (ns/kind/name)
	index map[string]map[string]map[string]space.Object
	// mu synchronizes access to Top
	mu *sync.RWMutex
}

// New creates a new topology and returns it.
func New(a space.Plan) (*Top, error) {
	return &Top{
		space:   a,
		objects: make(map[string]space.Object),
		index:   make(map[string]map[string]map[string]space.Object),
		mu:      &sync.RWMutex{},
	}, nil
}

// Plan returns topology Plan.
func (t Top) Plan() space.Plan {
	return t.space
}

// Objects returns all space objects in tpoology
func (t Top) Objects() []space.Object {
	t.mu.RLock()
	defer t.mu.RUnlock()

	objects := make([]space.Object, len(t.objects))

	i := 0

	for _, object := range t.objects {
		objects[i] = object
		i++
	}

	return objects
}

// add adds a new object to topology
func (t *Top) add(o space.Object) error {
	t.objects[o.UID().Value()] = o

	ns := o.Namespace()

	if t.index[ns] == nil {
		t.index[ns] = make(map[string]map[string]space.Object)
	}

	kind := o.Resource().Kind()

	if t.index[ns][kind] == nil {
		t.index[ns][kind] = make(map[string]space.Object)
	}

	name := o.Name()

	t.index[ns][kind][name] = o

	return nil
}

// Add adds o to topology with the given options.
// If an object already exists in topology and MergeLinks option is enabled
// the existing object links are merged with the links of o.
func (t *Top) Add(o space.Object, opts space.AddOptions) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.objects[o.UID().Value()]; !ok {
		if err := t.add(o); err != nil {
			return err
		}

		for _, l := range o.Links() {
			lopts := space.LinkOptions{
				Merge:    opts.MergeLinks,
				Metadata: l.Metadata(),
			}

			if to, ok := t.objects[l.To().Value()]; ok {
				if err := to.Link(o.UID(), lopts); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return nil
}

// getNamespaceKindObjects returns all objects in given namespace with given kind matching query q.
func (t Top) getNamespaceKindObjects(ns, kind string, q query.Query) ([]space.Object, error) {
	if m := q.Matcher(query.PName); m != nil {
		name, ok := m.Predicate().Value().(string)
		if ok && len(name) > 0 {
			o, ok := t.index[ns][kind][name]
			if !ok {
				return []space.Object{}, nil
			}
			return []space.Object{o}, nil
		}
	}

	objects := make([]space.Object, len(t.index[ns][kind]))

	i := 0
	for _, o := range t.index[ns][kind] {
		objects[i] = o
		i++
	}

	return objects, nil
}

// getNamespaceObjects returns all objects in namespaces ns matching given query
func (t Top) getNamespaceObjects(ns string, q query.Query) ([]space.Object, error) {
	if m := q.Matcher(query.PKind); m != nil {
		kind, ok := m.Predicate().Value().(string)
		if ok && len(kind) > 0 {
			return t.getNamespaceKindObjects(ns, kind, q)
		}
	}

	// nolint:prealloc
	var objects []space.Object
	for kind := range t.index[ns] {
		objs, err := t.getNamespaceKindObjects(ns, kind, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

// getAllNamespacedObjects returns all objects from all namespaces
func (t Top) getAllNamespacedObjects(q query.Query) ([]space.Object, error) {
	// nolint:prealloc
	var objects []space.Object

	for ns := range t.index {
		objs, err := t.getNamespaceObjects(ns, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

// Get queries the mapped objects and returns the results
func (t Top) Get(q query.Query) ([]space.Object, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if m := q.Matcher(query.PUID); m != nil {
		uid, ok := m.Predicate().Value().(uuid.UID)
		if ok && len(uid.Value()) > 0 {
			o, ok := t.objects[uid.Value()]
			if !ok {
				return []space.Object{}, nil
			}
			return []space.Object{o}, nil
		}
	}

	if m := q.Matcher(query.PNamespace); m != nil {
		ns, ok := m.Predicate().Value().(string)
		if ok && len(ns) > 0 {
			return t.getNamespaceObjects(ns, q)
		}
	}

	return t.getAllNamespacedObjects(q)
}
