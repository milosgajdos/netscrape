package simple

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/plan"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Simple is an in-memory scraper plan.
type Simple struct {
	// index of resources by their UID
	index map[string]plan.Resource
	// mu synchronizes access to Plan
	mu *sync.RWMutex
}

// NewSimple creates a new Plan and returns it.
func NewSimple(opts ...plan.Option) (*Simple, error) {
	popts := plan.Options{}
	for _, apply := range opts {
		apply(&popts)
	}

	return &Simple{
		index: make(map[string]plan.Resource),
		mu:    &sync.RWMutex{},
	}, nil
}

func (p *Simple) add(ctx context.Context, r plan.Resource, opts ...plan.Option) error {
	p.index[r.UID().String()] = r
	return nil
}

// Add adds r to plan.
func (p *Simple) Add(ctx context.Context, r plan.Resource, opts ...plan.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.add(ctx, r, opts...)
}

func (p Simple) getAll(ctx context.Context, opts ...plan.Option) ([]plan.Resource, error) {
	res := make([]plan.Resource, len(p.index))

	i := 0
	for _, r := range p.index {
		res[i] = r
		i++
	}
	return res, nil
}

// GetAll returns all resource
func (p *Simple) GetAll(ctx context.Context, opts ...plan.Option) ([]plan.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.getAll(ctx, opts...)
}

func (p *Simple) get(ctx context.Context, uid uuid.UID, opts ...plan.Option) (plan.Resource, error) {
	r, ok := p.index[uid.String()]
	if !ok {
		return nil, plan.ErrResourceNotFound
	}
	return r, nil
}

// Get returns the resource with the given uid.
func (p *Simple) Get(ctx context.Context, uid uuid.UID, opts ...plan.Option) (plan.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.get(ctx, uid, opts...)
}

func (p *Simple) delete(ctx context.Context, uid uuid.UID, opts ...plan.Option) error {
	delete(p.index, uid.String())
	return nil
}

// Delete removes resource with the given uid from the plan.
func (p *Simple) Delete(ctx context.Context, uid uuid.UID, opts ...plan.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.delete(ctx, uid, opts...)
}
