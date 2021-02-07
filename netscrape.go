package netscrape

import (
	"context"
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
	// simple loop for the sake of better performance
	for _, f := range n.fx {
		if f(o) {
			return true
		}
	}

	return false
}

func (n *netscraper) linkEntities(ctx context.Context, g graph.Graph, from, to space.Entity, opts ...graph.Option) error {
	gl, ok := g.(graph.NodeLinker)
	if !ok {
		return fmt.Errorf("link entities: %w", graph.ErrUnsupported)
	}

	if _, err := gl.Link(ctx, from.UID(), to.UID(), opts...); err != nil {
		return err
	}

	return nil
}

func (n *netscraper) addEntity(ctx context.Context, g graph.Graph, e space.Entity) error {
	ga, ok := g.(graph.NodeAdder)
	if !ok {
		return fmt.Errorf("add entity: %w", graph.ErrUnsupported)
	}

	from, err := ga.NewNode(ctx, e)
	if err != nil {
		return fmt.Errorf("new node: %v", err)
	}

	if err := n.s.Add(ctx, from); err != nil {
		return fmt.Errorf("store node: %w", err)
	}

	return nil
}

// link links entity e with its topology peers.
func (n *netscraper) link(ctx context.Context, g graph.Graph, e space.Entity, peers []space.Entity, opts ...graph.Option) error {
	if err := n.addEntity(ctx, g, e); err != nil {
		return err
	}

	for _, peer := range peers {
		if err := n.addEntity(ctx, g, peer); err != nil {
			return err
		}

		if err := n.linkEntities(ctx, g, e, peer, opts...); err != nil {
			return err
		}
	}

	return nil
}

// buildGraph builds a graph from given topology skipping entities that match filters.
func (n *netscraper) buildGraph(ctx context.Context, top space.Top, fx ...Filter) error {
	g, err := n.s.Graph(ctx)
	if err != nil {
		return err
	}

	entities, err := top.Entities(ctx)
	if err != nil {
		return err
	}

	for _, ent := range entities {
		if n.skip(ent, fx...) {
			continue
		}

		if err := n.addEntity(ctx, g, ent); err != nil {
			return err
		}

		links, err := top.Links(ctx, ent.UID())
		if err != nil {
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

			if err := n.link(ctx, g, ent, peers, graph.WithAttrs(a)); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run runs netscraping using scraper s on the origin o with filters fx.
// It first creates a space.Plan for the given origin and then maps it into space Topology.
// The topology is used for building a graph which is stored in the configured store.
func (n *netscraper) Run(ctx context.Context, s space.Scraper, o space.Origin, fx ...Filter) error {
	plan, err := s.Plan(ctx, o)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}

	top, err := s.Map(ctx, plan)
	if err != nil {
		return fmt.Errorf("map: %w", err)
	}

	return n.buildGraph(ctx, top, fx...)
}

// Store returns Store handle.
func (n *netscraper) Store() store.Store {
	return n.s
}
