package netscrape

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/plan"
)

// Scraper scrapes data.
type Scraper interface {
	// Scrape scrapes data following the given plan.
	Scrape(context.Context, plan.Plan, ...Option) error
}
