package marshal

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/internal"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	resName    = "nodeResName"
	resType    = "nodeResType"
	resGroup   = "nodeResGroup"
	resVersion = "nodeResVersion"
	resKind    = "nodeResKind"
	entUID     = "testID"
	entType    = "testType"
	entName    = "testName"
	entNs      = "testNs"
	entDOTID   = "dotID"

	// NOTE: we could list all files in testadata dir programmatically, yeah
	entJsonPath  = "testdata/entity.json"
	resJsonPath  = "testdata/resource.json"
	linkJsonPath = "testdata/link.json"
)

func MustReadFile(t *testing.T, path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return content
}

func MustResource(t *testing.T, opts ...resource.Option) space.Resource {
	r, err := internal.NewTestResource(resType, resName, resGroup, resVersion, resKind, false, opts...)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}
	return r
}

func MustEntity(t *testing.T, opts ...entity.Option) space.Entity {
	r := MustResource(t)

	e, err := internal.NewTestEntity(entUID, entType, entName, entNs, r, opts...)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	return e
}

func MustEntityNoResource(t *testing.T, opts ...entity.Option) space.Entity {
	e, err := internal.NewTestEntity(entUID, entType, entName, entNs, nil, opts...)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	return e
}

func MustLink(t *testing.T, uid1, uid2 string, opts ...link.Option) space.Link {
	uuid1, err := uuid.NewFromString(uid1)
	if err != nil {
		t.Fatalf("failed to create uid from %s: %v", uid1, err)
	}

	uuid2, err := uuid.NewFromString(uid2)
	if err != nil {
		t.Fatalf("failed to create uid from %s: %v", uid1, err)
	}

	l, err := link.New(uuid1, uuid2, opts...)
	if err != nil {
		t.Fatalf("failed to create link: %v", err)
	}
	return l
}

func TestMarshal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("ErrUnsupportedFormat", func(t *testing.T) {
		e := MustEntity(t)

		if _, err := Marshal(Format(-1000), e); !errors.Is(err, ErrUnsupportedFormat) {
			t.Fatalf("expected error: %v, got: %v", ErrUnsupportedFormat, err)
		}
	})

	t.Run("ErrUnsuportedType", func(t *testing.T) {
		e := struct {
			Foo int
		}{Foo: 1}

		if _, err := Marshal(JSON, e); !errors.Is(err, ErrUnsuportedType) {
			t.Fatalf("expected error: %v, got: %v", ErrUnsuportedType, err)
		}
	})

	t.Run("JSONEntity", func(t *testing.T) {
		e := MustEntity(t)
		// TODO: add a check for expected JSON output
		if _, err := Marshal(JSON, e); err != nil {
			t.Fatalf("failed to marshal entity: %v", err)
		}
	})

	t.Run("JSONEntityNoResource", func(t *testing.T) {
		e := MustEntityNoResource(t)
		// TODO: add a check for expected JSON output
		if _, err := Marshal(JSON, e); err != nil {
			t.Fatalf("failed to marshal entity: %v", err)
		}
	})

	t.Run("JSONResource", func(t *testing.T) {
		r := MustResource(t)
		// TODO: add a check for expected JSON output
		if _, err := Marshal(JSON, r); err != nil {
			t.Fatalf("failed to marshal resource: %v", err)
		}
	})

	t.Run("JSONLink", func(t *testing.T) {
		l := MustLink(t, "foo1", "foo2")
		// TODO: add a check for expected JSON output
		if _, err := Marshal(JSON, l); err != nil {
			t.Fatalf("failed to marshal link: %v", err)
		}
	})
}

func TestUnmarshal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("ErrUnsupportedFormat", func(t *testing.T) {
		e := MustEntity(t)

		if err := Unmarshal(Format(-1000), []byte{}, e); !errors.Is(err, ErrUnsupportedFormat) {
			t.Fatalf("expected error: %v, got: %v", ErrUnsupportedFormat, err)
		}
	})

	t.Run("ErrUnsuportedType", func(t *testing.T) {
		e := struct {
			Foo int
		}{Foo: 1}

		if err := Unmarshal(JSON, []byte{}, e); !errors.Is(err, ErrUnsuportedType) {
			t.Fatalf("expected error: %v, got: %v", ErrUnsuportedType, err)
		}
	})

	t.Run("JSONEntity", func(t *testing.T) {
		b := MustReadFile(t, entJsonPath)

		var e Entity
		if err := Unmarshal(JSON, b, &e); err != nil {
			t.Fatalf("failed to unmarshal data to entity: %v", err)
		}
	})

	t.Run("JSONResource", func(t *testing.T) {
		b := MustReadFile(t, resJsonPath)

		var r Resource
		if err := Unmarshal(JSON, b, &r); err != nil {
			t.Fatalf("failed to unmarshal data to resource: %v", err)
		}
	})

	t.Run("JSONLink", func(t *testing.T) {
		b := MustReadFile(t, linkJsonPath)

		var l Link
		if err := Unmarshal(JSON, b, &l); err != nil {
			t.Fatalf("failed to unmarshal data to link: %v", err)
		}
	})
}
