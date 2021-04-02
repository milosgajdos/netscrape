package scraper

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/scraper/plan"
)

// Scraper scrapes data.
type Scraper interface {
	// Scrape scrapes data following the provided plan.
	Scrape(context.Context, plan.Plan, ...Option) error
}
