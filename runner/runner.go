package runner

import (
	"context"

	"github.com/milosgajdos/netscrape"
	memgraph "github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/plan"
	memplan "github.com/milosgajdos/netscrape/pkg/plan/memory"
	"github.com/milosgajdos/netscrape/pkg/store"
	memstore "github.com/milosgajdos/netscrape/pkg/store/memory"
)

// Runner scrapes data into networks.
type Runner struct {
	s store.Store
	p plan.Plan
}

// New creates a new netscraper and returns it.
// By default an in-memory store is created backed by
// memory.WUG (in-memory Weighted Undirected Graph).
func New(opts ...netscrape.Option) (*Runner, error) {
	ropts := netscrape.Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	s := ropts.Store
	if s == nil {
		var err error

		g, err := memgraph.NewWUG()
		if err != nil {
			return nil, err
		}

		s, err = memstore.NewStore(memstore.WithGraph(g))
		if err != nil {
			return nil, err
		}
	}

	p := ropts.Plan
	if p == nil {
		var err error

		p, err = memplan.New()
		if err != nil {
			return nil, err
		}
	}

	return &Runner{
		s: s,
		p: p,
	}, nil
}

// Run runs netscraping using scraper s following the plan p.
func (r *Runner) Run(ctx context.Context, s netscrape.Scraper, opts ...netscrape.Option) error {
	return s.Scrape(ctx, r.p, r.s, opts...)
}

// Store returns store handle.
func (r *Runner) Store() store.Store {
	return r.s
}
