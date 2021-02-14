package top

import (
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	resName    = "nodeResName"
	resGroup   = "nodeResGroup"
	resVersion = "nodeResVersion"
	resKind    = "nodeResKind"
	entUID     = "testID"
	entName    = "testName"
	entNs      = "testNs"
	ent2UID    = "test2ID"
	ent2Name   = "test2Name"
	ent2Ns     = "test2Ns"
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