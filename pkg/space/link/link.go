package link

import (
	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Link links two unique space objects.
type Link struct {
	uid  uuid.UID
	from uuid.UID
	to   uuid.UID
	md   metadata.Metadata
}

// New creates a new link between two objects and returns it.
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

	md := lopts.Metadata
	if md == nil {
		var err error
		md, err = metadata.New()
		if err != nil {
			return nil, err
		}
	}

	return &Link{
		uid:  uid,
		from: from,
		to:   to,
		md:   md,
	}, nil
}

// UID returns link uid.
func (l Link) UID() uuid.UID {
	return l.uid
}

// From returns uid of Link origin.
func (l Link) From() uuid.UID {
	return l.from
}

// To returns uid of Link end.
func (l Link) To() uuid.UID {
	return l.to
}

// Metadata returns Link metadata.
func (l Link) Metadata() metadata.Metadata {
	return l.md
}
