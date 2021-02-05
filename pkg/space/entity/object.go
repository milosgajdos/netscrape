package entity

import (
	"sync"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Entity is a space entity.
type Entity struct {
	uid  uuid.UID
	name string
	ns   string
	res  space.Resource
	// links indexes all links to this object
	// for faster link lookups
	// NOTE: index key is the UID of the *link*
	links map[string]space.Link
	// olinks indexes links from this object to
	// other objects for faster lookups.
	// NOTE: index key is the UID of the *object* on the opposite end
	olinks map[string]space.Link
	attrs  attrs.Attrs
	// mu synchronizes access to Top
	mu *sync.RWMutex
}

// New creates a new Entity and returns it.
func New(name, ns string, res space.Resource, opts ...Option) (*Entity, error) {
	oopts := Options{}
	for _, apply := range opts {
		apply(&oopts)
	}

	uid := oopts.UID
	if uid == nil {
		var err error
		uid, err = uuid.New()
		if err != nil {
			return nil, err
		}
	}

	a := oopts.Attrs
	if a == nil {
		var err error
		a, err = attrs.New()
		if err != nil {
			return nil, err
		}
	}

	return &Entity{
		uid:    uid,
		name:   name,
		ns:     ns,
		res:    res,
		links:  make(map[string]space.Link),
		olinks: make(map[string]space.Link),
		attrs:  a,
		mu:     &sync.RWMutex{},
	}, nil
}

// UID returns object uid.
func (o Entity) UID() uuid.UID {
	return o.uid
}

// Name returns object name.
func (o Entity) Name() string {
	return o.name
}

// Namespace returns object namespace.
func (o Entity) Namespace() string {
	return o.ns
}

// Resource returns the resource the object is an instance of.
func (o Entity) Resource() space.Resource {
	return o.res
}

// link creates a new link to object to.
func (o *Entity) link(to uuid.UID, opts ...space.Option) error {
	lopts := space.Options{}
	for _, apply := range opts {
		apply(&lopts)
	}

	nopts := []link.Option{
		link.WithUID(lopts.UID),
		link.WithAttrs(lopts.Attrs),
		link.WithMerge(lopts.Merge),
	}

	link, err := link.New(o.uid, to, nopts...)
	if err != nil {
		return err
	}

	if _, ok := o.links[link.UID().Value()]; !ok {
		o.links[link.UID().Value()] = link
	}

	o.olinks[to.Value()] = link

	return nil
}

// Link links object to object to with the given UID.
// If WithMergeAttrs option is set to true, existing link attributes
// linking to the same object are merged in with the passed in attributes.
// NOTE: the original attributes can be overridden in place.
func (o *Entity) Link(to uuid.UID, opts ...space.Option) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	lopts := space.Options{}
	for _, apply := range opts {
		apply(&lopts)
	}

	l, ok := o.olinks[to.Value()]
	if !ok {
		return o.link(to, opts...)
	}

	if lopts.Merge {
		if lopts.Attrs != nil {
			for _, k := range lopts.Attrs.Keys() {
				l.Attrs().Set(k, lopts.Attrs.Get(k))
			}
		}
	}

	return nil
}

// Links returns a slice of all object links.
func (o Entity) Links() []space.Link {
	o.mu.Lock()
	defer o.mu.Unlock()

	links := make([]space.Link, len(o.links))

	i := 0
	for _, link := range o.links {
		links[i] = link
		i++
	}

	return links
}

// Attrs returns attributes.
// NOTE: Attrs is not thread-safe
func (o *Entity) Attrs() attrs.Attrs {
	return o.attrs
}
