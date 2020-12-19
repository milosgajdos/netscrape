package netscrape

import (
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
)

type netscraper struct {
	store store.Store
	opts  Options
}

// New creates a new netscraper and returns it
func New(store store.Store, opts ...Option) (NetScraper, error) {
	o := Options{}
	for _, apply := range opts {
		apply(&o)
	}

	return &netscraper{
		store: store,
		opts:  o,
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

func (n *netscraper) link(o space.Object, neighbs []space.Object, opts graph.LinkOptions) error {
	ga := n.store.Graph().(graph.NodeAdder)
	gl := n.store.Graph().(graph.NodeLinker)

	from, err := ga.NewNode(o, graph.NodeOptions{})
	if err != nil {
		return err
	}

	if err := n.store.Add(from, store.AddOptions{}); err != nil {
		return err
	}

	for _, neighb := range neighbs {
		to, err := ga.NewNode(neighb, graph.NodeOptions{})
		if err != nil {
			return err
		}

		if err := n.store.Add(to, store.AddOptions{}); err != nil {
			return err
		}

		if _, err := gl.Link(from.UID(), to.UID(), opts); err != nil {
			return err
		}
	}

	return nil
}

// buildGraph builds a graph from given topology.
// It skips adding nodes to graph for topology objects which match any of filters.
func (n *netscraper) buildGraph(top space.Top, filters ...Filter) error {
	g := n.store.Graph().(graph.NodeAdder)

	for _, object := range top.Objects() {
		if skip(object, filters...) {
			continue
		}

		if len(object.Links()) == 0 {
			node, err := g.NewNode(object, graph.NodeOptions{})
			if err != nil {
				return fmt.Errorf("faled to create node: %v", err)
			}

			if err := n.store.Add(node, store.AddOptions{}); err != nil {
				return fmt.Errorf("adding node: %w", err)
			}

			continue
		}

		for _, link := range object.Links() {
			uid := link.To()

			q := base.Build().Add(query.UID(uid), query.UUIDEqFunc(uid))

			neighbs, err := top.Get(q)
			if err != nil {
				return err
			}

			a, err := attrs.New()
			if err != nil {
				return err
			}

			w := graph.DefaultWeight

			if weight := link.Metadata().Get("weight"); weight != nil {
				if val, ok := weight.(float64); ok {
					w = val
				}
			}
			a.Set("weight", fmt.Sprintf("%f", w))

			if rel := link.Metadata().Get("relation"); rel != nil {
				if r, ok := rel.(string); ok {
					a.Set("relation", r)
				}
			}

			if err := n.link(object, neighbs, graph.LinkOptions{Attrs: a, Weight: w}); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run runs netscaping and returns error if it fails.
func (n *netscraper) Run(s space.Scraper, o space.Origin, fx ...Filter) error {
	plan, err := s.Plan(o)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}

	top, err := s.Map(plan)
	if err != nil {
		return fmt.Errorf("map: %w", err)
	}

	if _, ok := n.store.Graph().(graph.NodeAdder); !ok {
		return fmt.Errorf("unable to build graph: %w", graph.ErrUnsupported)
	}

	return n.buildGraph(top, fx...)
}

// Store returns Store handle.
func (n *netscraper) Store() store.Store {
	return n.store
}
