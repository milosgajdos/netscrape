package memory

import (
	"context"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
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

func newTestEntity(uid, name, ns string, res space.Resource, opts ...entity.Option) (space.Entity, error) {
	u, err := uuid.NewFromString(uid)
	if err != nil {
		return nil, err
	}

	opts = append(opts, entity.WithUID(u))

	return entity.New(name, ns, res, opts...)
}

func makeTestSpaceEntities(path string) (map[string]space.Entity, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var testEntities []types.Entity
	if err := yaml.Unmarshal(data, &testEntities); err != nil {
		return nil, err
	}

	entities := make(map[string]space.Entity)

	for _, o := range testEntities {
		a, err := attrs.NewFromMap(o.Resource.Attrs)
		if err != nil {
			return nil, err
		}

		res, err := resource.New(
			o.Resource.Name,
			o.Resource.Group,
			o.Resource.Version,
			o.Resource.Kind,
			o.Resource.Namespaced,
			resource.WithAttrs(a),
		)
		if err != nil {
			return nil, err
		}

		a, err = attrs.NewFromMap(o.Attrs)
		if err != nil {
			return nil, err
		}

		uid, err := uuid.NewFromString(o.UID)
		if err != nil {
			return nil, err
		}

		obj, err := entity.New(o.Name, o.Namespace, res, entity.WithUID(uid), entity.WithAttrs(a))
		if err != nil {
			return nil, err
		}

		for _, l := range o.Links {
			toUID, err := uuid.NewFromString(l.To)
			if err != nil {
				return nil, err
			}

			a, err = attrs.NewFromMap(l.Attrs)
			if err != nil {
				return nil, err
			}

			if err := obj.Link(toUID, space.WithAttrs(a)); err != nil {
				return nil, err
			}
		}

		entities[o.UID] = obj
	}

	return entities, nil
}

func makeTestGraph(path string) (*WUG, error) {
	g, err := NewWUG()
	if err != nil {
		return nil, err
	}

	entities, err := makeTestSpaceEntities(path)
	if err != nil {
		return nil, err
	}

	for _, ent := range entities {
		n, err := g.NewNode(context.TODO(), ent)
		if err != nil {
			return nil, err
		}

		if err := g.AddNode(context.TODO(), n); err != nil {
			return nil, err
		}

		for _, link := range ent.Links() {
			ent2 := entities[link.To().Value()]

			n2, err := g.NewNode(context.TODO(), ent2)
			if err != nil {
				return nil, err
			}

			if err := g.AddNode(context.TODO(), n2); err != nil {
				return nil, err
			}

			a, err := attrs.New()
			if err != nil {
				return nil, err
			}

			if relation := link.Attrs().Get("relation"); relation != "" {
				a.Set("relation", relation)
			}

			if _, err = g.Link(context.TODO(), n.UID(), n2.UID(), graph.WithAttrs(a)); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}
