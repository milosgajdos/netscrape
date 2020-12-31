package memory

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/space"
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

func newTestObject(uid, name, ns string, res space.Resource, opts ...object.Option) (space.Object, error) {
	u, err := uuid.NewFromString(uid)
	if err != nil {
		return nil, err
	}

	return object.New(u, name, ns, res, opts...)
}

func makeTestSpaceObjects(path string) (map[string]space.Object, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var testObjects []types.Object
	if err := yaml.Unmarshal(data, &testObjects); err != nil {
		return nil, err
	}

	objects := make(map[string]space.Object)

	for _, o := range testObjects {
		m, err := metadata.NewFromMap(o.Resource.Metadata)
		if err != nil {
			return nil, err
		}

		res, err := resource.New(
			o.Resource.Name,
			o.Resource.Group,
			o.Resource.Version,
			o.Resource.Kind,
			o.Resource.Namespaced,
			resource.Metadata(m),
		)
		if err != nil {
			return nil, err
		}

		m, err = metadata.NewFromMap(o.Metadata)
		if err != nil {
			return nil, err
		}

		uid, err := uuid.NewFromString(o.UID)
		if err != nil {
			return nil, err
		}

		obj, err := object.New(uid, o.Name, o.Namespace, res, object.Metadata(m))
		if err != nil {
			return nil, err
		}

		for _, l := range o.Links {
			toUID, err := uuid.NewFromString(l.To)
			if err != nil {
				return nil, err
			}

			m, err = metadata.NewFromMap(l.Metadata)
			if err != nil {
				return nil, err
			}

			if err := obj.Link(toUID, space.LinkOptions{Metadata: m}); err != nil {
				return nil, err
			}
		}

		objects[o.UID] = obj
	}

	return objects, nil
}

func makeTestGraph(path string) (*WUG, error) {
	g, err := NewWUG("test", graph.Options{})
	if err != nil {
		return nil, err
	}

	objects, err := makeTestSpaceObjects(path)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		n, err := g.NewNode(object, graph.NodeOptions{})
		if err != nil {
			return nil, err
		}

		if err := g.AddNode(n); err != nil {
			return nil, err
		}

		for _, link := range object.Links() {
			object2 := objects[link.To().Value()]

			n2, err := g.NewNode(object2, graph.NodeOptions{})
			if err != nil {
				return nil, err
			}

			if err := g.AddNode(n2); err != nil {
				return nil, err
			}

			a, err := attrs.New()
			if err != nil {
				return nil, err
			}

			if relation, ok := link.Metadata().Get("relation").(string); ok {
				a.Set("relation", relation)
			}

			if _, err = g.Link(n.UID(), n2.UID(), graph.LinkOptions{Attrs: a}); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}
