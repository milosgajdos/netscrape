package plan

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
)

// Plan is a space Plan.
type Plan struct {
	// origin is Space origin
	origin space.Origin
	// resources indexes discovered Space resources
	// the index follows this pattern: group/version/kind
	resources map[string]map[string]map[string]space.Resource
	// mu synchronizes access to Space
	mu *sync.RWMutex
}

// New creates a new Plan and returns it.
func New(o space.Origin) (*Plan, error) {
	return &Plan{
		origin:    o,
		resources: make(map[string]map[string]map[string]space.Resource),
		mu:        &sync.RWMutex{},
	}, nil
}

// Origin returns Space origin.
func (a Plan) Origin(ctx context.Context) (space.Origin, error) {
	return a.origin, nil
}

// Add adds r to Space.
func (a *Plan) Add(ctx context.Context, r space.Resource, opts space.AddOptions) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	group := r.Group()

	if a.resources[group] == nil {
		a.resources[group] = make(map[string]map[string]space.Resource)
	}

	version := r.Version()

	if a.resources[group][version] == nil {
		a.resources[group][version] = make(map[string]space.Resource)
	}

	kind := r.Kind()

	a.resources[group][version][kind] = r

	return nil
}

// Resources returns all Space resources.
func (a Plan) Resources(ctx context.Context) ([]space.Resource, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var resources []space.Resource

	for _, groups := range a.resources {
		for _, versions := range groups {
			for _, r := range versions {
				resources = append(resources, r)
			}
		}
	}

	return resources, nil
}

func matchName(r space.Resource, q query.Query) bool {
	if m := q.Matcher(query.PName); m != nil {
		name, ok := m.Predicate().Value().(string)
		if ok && len(name) > 0 {
			return name == r.Name()
		}
	}

	// NOTE: if missing the name we assume Any name
	return true
}

func (a Plan) getGroupVersionResources(group, version string, q query.Query) ([]space.Resource, error) {
	if m := q.Matcher(query.PKind); m != nil {
		kind, ok := m.Predicate().Value().(string)
		if ok && len(kind) > 0 {
			r, ok := a.resources[group][version][kind]
			if !ok {
				return []space.Resource{}, nil
			}

			if !matchName(r, q) {
				return []space.Resource{}, nil
			}

			return []space.Resource{r}, nil
		}
	}

	var resources []space.Resource

	// NOTE: missing Kind matcher implies ANY kind
	for kind := range a.resources[group][version] {
		r := a.resources[group][version][kind]
		if matchName(r, q) {
			resources = append(resources, r)
		}
	}

	return resources, nil
}

// getGroupResources returns all resource in group g matching q.
func (a Plan) getGroupResources(g string, q query.Query) ([]space.Resource, error) {
	if m := q.Matcher(query.PVersion); m != nil {
		v, ok := m.Predicate().Value().(string)
		if ok && len(v) > 0 {
			return a.getGroupVersionResources(g, v, q)
		}
	}

	// nolint:prealloc
	var resources []space.Resource

	// NOTE: missing Version matcher implies ANY version
	for v := range a.resources[g] {
		rx, err := a.getGroupVersionResources(g, v, q)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rx...)
	}

	return resources, nil
}

// getAllGroupedResources returns all Resources in all groups matching q.
func (a Plan) getAllGroupedResources(q query.Query) ([]space.Resource, error) {
	// nolint:prealloc
	var resources []space.Resource

	for g := range a.resources {
		rx, err := a.getGroupResources(g, q)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rx...)
	}

	return resources, nil
}

// Get returns all resources matching the given query.
func (a Plan) Get(ctx context.Context, q query.Query) ([]space.Resource, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if m := q.Matcher(query.PGroup); m != nil {
		g, ok := m.Predicate().Value().(string)
		if ok && len(g) > 0 {
			return a.getGroupResources(g, q)
		}
	}

	return a.getAllGroupedResources(q)
}
