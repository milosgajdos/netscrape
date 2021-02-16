package top

import (
	"context"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/object"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/space/types"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// NewMock creates mock Top from entities and resrouces
// defined in given path and returns it.
func NewMock(path string) (*Top, error) {
	t, err := New()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entities []types.Object
	if err := yaml.Unmarshal(data, &entities); err != nil {
		return nil, err
	}

	ctx := context.Background()

	for _, o := range entities {
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

		a, err := attrs.New()
		if err != nil {
			return nil, err
		}

		uid, err := uuid.NewFromString(o.UID)
		if err != nil {
			return nil, err
		}

		ent, err := object.New(
			o.Name,
			o.Namespace,
			r,
			object.WithUID(uid),
			object.WithAttrs(a),
		)

		if err != nil {
			return nil, err
		}

		if err := t.Add(ctx, ent); err != nil {
			return nil, err
		}

		for _, l := range o.Links {
			a, err := attrs.NewFromMap(l.Attrs)
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

			if err := t.Link(ctx, ent.UID(), toUID, space.WithUID(lUID), space.WithAttrs(a)); err != nil {
				return nil, err
			}
		}
	}

	return t, nil
}
