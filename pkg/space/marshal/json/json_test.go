package json

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/internal"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/space/marshal"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

const (
	entPath  = "testdata/entity.json"
	objPath  = "testdata/object.json"
	resPath  = "testdata/resource.json"
	linkPath = "testdata/link.json"
)

func MustReadFile(t *testing.T, path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return content
}

func MustResource(t *testing.T, opts ...entity.Option) space.Resource {
	r, err := internal.NewTestResource(opts...)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}
	return r
}

func MustObject(t *testing.T, opts ...entity.Option) space.Entity {
	e, err := internal.NewTestObject(opts...)
	if err != nil {
		t.Fatalf("failed to create object: %v", err)
	}
	return e
}

func MustEntity(t *testing.T, opts ...entity.Option) space.Entity {
	e, err := internal.NewTestEntity(opts...)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	return e
}

func MustLink(t *testing.T, uid1, uid2 string, opts ...link.Option) space.Link {
	l, err := internal.NewTestLink(opts...)
	if err != nil {
		t.Fatalf("failed to create link: %v", err)
	}
	return l
}

func MustAttrs(t *testing.T) attrs.Attrs {
	m := map[string]string{
		"foo": "bar",
	}
	return memattrs.NewFromMap(m)
}

func MustMarshaler(t *testing.T) *Marshaler {
	m, err := NewMarshaler()
	if err != nil {
		t.Fatalf("failed to create a new JSON Marshaler: %v", err)
	}
	return m
}

func TestNewMarshaler(t *testing.T) {
	if _, err := NewMarshaler(); err != nil {
		t.Fatalf("failed to create a new JSON Marshaler: %v", err)
	}
}

func TestMarshal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("ErrUnsuportedType", func(t *testing.T) {
		m := MustMarshaler(t)
		e := struct {
			Foo int
		}{Foo: 1}

		if _, err := m.Marshal(e); !errors.Is(err, marshal.ErrUnsuportedType) {
			t.Fatalf("expected error: %v, got: %v", marshal.ErrUnsuportedType, err)
		}
	})

	t.Run("Resource", func(t *testing.T) {
		m := MustMarshaler(t)
		a := MustAttrs(t)
		opts := []entity.Option{
			entity.WithUID(memuid.NewFromString(internal.ResUID)),
			entity.WithAttrs(a),
		}

		r := MustResource(t, opts...)
		if _, err := m.Marshal(r); err != nil {
			t.Fatalf("failed to marshal resource: %v", err)
		}
	})

	t.Run("Object", func(t *testing.T) {
		m := MustMarshaler(t)
		a := MustAttrs(t)
		opts := []entity.Option{
			entity.WithUID(memuid.NewFromString(internal.ObjUID)),
			entity.WithAttrs(a),
		}

		o := MustObject(t, opts...)
		if _, err := m.Marshal(o); err != nil {
			t.Fatalf("failed to marshal entity: %v", err)
		}
	})

	t.Run("Entity", func(t *testing.T) {
		m := MustMarshaler(t)
		a := MustAttrs(t)
		opts := []entity.Option{
			entity.WithUID(memuid.NewFromString(internal.EntUID)),
			entity.WithAttrs(a),
		}

		e := MustEntity(t, opts...)
		if _, err := m.Marshal(e); err != nil {
			t.Fatalf("failed to marshal entity: %v", err)
		}
	})

	t.Run("Link", func(t *testing.T) {
		m := MustMarshaler(t)
		a := MustAttrs(t)
		opts := []link.Option{
			link.WithUID(memuid.NewFromString(internal.LinkUID)),
			link.WithAttrs(a),
		}

		l := MustLink(t, "foo1", "foo2", opts...)
		if _, err := m.Marshal(l); err != nil {
			t.Fatalf("failed to marshal link: %v", err)
		}
	})
}

func TestUnmarshal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("ErrUnsuportedType", func(t *testing.T) {
		m := MustMarshaler(t)
		e := struct {
			Foo int
		}{Foo: 1}

		if err := m.Unmarshal([]byte{}, e); !errors.Is(err, marshal.ErrUnsuportedType) {
			t.Fatalf("expected error: %v, got: %v", marshal.ErrUnsuportedType, err)
		}
	})

	t.Run("Resource", func(t *testing.T) {
		m := MustMarshaler(t)
		b := MustReadFile(t, resPath)

		var r space.Resource
		if err := m.Unmarshal(b, &r); err != nil {
			t.Fatalf("failed to unmarshal data to resource: %v", err)
		}
	})

	t.Run("Object", func(t *testing.T) {
		m := MustMarshaler(t)
		b := MustReadFile(t, objPath)

		var o space.Object
		if err := m.Unmarshal(b, &o); err != nil {
			t.Fatalf("failed to unmarshal data to object: %v", err)
		}
	})

	t.Run("Entity", func(t *testing.T) {
		m := MustMarshaler(t)
		b := MustReadFile(t, entPath)

		var e space.Entity
		if err := m.Unmarshal(b, &e); err != nil {
			t.Fatalf("failed to unmarshal data to entity: %v", err)
		}
	})

	t.Run("Link", func(t *testing.T) {
		m := MustMarshaler(t)
		b := MustReadFile(t, linkPath)

		var l space.Link
		if err := m.Unmarshal(b, &l); err != nil {
			t.Fatalf("failed to unmarshal data to link: %v", err)
		}
	})
}
