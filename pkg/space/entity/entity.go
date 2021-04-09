package entity

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

const (
	Partial = "_partial"
)

// Entity is a generic space entity.
type Entity struct {
	uid   uuid.UID
	typ   string
	dotid string
	attrs attrs.Attrs
}

// New creates a new entity and returns it
func New(typ string, opts ...Option) (*Entity, error) {
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
	}

	return &Entity{
		uid:   uid,
		typ:   typ,
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
