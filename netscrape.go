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

// New creates a netscraper and returns it.
// If no store option is given memory store is created
// backed by memory WUG (Weighted Undirected Graph).
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
func skip(o space.Object, filters ...Filter) bool {
	for _, filter := range filters {
		if filter(o) {
			return false
		}
	}

	return false
}

func (n *netscraper) linkObjects(ctx context.Context, g graph.Graph, from, to space.Object, opts ...graph.Option) error {
	gl, ok := g.(graph.NodeLinker)
	if !ok {
		return fmt.Errorf("link objects: %w", graph.ErrUnsupported)
	}

	if _, err := gl.Link(ctx, from.UID(), to.UID(), opts...); err != nil {
		return err
	}

	return nil
}

func (n *netscraper) addObject(ctx context.Context, g graph.Graph, o space.Object) error {
	ga, ok := g.(graph.NodeAdder)
	if !ok {
		return fmt.Errorf("add objects: %w", graph.ErrUnsupported)
	}

	from, err := ga.NewNode(ctx, o)
	if err != nil {
		return fmt.Errorf("create node: %v", err)
	}

	if err := n.s.Add(ctx, from); err != nil {
		return fmt.Errorf("add node: %w", err)
	}

	return nil
}

// link links object o with its topology peers.
func (n *netscraper) link(ctx context.Context, g graph.Graph, o space.Object, peers []space.Object, opts ...graph.Option) error {
	if err := n.addObject(ctx, g, o); err != nil {
		return err
	}

	for _, peer := range peers {
		if err := n.addObject(ctx, g, peer); err != nil {
			return err
		}

		if err := n.linkObjects(ctx, g, o, peer, opts...); err != nil {
			return err
		}
	}

	return nil
}

// buildGraph builds a graph from given topology skipping objects that match filters.
func (n *netscraper) buildGraph(ctx context.Context, top space.Top, fx ...Filter) error {
	g, err := n.s.Graph(ctx)
	if err != nil {
		return err
	}

	objects, err := top.Objects(ctx)
	if err != nil {
		return err
	}

	for _, object := range objects {
		if skip(object, append(n.fx, fx...)...) {
			continue
		}

		if len(object.Links()) == 0 {
			if err := n.addObject(ctx, g, object); err != nil {
				return err
			}

			continue
		}

		for _, link := range object.Links() {
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

			if err := n.link(ctx, g, object, peers, graph.WithAttrs(a)); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run runs netscraping using scraper s on the origin o with filters fx.
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
