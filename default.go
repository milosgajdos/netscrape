package netscrape

import (
	"context"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	memgraph "github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
	memstore "github.com/milosgajdos/netscrape/pkg/store/memory"
)

type netscraper struct {
	s store.Store
}

// New creates a new netscraper and returns it.
// If no store option has been provided it uses memory store
// backed by memory WUG (Weighted Undirected Graph).
func New(opts ...Option) (NetScraper, error) {
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
		s: s,
	}, nil
}

// skip returns true if o matches any of the filters.
func skip(o space.Object, filters ...Filter) bool {
	if len(filters) == 0 {
		return false
	}

	for _, filter := range filters {
		if filter(o) {
			return false
		}
	}

	return true
}

func (n *netscraper) linkObjectNodes(ctx context.Context, from, to space.Object, opts ...graph.Option) error {
	g, err := n.s.Graph(ctx)
	if err != nil {
		return err
	}

	gl := g.(graph.NodeLinker)

	if _, err := gl.Link(ctx, from.UID(), to.UID(), opts...); err != nil {
		return err
	}

	return nil
}

func (n *netscraper) addObjectNode(ctx context.Context, o space.Object) error {
	g, err := n.s.Graph(ctx)
	if err != nil {
		return err
	}

	ga := g.(graph.NodeAdder)

	from, err := ga.NewNode(ctx, o)
	if err != nil {
		return fmt.Errorf("create node: %v", err)
	}

	if err := n.s.Add(ctx, from); err != nil {
		return fmt.Errorf("add node: %w", err)
	}

	return nil
}

func (n *netscraper) link(ctx context.Context, o space.Object, peers []space.Object, opts ...graph.Option) error {
	if err := n.addObjectNode(ctx, o); err != nil {
		return err
	}

	for _, peer := range peers {
		if err := n.addObjectNode(ctx, peer); err != nil {
			return err
		}

		if err := n.linkObjectNodes(ctx, o, peer, opts...); err != nil {
			return err
		}
	}

	return nil
}

// buildGraph builds a graph from given topology.
// It skips adding nodes that match any of the passed in filters.
func (n *netscraper) buildGraph(ctx context.Context, top space.Top, filters ...Filter) error {
	objects, err := top.Objects(ctx)
	if err != nil {
		return err
	}

	for _, object := range objects {
		if skip(object, filters...) {
			continue
		}

		if len(object.Links()) == 0 {
			if err := n.addObjectNode(ctx, object); err != nil {
				return err
			}

			continue
		}

		for _, link := range object.Links() {
			uid := link.To()

			q := base.Build().
				Add(query.UID(uid), query.UUIDEqFunc(uid))

			// NOTE: this must return a single node
			peers, err := top.Get(ctx, q)
			if err != nil {
				return err
			}

			a := attrs.NewCopyFrom(link.Attrs())
			if w := a.Get("weight"); w == "" {
				a.Set("weight", fmt.Sprintf("%f", graph.DefaultWeight))
			}

			if err := n.link(ctx, object, peers, graph.WithAttrs(a)); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run runs netscraping and returns error if it fails.
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
