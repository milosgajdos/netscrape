package simple

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/scraper"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Simple is an in-memory scraper plan.
type Simple struct {
	// index of resources by their UID
	index map[string]space.Resource
	// mu synchronizes access to Plan
	mu *sync.RWMutex
}

// NewSimple creates a new Plan and returns it.
func NewSimple(opts ...scraper.Option) (*Simple, error) {
	popts := scraper.Options{}
	for _, apply := range opts {
		apply(&popts)
	}

	return &Simple{
		index: make(map[string]space.Resource),
		mu:    &sync.RWMutex{},
	}, nil
}

func (p *Simple) add(ctx context.Context, r space.Resource, opts ...scraper.Option) error {
	p.index[r.UID().String()] = r
	return nil
}

// Add adds r to plan.
func (p *Simple) Add(ctx context.Context, r space.Resource, opts ...scraper.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.add(ctx, r, opts...)
}

func (p Simple) getAll(ctx context.Context, opts ...scraper.Option) ([]space.Resource, error) {
	res := make([]space.Resource, len(p.index))

	i := 0
	for _, r := range p.index {
		res[i] = r
		i++
	}
	return res, nil
}

// GetAll returns all resource
func (p *Simple) GetAll(ctx context.Context, opts ...scraper.Option) ([]space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.getAll(ctx, opts...)
}

func (p *Simple) get(ctx context.Context, uid uuid.UID, opts ...scraper.Option) (space.Resource, error) {
	r, ok := p.index[uid.String()]
	if !ok {
		return nil, scraper.ErrResourceNotFound
	}
	return r, nil
}

// Get returns the resource with the given uid.
func (p *Simple) Get(ctx context.Context, uid uuid.UID, opts ...scraper.Option) (space.Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.get(ctx, uid, opts...)
}

func (p *Simple) delete(ctx context.Context, uid uuid.UID, opts ...scraper.Option) error {
	delete(p.index, uid.String())
	return nil
}

// Delete removes resource with the given uid from the plan.
func (p *Simple) Delete(ctx context.Context, uid uuid.UID, opts ...scraper.Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.delete(ctx, uid, opts...)
}
