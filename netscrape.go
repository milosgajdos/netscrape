package netscrape

import (
	"context"
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	memgraph "github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/store"
	memstore "github.com/milosgajdos/netscrape/pkg/store/memory"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// netscraper scrapes data into networks
type netscraper struct {
	s store.Store
}

// New creates a new netscraper and returns it.
// By default an in-memory store is created backed by
// memory.WUG (in-memory Weighted Undirected Graph).
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

		s, err = memstore.NewStore(memstore.WithGraph(g))
		if err != nil {
			return nil, err
		}
	}

	return &netscraper{
		s: s,
	}, nil
}

func (n *netscraper) bulkLinkPeers(ctx context.Context, s store.BulkStore, top space.BulkTop) error {
	ents, err := top.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get entities: %w", err)
	}

	storeEnts := make([]store.Entity, len(ents))
	uids := make([]uuid.UID, len(ents))

	for i, e := range ents {
		storeEnts[i] = e
		uids[i] = e.UID()
	}

	if err := s.BulkAdd(ctx, storeEnts); err != nil {
		return fmt.Errorf("bulk store entities : %w", err)
	}

	links, err := top.BulkLinks(ctx, uids)
	if err != nil {
		return fmt.Errorf("get bulk links: %w", err)
	}

	for suid, lx := range links {
		uid, err := uuid.NewFromString(suid)
		if err != nil {
			return err
		}

		to := make([]uuid.UID, len(lx))
		for i, l := range lx {
			to[i] = l.To()
		}

		if err := s.BulkLink(ctx, uid, to); err != nil {
			return fmt.Errorf("store bulk link: %w", err)
		}
	}

	return nil
}

// bulkBuildNetwork is the same as buildNetwork, but stores the nodes in bulks rather than one by one.
func (n *netscraper) bulkBuildNetwork(ctx context.Context, s store.BulkStore, top space.Top) error {
	if bt, ok := top.(space.BulkTop); ok {
		return n.bulkLinkPeers(ctx, s, bt)
	}

	ents, err := top.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get entities: %w", err)
	}

	storeEnts := make([]store.Entity, len(ents))

	for i, e := range ents {
		storeEnts[i] = e
	}

	if err := s.BulkAdd(ctx, storeEnts); err != nil {
		return fmt.Errorf("bulk store entities : %w", err)
	}

	for _, e := range ents {
		lx, err := top.Links(ctx, e.UID())
		if err != nil {
			return err
		}

		to := make([]uuid.UID, len(lx))
		for i, l := range lx {
			to[i] = l.To()
		}

		if err := s.BulkLink(ctx, e.UID(), to); err != nil {
			return fmt.Errorf("store bulk link: %w", err)
		}
	}

	return nil
}

// buildNetwork builds a network from given topology top skipping entities that match filters fx.
func (n *netscraper) buildNetwork(ctx context.Context, top space.Top) error {
	entities, err := top.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get entities: %w", err)
	}

	for _, ent := range entities {
		if err := n.s.Add(ctx, ent); err != nil {
			return fmt.Errorf("store entity: %w", err)
		}

		links, err := top.Links(ctx, ent.UID())
		if err != nil {
			return err
		}

		for _, link := range links {
			peer, err := top.Get(ctx, link.To())
			if err != nil {
				return err
			}

			if err := n.s.Add(ctx, peer); err != nil {
				return fmt.Errorf("store peer: %w", err)
			}

			a := attrs.NewCopyFrom(link.Attrs())
			if w := a.Get(attrs.Weight); w == "" {
				a.Set(attrs.Weight, fmt.Sprintf("%f", graph.DefaultWeight))
			}

			if err := n.s.Link(ctx, ent.UID(), peer.UID(), store.WithAttrs(a)); err != nil {
				return fmt.Errorf("link peers: %w", err)
			}
		}
	}

	return nil
}

// Run runs netscraping using scraper s on origin o with filters fx.
// It first creates a space.Plan for the given origin and then maps it into space.Top.
// The topology is used for building a graph which is stored in configured store.
func (n *netscraper) Run(ctx context.Context, s space.Scraper, o space.Origin) error {
	plan, err := s.Plan(ctx, o)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}

	top, err := s.Map(ctx, plan)
	if err != nil {
		return fmt.Errorf("map: %w", err)
	}

	bs, ok := n.s.(store.BulkStore)
	if ok {
		return n.bulkBuildNetwork(ctx, bs, top)
	}

	return n.buildNetwork(ctx, top)
}

// Store returns store handle.
func (n *netscraper) Store() store.Store {
	return n.s
}
