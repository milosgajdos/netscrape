package top

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/object"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/space/types"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// NewMock creates mock Top from objects and resrouces
// defined in given path and returns it.
func NewMock(a space.Plan, path string) (*Top, error) {
	t, err := New(a)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var objects []types.Object
	if err := yaml.Unmarshal(data, &objects); err != nil {
		return nil, err
	}

	for _, o := range objects {
		r, err := resource.New(
			o.Resource.Name,
			o.Resource.Kind,
			o.Resource.Group,
			o.Resource.Version,
			o.Resource.Namespaced,
		)
		if err != nil {
			return nil, err
		}

		md, err := metadata.New()
		if err != nil {
			return nil, err
		}

		uid, err := uuid.NewFromString(o.UID)
		if err != nil {
			return nil, err
		}

		obj, err := object.New(
			uid,
			o.Name,
			o.Namespace,
			r,
			object.Metadata(md),
		)

		if err != nil {
			return nil, err
		}

		for _, l := range o.Links {
			md, err := metadata.NewFromMap(l.Metadata)
			if err != nil {
				return nil, err
			}

			lUID, err := uuid.NewFromString(l.UID)
			if err != nil {
				return nil, err
			}

			toUID, err := uuid.NewFromString(l.To)
			if err != nil {
				return nil, err
			}

			opts := space.LinkOptions{UID: lUID, Metadata: md}

			if err := obj.Link(toUID, opts); err != nil {
				return nil, err
			}
		}

		if err := t.Add(obj, space.AddOptions{}); err != nil {
			return nil, err
		}
	}

	return t, nil
}
