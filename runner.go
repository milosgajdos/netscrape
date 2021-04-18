package netscrape

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/plan"
)

// Runner runs netscraping.
type Runner struct {
	opts Options
}

// NewRunner creates a new Runner and returns it.
func NewRunner(opts ...Option) (*Runner, error) {
	ropts := Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	return &Runner{
		opts: ropts,
	}, nil
}

// Run runs netscraping using scraper s.
func (r *Runner) Run(ctx context.Context, p plan.Plan, s Scraper, opts ...Option) error {
	sopts := Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	return s.Scrape(ctx, p, WithBroker(sopts.Broker))
}
