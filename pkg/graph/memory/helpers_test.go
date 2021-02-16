package memory

import (
	"context"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/space/object"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/space/types"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	nodeResName    = "nodeResName"
	nodeResGroup   = "nodeResGroup"
	nodeResVersion = "nodeResVersion"
	nodeResKind    = "nodeResKind"
	nodeGID        = 123
	nodeID         = "testID"
	nodeName       = "testName"
	nodeNs         = "testNs"
)

func newTestResource(name, group, version, kind string, namespaced bool, opts ...resource.Option) (space.Resource, error) {
	return resource.New(name, group, version, kind, namespaced, opts...)
}

func newTestEntity(uid, name, ns string, res space.Resource, opts ...object.Option) (space.Object, error) {
	u, err := uuid.NewFromString(uid)
	if err != nil {
		return nil, err
	}

	opts = append(opts, object.WithUID(u))

	return object.New(name, ns, res, opts...)
}

type testSpace struct {
	entities map[string]space.Object
	links    map[string]map[string]space.Link
}

func makeTestSpace(path string) (*testSpace, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var testEntities []types.Object
	if err := yaml.Unmarshal(data, &testEntities); err != nil {
		return nil, err
	}

	entities := make(map[string]space.Object)
	elinks := make(map[string]map[string]space.Link)

	for _, e := range testEntities {
		a, err := attrs.NewFromMap(e.Resource.Attrs)
		if err != nil {
			return nil, err
		}

		res, err := resource.New(
			e.Resource.Name,
			e.Resource.Group,
			e.Resource.Version,
			e.Resource.Kind,
			e.Resource.Namespaced,
			resource.WithAttrs(a),
		)
		if err != nil {
			return nil, err
		}

		a, err = attrs.NewFromMap(e.Attrs)
		if err != nil {
			return nil, err
		}

		uid, err := uuid.NewFromString(e.UID)
		if err != nil {
			return nil, err
		}

		ent, err := object.New(e.Name, e.Namespace, res, object.WithUID(uid), object.WithAttrs(a))
		if err != nil {
			return nil, err
		}

		for _, l := range e.Links {
			to, err := uuid.NewFromString(l.To)
			if err != nil {
				return nil, err
			}

			a, err = attrs.NewFromMap(l.Attrs)
			if err != nil {
				return nil, err
			}

			if elinks[uid.Value()] == nil {
				elinks[uid.Value()] = make(map[string]space.Link)
			}

			if _, ok := elinks[uid.Value()][to.Value()]; !ok {
				link, err := link.New(uid, to, link.WithAttrs(a))
				if err != nil {
					return nil, err
				}

				elinks[uid.Value()][to.Value()] = link
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

	t, err := makeTestSpace(path)
	if err != nil {
		return nil, err
	}

	for _, ent := range t.entities {
		n, err := g.NewNode(context.TODO(), ent)
		if err != nil {
			return nil, err
		}

		if err := g.AddNode(context.TODO(), n); err != nil {
			return nil, err
		}

		for _, link := range t.links[ent.UID().Value()] {
			ent2, ok := t.entities[link.To().Value()]
			if !ok {
				continue
			}

			n2, err := g.NewNode(context.TODO(), ent2)
			if err != nil {
				return nil, err
			}

			if err := g.AddNode(context.TODO(), n2); err != nil {
				return nil, err
			}

			a := attrs.NewCopyFrom(link.Attrs())

			if _, err = g.Link(context.TODO(), n.UID(), n2.UID(), graph.WithAttrs(a)); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}
