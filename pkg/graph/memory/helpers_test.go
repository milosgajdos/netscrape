package memory

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/space/marshal"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	nodeResType    = "nodeResType"
	nodeResName    = "nodeResName"
	nodeResGroup   = "nodeResGroup"
	nodeResVersion = "nodeResVersion"
	nodeResKind    = "nodeResKind"
	nodeType       = "testType"
	nodeID         = "testID"
	nodeName       = "testName"
	nodeNs         = "testNs"
)

type testSpace struct {
	entities map[string]space.Entity
	links    map[string]map[string]space.Link
}

func makeTestSpace(path string) (*testSpace, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var testEntities []marshal.LinkedEntity
	if err := yaml.Unmarshal(data, &testEntities); err != nil {
		return nil, err
	}

	entities := make(map[string]space.Entity)
	elinks := make(map[string]map[string]space.Link)

	for _, e := range testEntities {
		resAttrs := memattrs.NewFromMap(e.Resource.Attrs)
		res, err := resource.New(
			e.Resource.Type,
			e.Resource.Name,
			e.Resource.Group,
			e.Resource.Version,
			e.Resource.Kind,
			e.Resource.Namespaced,
			resource.WithAttrs(resAttrs),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create new resource: %v", err)
		}

		a := memattrs.NewFromMap(e.Attrs)

		uid, err := uuid.NewFromString(e.UID)
		if err != nil {
			return nil, err
		}

		ent, err := entity.New(e.Type, e.Name, e.Namespace, res, entity.WithUID(uid), entity.WithAttrs(a))
		if err != nil {
			return nil, err
		}

		for _, l := range e.Links {
			to, err := uuid.NewFromString(l.To)
			if err != nil {
				return nil, err
			}

			if elinks[uid.String()] == nil {
				elinks[uid.String()] = make(map[string]space.Link)
			}

			if _, ok := elinks[uid.String()][to.String()]; !ok {
				a := memattrs.NewFromMap(l.Attrs)
				link, err := link.New(uid, to, link.WithAttrs(a))
				if err != nil {
					return nil, err
				}

				elinks[uid.String()][to.String()] = link
			}
		}

		entities[e.UID] = ent
	}

	return &testSpace{
		entities: entities,
		links:    elinks,
	}, nil
}

func makeTestGraph(path string) (*WUG, error) {
	g, err := NewWUG()
	if err != nil {
		return nil, err
	}

	// NOTE: makeTestSpace builds a map of entities
	// from testdata; these are marshal.LinkedEntity-s
	// so we know they're proper "full"-y initialized entities
	// hence it's ok to do Upsert into graph when adding them in
	t, err := makeTestSpace(path)
	if err != nil {
		return nil, err
	}

	for _, ent := range t.entities {
		n, err := g.NewNode(context.Background(), ent)
		if err != nil {
			return nil, err
		}

		if err := g.AddNode(context.Background(), n, graph.WithUpsert()); err != nil {
			return nil, err
		}

		for _, link := range t.links[ent.UID().String()] {
			ent2, ok := t.entities[link.To().String()]
			if !ok {
				continue
			}

			n2, err := g.NewNode(context.Background(), ent2)
			if err != nil {
				return nil, err
			}

			if err := g.AddNode(context.Background(), n2, graph.WithUpsert()); err != nil {
				return nil, err
			}

			a := memattrs.NewCopyFrom(link.Attrs())

			if _, err = g.Link(context.Background(), n.UID(), n2.UID(), graph.WithAttrs(a)); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}
