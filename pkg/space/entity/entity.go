package entity

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is a space entity.
type Entity struct {
	uid   uuid.UID
	name  string
	ns    string
	res   space.Resource
	attrs attrs.Attrs
}

// New creates a new Entity and returns it.
func New(name, ns string, res space.Resource, opts ...Option) (*Entity, error) {
	oopts := Options{}
	for _, apply := range opts {
		apply(&oopts)
	}

	uid := oopts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	a := oopts.Attrs
	if a == nil {
		var err error
		a, err = attrs.New()
		if err != nil {
			return nil, err
		}
	}

	return &Entity{
		uid:   uid,
		name:  name,
		ns:    ns,
		res:   res,
		attrs: a,
	}, nil
}

// UID returns UID.
func (o Entity) UID() uuid.UID {
	return o.uid
}

// Name returns human readable Entity name.
func (o Entity) Name() string {
	return o.name
}

// Namespace returns entity namespace.
func (o Entity) Namespace() string {
	return o.ns
}

// Resource returns resource the entity is an instance of.
func (o Entity) Resource() space.Resource {
	return o.res
}

// Attrs returns attributes.
func (o *Entity) Attrs() attrs.Attrs {
	return o.attrs
}
