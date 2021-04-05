package netscrape

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/scraper"
)

// Runner runs netscraping.
type Runner struct {
	// NOTE: these options are not currently used.
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
func (r *Runner) Run(ctx context.Context, p scraper.Plan, s scraper.Scraper, opts ...Option) error {
	sopts := Options{}
	for _, apply := range opts {
		apply(&sopts)
	}

	return s.Scrape(ctx, p, scraper.WithBroker(sopts.Broker))
}
