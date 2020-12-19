package entity

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
)

// Options are store options.
type Options struct {
	Attrs attrs.Attrs
}

// Option sets options.
type Option func(*Options)

// Attrs sets entity attributes.
func Attrs(a attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}
