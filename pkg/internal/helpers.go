package internal

import (
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

func NewTestResource(typ, name, group, version, kind string, namespaced bool, opts ...resource.Option) (space.Resource, error) {
	return resource.New(typ, name, group, version, kind, namespaced, opts...)
}

func NewTestEntity(uid, typ, name, ns string, res space.Resource, opts ...entity.Option) (space.Entity, error) {
	u, err := uuid.NewFromString(uid)
	if err != nil {
		return nil, err
	}

	opts = append(opts, entity.WithUID(u))

	return entity.New(typ, name, ns, res, opts...)
}
