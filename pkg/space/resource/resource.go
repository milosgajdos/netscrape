package resource

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Resource implements a generic Space resource.
type Resource struct {
	uid        uuid.UID
	name       string
	group      string
	version    string
	kind       string
	namespaced bool
	attrs      attrs.Attrs
}

// New creates a new generic resource and returns it.
func New(name, group, version, kind string, namespaced bool, opts ...Option) (*Resource, error) {
	ropts := Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	uid := ropts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	a := ropts.Attrs
	if a == nil {
		var err error
		a, err = attrs.New()
		if err != nil {
			return nil, err
		}
	}

	return &Resource{
		uid:        uid,
		name:       name,
		group:      group,
		version:    version,
		kind:       kind,
		namespaced: namespaced,
		attrs:      a,
	}, nil
}

// UID returns UID.
func (r Resource) UID() uuid.UID {
	return r.uid
}

// Name returns resource name.
func (r Resource) Name() string {
	return r.name
}

// Group returns resource group.
func (r Resource) Group() string {
	return r.group
}

// Version returns resource version.
func (r Resource) Version() string {
	return r.version
}

// Kind returns resource kind.
func (r Resource) Kind() string {
	return r.kind
}

// Namespaced returns true if the resource is namespaced.
func (r Resource) Namespaced() bool {
	return r.namespaced
}

// Paths returns all possible variations resource paths.
func (r Resource) Paths() []string {
	resNames := []string{strings.ToLower(r.name)}

	// nolint:prealloc
	var names []string
	for _, name := range resNames {
		names = append(names,
			name,
			strings.Join([]string{name, r.group}, "/"),
			strings.Join([]string{name, r.group, r.version}, "/"),
		)
	}

	return names
}

// Attrs returns attributes.
func (r Resource) Attrs() attrs.Attrs {
	return r.attrs
}

// DOTID returns DOTID string
func (r Resource) DOTID() string {
	return strings.Join([]string{
		r.Group(),
		r.Version(),
		r.Kind()}, "/")
}
