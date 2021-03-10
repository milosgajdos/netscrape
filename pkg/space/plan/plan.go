package plan

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Plan is a space Plan.
type Plan struct {
	// origin is space origin
	origin space.Origin
	// index of resources by their UID
	index map[string]space.Resource
	// mu synchronizes access to Plan
	mu *sync.RWMutex
}

// New creates a new Plan and returns it.
func New(o space.Origin) (*Plan, error) {
	return &Plan{
		origin: o,
		index:  make(map[string]space.Resource),
		mu:     &sync.RWMutex{},
	}, nil
}

// Origin returns space origin.
func (p Plan) Origin(ctx context.Context) (space.Origin, error) {
	return p.origin, nil
}

func (p *Plan) add(ctx context.Context, r space.Resource, opts ...space.Option) error {
	p.index[r.UID().Value()] = r

	return nil
}

// Add adds r to plan.
func (p *Plan) Add(ctx context.Context, r space.Resource, opts ...space.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.add(ctx, r, opts...)
}

func (p *Plan) getAll(ctx context.Context) ([]space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	res := make([]space.Resource, len(p.index))

	i := 0
	for _, r := range p.index {
		res[i] = r
		i++
	}

	return res, nil
}

// GetAll returns all resource
func (p *Plan) GetAll(ctx context.Context) ([]space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.getAll(ctx)
}

func (p *Plan) get(ctx context.Context, uid uuid.UID, opts ...space.Option) (space.Resource, error) {
	r, ok := p.index[uid.Value()]
	if !ok {
		return nil, space.ErrResourceNotFound
	}

	return r, nil
}

// Get returns the resource with the given uid.
func (p *Plan) Get(ctx context.Context, uid uuid.UID, opts ...space.Option) (space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.get(ctx, uid, opts...)
}

func (p *Plan) delete(ctx context.Context, uid uuid.UID, opts ...space.Option) error {
	delete(p.index, uid.Value())

	return nil
}

// Delete removes resource with the given uid from the plan.
func (p *Plan) Delete(ctx context.Context, uid uuid.UID, opts ...space.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.delete(ctx, uid, opts...)
}
