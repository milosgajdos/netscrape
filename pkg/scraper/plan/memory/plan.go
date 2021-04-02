package memory

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/scraper/plan"
	"github.com/milosgajdos/netscrape/pkg/scraper/plan/origin"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	unknownOrigin = "origin://unknown"
)

// Plan is a space Plan.
type Plan struct {
	// origin is space origin
	origin plan.Origin
	// index of resources by their UID
	index map[string]space.Resource
	// mu synchronizes access to Plan
	mu *sync.RWMutex
}

// New creates a new Plan and returns it.
func New(opts ...plan.Option) (*Plan, error) {
	popts := plan.Options{}
	for _, apply := range opts {
		apply(&popts)
	}

	o := popts.Origin
	if o == nil {
		var err error

		o, err = origin.New(unknownOrigin)
		if err != nil {
			return nil, err
		}
	}

	return &Plan{
		origin: o,
		index:  make(map[string]space.Resource),
		mu:     &sync.RWMutex{},
	}, nil
}

// Origin returns space origin.
func (p Plan) Origin(ctx context.Context) (plan.Origin, error) {
	return p.origin, nil
}

func (p *Plan) add(ctx context.Context, r space.Resource, opts ...plan.Option) error {
	p.index[r.UID().String()] = r
	return nil
}

// Add adds r to plan.
func (p *Plan) Add(ctx context.Context, r space.Resource, opts ...plan.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.add(ctx, r, opts...)
}

func (p Plan) getAll(ctx context.Context, opts ...plan.Option) ([]space.Resource, error) {
	res := make([]space.Resource, len(p.index))

	i := 0
	for _, r := range p.index {
		res[i] = r
		i++
	}
	return res, nil
}

// GetAll returns all resource
func (p *Plan) GetAll(ctx context.Context, opts ...plan.Option) ([]space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.getAll(ctx, opts...)
}

func (p *Plan) get(ctx context.Context, uid uuid.UID, opts ...plan.Option) (space.Resource, error) {
	r, ok := p.index[uid.String()]
	if !ok {
		return nil, plan.ErrResourceNotFound
	}
	return r, nil
}

// Get returns the resource with the given uid.
func (p *Plan) Get(ctx context.Context, uid uuid.UID, opts ...plan.Option) (space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.get(ctx, uid, opts...)
}

func (p *Plan) delete(ctx context.Context, uid uuid.UID, opts ...plan.Option) error {
	delete(p.index, uid.String())
	return nil
}

// Delete removes resource with the given uid from the plan.
func (p *Plan) Delete(ctx context.Context, uid uuid.UID, opts ...plan.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.delete(ctx, uid, opts...)
}
