package link

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Link links two space entities.
type Link struct {
	uid   uuid.UID
	from  uuid.UID
	to    uuid.UID
	attrs attrs.Attrs
}

// New creates a new link between two entities and returns it.
func New(from, to uuid.UID, opts ...Option) (*Link, error) {
	lopts := Options{}
	for _, apply := range opts {
		apply(&lopts)
	}

	uid := lopts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	a := lopts.Attrs
	if a == nil {
		var err error
		a, err = attrs.New()
		if err != nil {
			return nil, err
		}
	}

	return &Link{
		uid:   uid,
		from:  from,
		to:    to,
		attrs: a,
	}, nil
}

// UID returns link uid.
func (l Link) UID() uuid.UID {
	return l.uid
}

// From returns uid of link origin.
func (l Link) From() uuid.UID {
	return l.from
}

// To returns uid of link end.
func (l Link) To() uuid.UID {
	return l.to
}

// Attrs returns attributes.
func (l Link) Attrs() attrs.Attrs {
	return l.attrs
}
