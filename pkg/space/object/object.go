package object

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Object is a space object
type Object struct {
	uid  uuid.UID
	name string
	ns   string
	res  space.Resource
	// links indexes all links to this object
	// for faster link lookups
	links map[string]space.Link
	// olinks indexes links from this object to
	// other object for faster object lookups.
	olinks map[string]space.Link
	attrs  attrs.Attrs
}

// New creates a new Object and returns it.
func New(name, ns string, res space.Resource, opts ...Option) (*Object, error) {
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

	return &Object{
		uid:    uid,
		name:   name,
		ns:     ns,
		res:    res,
		links:  make(map[string]space.Link),
		olinks: make(map[string]space.Link),
		attrs:  a,
	}, nil
}

// UID returns object uid.
func (o Object) UID() uuid.UID {
	return o.uid
}

// Name returns object name.
func (o Object) Name() string {
	return o.name
}

// Namespace returns object namespace.
func (o Object) Namespace() string {
	return o.ns
}

// Resource returns the resource the object is an instance of.
func (o Object) Resource() space.Resource {
	return o.res
}

// link creates a new link to object to.
func (o *Object) link(u uuid.UID, opts ...space.Option) error {
	lopts := space.Options{}
	for _, apply := range opts {
		apply(&lopts)
	}

	nopts := []link.Option{
		link.WithUID(lopts.UID),
		link.WithAttrs(lopts.Attrs),
		link.WithMerge(lopts.Merge),
	}

	link, err := link.New(o.uid, u, nopts...)
	if err != nil {
		return err
	}

	if _, ok := o.links[link.UID().Value()]; !ok {
		o.links[link.UID().Value()] = link
	}

	o.olinks[u.Value()] = link

	return nil
}

// Link links object to object to with the given UID.
// If link merging is requested, the new link will contain
// all the attributes of the existing link with addition to the attributes
// that are not in the original link. The original attributes are updated.
func (o *Object) Link(to uuid.UID, opts ...space.Option) error {
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
func (o Object) Links() []space.Link {
	links := make([]space.Link, len(o.links))

	i := 0
	for _, link := range o.links {
		links[i] = link
		i++
	}

	return links
}

// Attrs returns attributes.
func (o *Object) Attrs() attrs.Attrs {
	return o.attrs
}
