package entity

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

// Object is a space object.
type Object struct {
	uid   uuid.UID
	typ   string
	name  string
	ns    string
	res   space.Resource
	dotid string
	attrs attrs.Attrs
}

// NewObject creates a new object and returns it.
func NewObject(typ, name, ns string, res space.Resource, opts ...Option) (*Object, error) {
	eopts := Options{}
	for _, apply := range opts {
		apply(&eopts)
	}

	uid := eopts.UID
	if uid == nil {
		uid = memuid.New()
	}

	a := eopts.Attrs
	if a == nil {
		a = memattrs.New()
	}

	dotid := eopts.DOTID
	if dotid == "" {
		dotid = uid.String()
		if res != nil {
			dotid = strings.Join([]string{
				res.Group(),
				res.Version(),
				res.Kind(),
				ns,
				name}, "/")
		}
	}

	return &Object{
		uid:   uid,
		typ:   typ,
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

// Type returns object type.
func (o Object) Type() string {
	return o.typ
}

// Attrs returns attributes.
func (e *Object) Attrs() attrs.Attrs {
	return e.attrs
}

// Name returns human readable object name.
func (o Object) Name() string {
	return o.name
}

// Namespace returns object namespace.
func (o Object) Namespace() string {
	return o.ns
}

// Resource returns object resource.
func (o Object) Resource() space.Resource {
	return o.res
}

// DOTID returns DOTID string.
func (o Object) DOTID() string {
	return o.dotid
}

// SetDOTID sets DOTID.
func (o *Object) SetDOTID(dotid string) {
	o.dotid = dotid
}
