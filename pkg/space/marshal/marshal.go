package marshal

import (
	"context"
	"encoding/json"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/space/resource"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

// Format defines encoding data format.
type Format int

const (
	// JSON data encoding.
	JSON Format = iota
)

// Marshal marshals m into format f.
func Marshal(f Format, m interface{}) ([]byte, error) {
	switch f {
	case JSON:
		return marshalJSON(m)
	default:
		return nil, ErrUnsupportedFormat
	}
}

// marshalJSON marshals m to JSON.
func marshalJSON(m interface{}) ([]byte, error) {
	switch v := m.(type) {
	case space.Entity:
		e, err := FromSpaceEntity(v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(e)
	case space.Resource:
		r, err := FromSpaceResource(v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(r)
	case space.Link:
		l, err := FromSpaceLink(v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(l)
	default:
		return nil, ErrUnsuportedType
	}
}

// Unmarshal unmarshals b in format f into m.
func Unmarshal(f Format, b []byte, m interface{}) error {
	switch f {
	case JSON:
		return unmarshalJSON(b, m)
	default:
		return ErrUnsupportedFormat
	}
}

// unmarshalJSON decodes JSON data stored in b into m.
func unmarshalJSON(b []byte, m interface{}) error {
	switch m.(type) {
	case *Entity, *Resource, *Link, *LinkedEntity:
		return json.Unmarshal(b, m)
	default:
		return ErrUnsuportedType
	}
}

// ToSpaceLink creates a new space.Link from Link and returns it.
func ToSpaceLink(l Link) (space.Link, error) {
	uid := memuid.NewFromString(l.UID)
	a := memattrs.NewFromMap(l.Attrs)

	opts := []link.Option{
		link.WithUID(uid),
		link.WithAttrs(a),
	}

	from := memuid.NewFromString(l.From)
	to := memuid.NewFromString(l.To)

	return link.New(from, to, opts...)
}

// ToSpaceResource creates a new space.Resource from Resources and returns it.
func ToSpaceResource(r Resource) (space.Resource, error) {
	uid := memuid.NewFromString(r.UID)
	a := memattrs.NewFromMap(r.Attrs)

	opts := []resource.Option{
		resource.WithUID(uid),
		resource.WithAttrs(a),
	}

	return resource.New(r.Type, r.Name, r.Group, r.Version, r.Kind, r.Namespaced, opts...)
}

// ToSpaceEntity creates a new space.Entity from Entity and returns it.
func ToSpaceEntity(e Entity) (space.Entity, error) {
	var r space.Resource

	if e.Resource != nil {
		var err error

		r, err = ToSpaceResource(*e.Resource)
		if err != nil {
			return nil, err
		}
	}

	uid := memuid.NewFromString(e.UID)
	a := memattrs.NewFromMap(e.Attrs)

	opts := []entity.Option{
		entity.WithUID(uid),
		entity.WithAttrs(a),
	}

	return entity.New(e.Type, e.Name, e.Namespace, r, opts...)

}

// FromSpaceResource creates a new Resource from space.Resource and returns it.
func FromSpaceResource(r space.Resource) (*Resource, error) {
	a, err := attrs.ToMap(context.Background(), r.Attrs())
	if err != nil {
		return nil, err
	}

	return &Resource{
		UID:        r.UID().String(),
		Type:       r.Type(),
		Name:       r.Name(),
		Group:      r.Group(),
		Version:    r.Version(),
		Kind:       r.Kind(),
		Namespaced: r.Namespaced(),
		Attrs:      a,
	}, nil
}

// FromSpaceEntity creates a new Entity from space.Entity and returns it.
func FromSpaceEntity(e space.Entity) (*Entity, error) {
	a, err := attrs.ToMap(context.Background(), e.Attrs())
	if err != nil {
		return nil, err
	}

	ent := &Entity{
		UID:       e.UID().String(),
		Type:      e.Type(),
		Name:      e.Name(),
		Namespace: e.Namespace(),
		Attrs:     a,
	}

	if e.Resource() != nil {
		r, err := FromSpaceResource(e.Resource())
		if err != nil {
			return nil, err
		}
		ent.Resource = r
	}

	return ent, nil
}

// FromSpaceLink creates a new Link from space.Link and returns it.
func FromSpaceLink(l space.Link) (*Link, error) {
	a, err := attrs.ToMap(context.Background(), l.Attrs())
	if err != nil {
		return nil, err
	}

	return &Link{
		UID:   l.UID().String(),
		From:  l.From().String(),
		To:    l.To().String(),
		Attrs: a,
	}, nil
}
