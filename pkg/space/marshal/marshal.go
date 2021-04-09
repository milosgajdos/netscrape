package marshal

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

// EntityToSpace creates a new space.Entity from e and returns it.
func EntityToSpace(e Entity) (space.Entity, error) {
	uid := memuid.NewFromString(e.UID)
	a := memattrs.NewFromMap(e.Attrs)

	opts := []entity.Option{
		entity.WithUID(uid),
		entity.WithAttrs(a),
	}

	return entity.New(e.Type, opts...)
}

// ResourceToSpace creates a new space.Resource from Resources and returns it.
func ResourceToSpace(r Resource) (space.Resource, error) {
	uid := memuid.NewFromString(r.UID)
	a := memattrs.NewFromMap(r.Attrs)

	opts := []entity.Option{
		entity.WithUID(uid),
		entity.WithAttrs(a),
	}

	return entity.NewResource(r.Type, r.Name, r.Group, r.Version, r.Kind, r.Namespaced, opts...)
}

// ObjectToSpace creates a new space.Object from Entity and returns it.
func ObjectToSpace(o Object) (space.Object, error) {
	var r space.Resource
	if o.Resource != nil {
		var err error
		r, err = ResourceToSpace(*o.Resource)
		if err != nil {
			return nil, err
		}
	}

	uid := memuid.NewFromString(o.UID)
	a := memattrs.NewFromMap(o.Attrs)

	opts := []entity.Option{
		entity.WithUID(uid),
		entity.WithAttrs(a),
	}

	return entity.NewObject(o.Type, o.Name, o.Namespace, r, opts...)

}

// LinkToSpace creates a new space.Link from Link and returns it.
func LinkToSpace(l Link) (space.Link, error) {
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

// EntityFromSpace creates new Entity from e and returns it.
func EntityFromSpace(e space.Entity) (*Entity, error) {
	a, err := attrs.ToMap(context.Background(), e.Attrs())
	if err != nil {
		return nil, err
	}

	return &Entity{
		UID:   e.UID().String(),
		Type:  e.Type(),
		Attrs: a,
	}, nil
}

// ResourceFromSpace creates a new Resource from space.Resource and returns it.
func ResourceFromSpace(r space.Resource) (*Resource, error) {
	a, err := attrs.ToMap(context.Background(), r.Attrs())
	if err != nil {
		return nil, err
	}

	return &Resource{
		Entity: Entity{
			UID:   r.UID().String(),
			Type:  r.Type(),
			Attrs: a,
		},
		Name:       r.Name(),
		Group:      r.Group(),
		Version:    r.Version(),
		Kind:       r.Kind(),
		Namespaced: r.Namespaced(),
	}, nil
}

// ObjectFromSpace creates a new Object from space.Object and returns it.
func ObjectFromSpace(e space.Object) (*Object, error) {
	a, err := attrs.ToMap(context.Background(), e.Attrs())
	if err != nil {
		return nil, err
	}

	ent := &Object{
		Entity: Entity{
			UID:   e.UID().String(),
			Type:  e.Type(),
			Attrs: a,
		},
		Name:      e.Name(),
		Namespace: e.Namespace(),
	}

	if e.Resource() != nil {
		r, err := ResourceFromSpace(e.Resource())
		if err != nil {
			return nil, err
		}
		ent.Resource = r
	}

	return ent, nil
}

// LinkFromSpace creates a new Link from space.Link and returns it.
func LinkFromSpace(l space.Link) (*Link, error) {
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
