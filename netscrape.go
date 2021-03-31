package netscrape

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/scraper"
)

// Runner runs netscraping.
type Runner struct {
	opts Options
}

// NewRunner creates a new netscraper runner and returns it.
func NewRunner(opts ...Option) (*Runner, error) {
	ropts := Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	return &Runner{
		opts: ropts,
	}, nil
}

// Run runs netscraping using scraper s following the plan p.
func (r *Runner) Run(ctx context.Context, s scraper.Scraper, opts ...Option) error {
	for _, apply := range opts {
		apply(&r.opts)
	}

	if r.opts.Plan == nil {
		return ErrMissingPlan
	}

	var scraperOpts []scraper.Option

	if r.opts.Broker != nil {
		scraperOpts = append(scraperOpts, scraper.WithBroker(r.opts.Broker))
	}

	if r.opts.Store != nil {
		scraperOpts = append(scraperOpts, scraper.WithStore(r.opts.Store))
	}

	return s.Scrape(ctx, r.opts.Plan, scraperOpts...)
}
