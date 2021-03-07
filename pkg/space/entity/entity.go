package entity

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is a space entity.
type Entity struct {
	uid   uuid.UID
	name  string
	ns    string
	dotid string
	res   space.Resource
	attrs attrs.Attrs
}

// New creates a new entity and returns it.
func New(name, ns string, res space.Resource, opts ...Option) (*Entity, error) {
	eopts := Options{}
	for _, apply := range opts {
		apply(&eopts)
	}

	uid := eopts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	a := eopts.Attrs
	if a == nil {
		var err error
		a, err = attrs.New()
		if err != nil {
			return nil, err
		}
	}

	dotid := eopts.DOTID
	if dotid == "" {
		dotid = uid.Value()
		if res != nil {
			dotid = strings.Join([]string{
				res.Group(),
				res.Version(),
				res.Kind(),
				ns,
				name}, "/")
		}
	}

	return &Entity{
		uid:   uid,
		name:  name,
		ns:    ns,
		res:   res,
		dotid: dotid,
		attrs: a,
	}, nil
}

// UID returns UID.
func (e Entity) UID() uuid.UID {
	return e.uid
}

// Type returns entity type.
func (e Entity) Type() entity.Type {
	return entity.EntityType
}

// Name returns human readable entity name.
func (e Entity) Name() string {
	return e.name
}

// Namespace returns entity namespace.
func (e Entity) Namespace() string {
	return e.ns
}

// Resource returns entity resource.
func (e Entity) Resource() space.Resource {
	return e.res
}

// Attrs returns attributes.
func (e *Entity) Attrs() attrs.Attrs {
	return e.attrs
}

// DOTID returns DOTID string.
func (e Entity) DOTID() string {
	return e.dotid
}

// SetDOTID sets DOTID.
func (e *Entity) SetDOTID(dotid string) {
	e.dotid = dotid
}
