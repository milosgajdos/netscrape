package object

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Object is a space object.
type Object struct {
	uid   uuid.UID
	name  string
	ns    string
	dotid string
	res   space.Resource
	attrs attrs.Attrs
}

// New creates a new object and returns it.
func New(name, ns string, res space.Resource, opts ...Option) (*Object, error) {
	if res == nil {
		return nil, ErrMissingResource
	}

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

	dotid := strings.Join([]string{
		res.Group(),
		res.Version(),
		res.Kind(),
		ns,
		name}, "/")

	return &Object{
		uid:   uid,
		name:  name,
		ns:    ns,
		res:   res,
		dotid: dotid,
		attrs: a,
	}, nil
}

// UID returns UID.
func (o Object) UID() uuid.UID {
	return o.uid
}

// Type returns entity type
func (o Object) Type() entity.Type {
	return entity.ObjectType
}

// Name returns human readable object name.
func (o Object) Name() string {
	return o.name
}

// Namespace returns object namespace.
func (o Object) Namespace() string {
	return o.ns
}

// Resource returns resource the object is an instance of.
func (o Object) Resource() space.Resource {
	return o.res
}

// Attrs returns attributes.
func (o *Object) Attrs() attrs.Attrs {
	return o.attrs
}

// DOTID returns DOTID string
func (o Object) DOTID() string {
	return o.dotid
}

// SetDOTID sets DOTID
func (o *Object) SetDOTID(dotid string) {
	o.dotid = dotid
}
