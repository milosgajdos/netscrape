package netscrape

import (
	"context"
	"errors"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	memgraph "github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
	memstore "github.com/milosgajdos/netscrape/pkg/store/memory"
)

// netscraper scrapes data into networks
type netscraper struct {
	s  store.Store
	fx []Filter
}

// New creates a new netscraper and returns it.
// If no store option is provided a memory store
// backed by memory.WUG (Weighted Undirected Graph)
// is used.
func New(opts ...Option) (*netscraper, error) {
	nopts := Options{}
	for _, apply := range opts {
		apply(&nopts)
	}

	s := nopts.Store
	if s == nil {
		var err error

		g, err := memgraph.NewWUG()
		if err != nil {
			return nil, err
		}

		s, err = memstore.New(g)
		if err != nil {
			return nil, err
		}
	}

	return &netscraper{
		s:  s,
		fx: nopts.Filters,
	}, nil
}

// skip returns true if o matches any of the filters.
func (n netscraper) skip(o space.Entity, fx ...Filter) bool {
	for _, f := range fx {
		if f(o) {
			return true
		}
	}

	// NOTE: we avoid appending n.fx to fx and iterating in
	// single loop for the sake of better performance
	for _, f := range n.fx {
		if f(o) {
			return true
		}
	}

	return false
}

// buildNetwork builds a network from given topology top skipping entities that match filters fx.
func (n *netscraper) buildNetwork(ctx context.Context, top space.Top, fx ...Filter) error {
	entities, err := top.Entities(ctx)
	if err != nil {
		return err
	}

	for _, ent := range entities {
		if n.skip(ent, fx...) {
			continue
		}

		if err := n.s.Add(ctx, ent); err != nil {
			return fmt.Errorf("store entity: %w", err)
		}

		links, err := top.Links(ctx, ent.UID())
		// don't return error if there are no outgoing links from ent
		if err != nil && !errors.Is(err, space.ErrEntityNotFound) {
			return err
		}

		for _, link := range links {
			uid := link.To()

			q := base.Build().Add(predicate.UID(uid))

			// NOTE: this must return a single node
			peers, err := top.Get(ctx, q)
			if err != nil {
				return err
			}

			a := attrs.NewCopyFrom(link.Attrs())
			if w := a.Get("weight"); w == "" {
				a.Set("weight", fmt.Sprintf("%f", graph.DefaultWeight))
			}

			for _, peer := range peers {
				if err := n.s.Add(ctx, peer); err != nil {
					return fmt.Errorf("store peer: %w", err)
				}

				if err := n.s.Link(ctx, ent.UID(), peer.UID(), store.WithAttrs(a)); err != nil {
					return fmt.Errorf("link peers: %w", err)
				}
			}
		}
	}

	return nil
}

// Run runs netscraping using scraper s on origin o with filters fx.
// It first creates a space.Plan for the given origin and then maps it into space.Top.
// The topology is used for building a graph which is stored in configured store.
func (n *netscraper) Run(ctx context.Context, s space.Scraper, o space.Origin, fx ...Filter) error {
	plan, err := s.Plan(ctx, o)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}

	top, err := s.Map(ctx, plan)
	if err != nil {
		return fmt.Errorf("map: %w", err)
	}

	return n.buildNetwork(ctx, top, fx...)
}

// Store returns store handle.
func (n *netscraper) Store() store.Store {
	return n.s
}
