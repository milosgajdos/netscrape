package space

import (
	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Options are Space options.
type Options struct {
	// Metadata options
	Metadata metadata.Metadata
}

// Option configures Options.
type Option func(*Options)

// AddOptions are Space add options
type AddOptions struct {
	// MergeLinks merges link with existing link.
	MergeLinks bool
}

// AddOption sets AddOptions.
type AddOption func(*AddOptions)

// LinkOptions are link options.
type LinkOptions struct {
	// UID is an optional link UID.
	UID uuid.UID
	// Merge merges link with existing link.
	Merge bool
	// Metadata options.
	Metadata metadata.Metadata
}

// LinkOption sets LinkOptions.
type LinkOption func(*LinkOptions)
