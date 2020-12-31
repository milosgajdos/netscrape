package resource

import "github.com/milosgajdos/netscrape/pkg/metadata"

// Options are Space options.
type Options struct {
	// Metadata options
	Metadata metadata.Metadata
}

// Option configures Options.
type Option func(*Options)

// Metadata set metadata option
func Metadata(m metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}
