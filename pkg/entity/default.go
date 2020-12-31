package entity

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// entity implements Entity.
type entity struct {
	uid   uuid.UID
	attrs attrs.Attrs
}

// NewWithUID creates a new entity with given uid and returns it.
func NewWithUID(uid uuid.UID, opts ...Option) (*entity, error) {
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	eopts := Options{}
	for _, apply := range opts {
		apply(&eopts)
	}

	return &entity{
		uid:   uid,
		attrs: eopts.Attrs,
	}, nil
}

// New creates a new entity and returns it.
func New(opts ...Option) (*entity, error) {
	u, err := uuid.New()
	if err != nil {
		return nil, err
	}

	return NewWithUID(u, opts...)
}

// UID returns entity UID.
func (e entity) UID() uuid.UID {
	return e.uid
}

// Attrs returns entity attributes.
func (e *entity) Attrs() attrs.Attrs {
	return e.attrs
}
