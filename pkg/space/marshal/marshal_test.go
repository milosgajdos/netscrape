package marshal

import (
	"testing"

	"github.com/milosgajdos/netscrape/pkg/internal"
)

func TestEntity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("EntityToSpace", func(t *testing.T) {
		e := Entity{
			UID:   internal.ResUID,
			Type:  internal.ResType,
			Attrs: map[string]string{"foo": "bar"},
		}

		if _, err := EntityToSpace(e); err != nil {
			t.Fatalf("error marshaling entity to space: %v", err)
		}
	})

	t.Run("EntityFromSpace", func(t *testing.T) {
		e := MustEntity(t)

		if _, err := EntityFromSpace(e); err != nil {
			t.Fatalf("error marshaling space to entity: %v", err)
		}
	})
}

func TestResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("ResourceToSpace", func(t *testing.T) {
		r := Resource{
			Entity: Entity{
				UID:   internal.ResUID,
				Type:  internal.ResType,
				Attrs: map[string]string{"foo": "bar"},
			},
			Name:       internal.ResName,
			Group:      internal.ResGroup,
			Version:    internal.ResVersion,
			Kind:       internal.ResKind,
			Namespaced: internal.ResNsd,
		}

		if _, err := ResourceToSpace(r); err != nil {
			t.Fatalf("error marshaling resource to space: %v", err)
		}
	})

	t.Run("ResourceFromSpace", func(t *testing.T) {
		r := MustResource(t)

		if _, err := ResourceFromSpace(r); err != nil {
			t.Fatalf("error marshaling space to resource: %v", err)
		}
	})
}

func TestObject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("ObjectToSpace", func(t *testing.T) {
		r := Resource{
			Entity: Entity{
				UID:   internal.ResUID,
				Type:  internal.ResType,
				Attrs: map[string]string{"foo": "bar"},
			},
			Name:       internal.ResName,
			Group:      internal.ResGroup,
			Version:    internal.ResVersion,
			Kind:       internal.ResKind,
			Namespaced: internal.ResNsd,
		}

		o := Object{
			Entity: Entity{
				UID:   internal.ObjUID,
				Type:  internal.ObjType,
				Attrs: map[string]string{"foo": "bar"},
			},
			Name:      internal.ObjName,
			Namespace: internal.ObjNs,
			Resource:  &r,
		}

		if _, err := ObjectToSpace(o); err != nil {
			t.Fatalf("error marshaling object to space: %v", err)
		}
	})

	t.Run("ObjectFromSpace", func(t *testing.T) {
		e := MustObject(t)

		if _, err := ObjectFromSpace(e); err != nil {
			t.Fatalf("error marshaling space to object: %v", err)
		}
	})
}

func TestLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("LinkToSpace", func(t *testing.T) {
		l := Link{
			UID:   internal.LinkUID,
			From:  internal.LinkFrom,
			To:    internal.LinkTo,
			Attrs: map[string]string{"foo": "bar"},
		}

		if _, err := LinkToSpace(l); err != nil {
			t.Fatalf("error marshaling link to space: %v", err)
		}
	})

	t.Run("LinkFromSpace", func(t *testing.T) {
		l := MustLink(t)

		if _, err := LinkFromSpace(l); err != nil {
			t.Fatalf("error marshalling space to link: %v", err)
		}
	})
}
