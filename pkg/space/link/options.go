package link

import (
	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Options are link options.
type Options struct {
	// UID is an optional link UID.
	UID uuid.UID
	// Merge merges link with existing link.
	Merge bool
	// Metadata options.
	Metadata metadata.Metadata
}

// Option sets LinkOptions.
type Option func(*Options)

// UID set UID option
func UID(u uuid.UID) Option {
	return func(o *Options) {
		o.UID = u
	}
}

// Merge set merge option
func Merge(m bool) Option {
	return func(o *Options) {
		o.Merge = m
	}
}

// Metadata set metadata option
func Metadata(m metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}
