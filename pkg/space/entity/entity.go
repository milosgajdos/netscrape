package entity

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

const (
	PartialType = "_partial"
	PartialName = "_partialName"
	PartialNs   = "_partialNs"
)

// Entity is a space entity.
type Entity struct {
	uid   uuid.UID
	typ   string
	name  string
	ns    string
	res   space.Resource
	dotid string
	attrs attrs.Attrs
}

// NewPartial creates a new partial entity and returns it.
// NOTE: Partial entity has no Resource associated with it.
func NewPartial(opts ...Option) (*Entity, error) {
	return New(PartialType, PartialName, PartialNs, nil, opts...)
}

// New creates a new entity and returns it.
func New(typ, name, ns string, res space.Resource, opts ...Option) (*Entity, error) {
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

	return &Entity{
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
func (e Entity) UID() uuid.UID {
	return e.uid
}

// Type returns entity type.
func (e Entity) Type() string {
	return e.typ
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
