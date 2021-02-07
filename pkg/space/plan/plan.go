package plan

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/space"
)

// Plan is a space Plan.
type Plan struct {
	// origin is space origin
	origin space.Origin
	// resources indexes discovered space resources
	// the index follows this pattern: group/version/kind
	resources map[string]map[string]map[string]space.Resource
	// mu synchronizes access to Plan
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

// Origin returns space origin.
func (p Plan) Origin(ctx context.Context) (space.Origin, error) {
	return p.origin, nil
}

// Add adds r to Plan
func (p *Plan) Add(ctx context.Context, r space.Resource, opts ...space.Option) error {
	oopts := space.Options{}
	for _, apply := range opts {
		apply(&oopts)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	group := r.Group()

	if p.resources[group] == nil {
		p.resources[group] = make(map[string]map[string]space.Resource)
	}

	version := r.Version()

	if p.resources[group][version] == nil {
		p.resources[group][version] = make(map[string]space.Resource)
	}

	kind := r.Kind()

	p.resources[group][version][kind] = r

	return nil
}

// Remove removes Resource from Plan.
func (p *Plan) Remove(ctx context.Context, r space.Resource, opts ...space.Option) error {
	oopts := space.Options{}
	for _, apply := range opts {
		apply(&oopts)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	group := r.Group()

	if p.resources[group] == nil {
		return nil
	}

	version := r.Version()

	if p.resources[group][version] == nil {
		return nil
	}

	kind := r.Kind()

	if _, ok := p.resources[group][version][kind]; ok {
		delete(p.resources[group][version], kind)
	}

	return nil
}

// Resources returns all Plan resources.
func (p Plan) Resources(ctx context.Context) ([]space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var resources []space.Resource

	for _, groups := range p.resources {
		for _, versions := range groups {
			for _, r := range versions {
				resources = append(resources, r)
			}
		}
	}

	return resources, nil
}

func matchName(r space.Resource, q query.Query) bool {
	if m := q.Matcher(query.Name); m != nil {
		name, ok := m.Predicate().Value().(string)
		if ok && len(name) > 0 {
			return name == r.Name()
		}
	}

	// NOTE: if missing the name we assume Any name
	return true
}

func (p Plan) getGroupVersionResources(group, version string, q query.Query) ([]space.Resource, error) {
	if m := q.Matcher(query.Kind); m != nil {
		kind, ok := m.Predicate().Value().(string)
		if ok && len(kind) > 0 {
			r, ok := p.resources[group][version][kind]
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
	for kind := range p.resources[group][version] {
		r := p.resources[group][version][kind]
		if matchName(r, q) {
			resources = append(resources, r)
		}
	}

	return resources, nil
}

// getGroupResources returns all resource in group g matching q.
func (p Plan) getGroupResources(g string, q query.Query) ([]space.Resource, error) {
	if m := q.Matcher(query.Version); m != nil {
		v, ok := m.Predicate().Value().(string)
		if ok && len(v) > 0 {
			return p.getGroupVersionResources(g, v, q)
		}
	}

	// nolint:prealloc
	var resources []space.Resource

	// NOTE: missing Version matcher implies ANY version
	for v := range p.resources[g] {
		rx, err := p.getGroupVersionResources(g, v, q)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rx...)
	}

	return resources, nil
}

// getAllGroupedResources returns all Resources in all groups matching q.
func (p Plan) getAllGroupedResources(q query.Query) ([]space.Resource, error) {
	// nolint:prealloc
	var resources []space.Resource

	for g := range p.resources {
		rx, err := p.getGroupResources(g, q)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rx...)
	}

	return resources, nil
}

// Get returns all resources matching the given query.
func (p Plan) Get(ctx context.Context, q query.Query) ([]space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if m := q.Matcher(query.Group); m != nil {
		g, ok := m.Predicate().Value().(string)
		if ok && len(g) > 0 {
			return p.getGroupResources(g, q)
		}
	}

	return p.getAllGroupedResources(q)
}
