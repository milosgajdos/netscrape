package entity

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is an arbitrary entity.
type Entity interface {
	// UID returns unique ID.
	UID() uuid.UID
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}
