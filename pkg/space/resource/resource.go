package resource

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/metadata"
)

// Resource implements a generic Space resource.
type Resource struct {
	name       string
	group      string
	version    string
	kind       string
	namespaced bool
	md         metadata.Metadata
}

// New creates a new generic resource and returns it.
func New(name, group, version, kind string, namespaced bool, opts ...Option) (*Resource, error) {
	ropts := Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	md := ropts.Metadata
	if md == nil {
		var err error
		md, err = metadata.New()
		if err != nil {
			return nil, err
		}
	}

	return &Resource{
		name:       name,
		group:      group,
		version:    version,
		kind:       kind,
		namespaced: namespaced,
		md:         md,
	}, nil
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

// Namespaced returns true if the resource objects are namespaced.
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

// Metadata returns resource metadata.
func (r Resource) Metadata() metadata.Metadata {
	return r.md
}
