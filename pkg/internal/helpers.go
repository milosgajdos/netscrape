package internal

import (
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"

	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

const (
	EntUID     = "testEntUID"
	EntType    = "testEntType"
	ResUID     = "testResUID"
	ResType    = "testResType"
	ResName    = "testResName"
	ResGroup   = "testResGroup"
	ResVersion = "testResVersion"
	ResKind    = "testResKind"
	ResNsd     = false
	ObjUID     = "testObjUID"
	ObjType    = "testObjType"
	ObjName    = "testObjName"
	ObjNs      = "testObjNs"
	LinkUID    = "testLinkUID"
	LinkFrom   = "testFromUID"
	LinkTo     = "testToUID"
)

func NewTestEntity(opts ...entity.Option) (space.Entity, error) {
	return entity.New(EntType, opts...)
}

func NewTestResource(opts ...entity.Option) (space.Resource, error) {
	return entity.NewResource(ResType, ResName, ResGroup, ResVersion, ResKind, ResNsd, opts...)
}

func NewNamedTestObject(name string, opts ...entity.Option) (space.Object, error) {
	r, err := entity.NewResource(ResType, ResName, ResGroup, ResVersion, ResKind, ResNsd)
	if err != nil {
		return nil, err
	}
	return entity.NewObject(ObjType, name, ObjNs, r, opts...)
}

func NewTestObject(opts ...entity.Option) (space.Object, error) {
	r, err := entity.NewResource(ResType, ResName, ResGroup, ResVersion, ResKind, ResNsd)
	if err != nil {
		return nil, err
	}
	return entity.NewObject(ObjType, ObjName, ObjNs, r, opts...)
}

func NewTestLink(opts ...link.Option) (space.Link, error) {
	from := memuid.NewFromString(LinkFrom)
	to := memuid.NewFromString(LinkTo)
	return link.New(from, to, opts...)
}
