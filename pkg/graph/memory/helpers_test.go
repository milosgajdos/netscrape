package memory

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/space/marshal"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func MustSet(ctx context.Context, a attrs.Attrs, k, v string, t *testing.T) {
	if err := a.Set(ctx, k, v); err != nil {
		t.Fatalf("failed to set val %s for key %s: %v", k, v, err)
	}
}

type testSpace struct {
	objects map[string]space.Object
	links   map[string]map[string]space.Link
}

func makeTestSpace(path string) (*testSpace, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var testObjects []marshal.LinkedObject
	if err := yaml.Unmarshal(data, &testObjects); err != nil {
		return nil, err
	}

	objects := make(map[string]space.Object)
	elinks := make(map[string]map[string]space.Link)

	for _, o := range testObjects {
		resAttrs := memattrs.NewFromMap(o.Resource.Attrs)
		res, err := entity.NewResource(
			o.Resource.Type,
			o.Resource.Name,
			o.Resource.Group,
			o.Resource.Version,
			o.Resource.Kind,
			o.Resource.Namespaced,
			entity.WithAttrs(resAttrs),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create new resource: %v", err)
		}

		a := memattrs.NewFromMap(o.Attrs)
		uid := memuid.NewFromString(o.UID)

		obj, err := entity.NewObject(o.Type, o.Name, o.Namespace, res, entity.WithUID(uid), entity.WithAttrs(a))
		if err != nil {
			return nil, err
		}

		for _, l := range o.Links {
			to := memuid.NewFromString(l.To)

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

		objects[o.UID] = obj
	}

	return &testSpace{
		objects: objects,
		links:   elinks,
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

	for _, obj := range t.objects {
		n, err := g.NewNode(context.Background(), obj)
		if err != nil {
			return nil, err
		}

		if err := g.AddNode(context.Background(), n, graph.WithUpsert()); err != nil {
			return nil, err
		}

		for _, link := range t.links[obj.UID().String()] {
			ent2, ok := t.objects[link.To().String()]
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

			a, err := memattrs.NewCopyFrom(context.Background(), link.Attrs())
			if err != nil {
				return nil, err
			}

			if _, err = g.Link(context.Background(), n.UID(), n2.UID(), graph.WithAttrs(a)); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}
