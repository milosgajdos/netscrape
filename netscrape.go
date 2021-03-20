package netscrape

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/plan"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Scraper scrapes data into graph.
type Scraper interface {
	// Scrap scrapes data followgin the given plan and stores it in the given store.
	Scrape(context.Context, plan.Plan, store.Store, ...Option) error
}

// Runner runs the netscraping.
type Runner interface {
	// Run runs netscraping with the given scraper following the given plan.
	Run(context.Context, Scraper, ...Option) error
}
