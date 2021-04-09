package memory

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/internal"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"

	memgraph "github.com/milosgajdos/netscrape/pkg/graph/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func MustNewStore(t *testing.T, opts ...Option) *Memory {
	s, err := NewStore(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func MustTestEntity(typ, name string, t *testing.T, opts ...entity.Option) store.Entity {
	e, err := internal.NewTestObject(opts...)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	return e
}

func MustMakeEntities(count int, t *testing.T) []store.Entity {
	ents := make([]store.Entity, count)

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("name%d", i)

		ents[i] = MustTestEntity("fooType", name, t)
	}
	return ents
}

func TestNewStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("NoOpts", func(t *testing.T) {
		if _, err := NewStore(); err != nil {
			t.Fatalf("failed creating new store: %v", err)
		}
	})

	t.Run("WithOpts", func(t *testing.T) {
		uid := "storeUID"
		suid := memuid.NewFromString(uid)

		s, err := NewStore(WithUID(suid))
		if err != nil {
			t.Fatalf("failed creating new store: %v", err)
		}

		if u := s.UID().String(); u != uid {
			t.Errorf("expected uid: %s, got: %s", uid, u)
		}

		g, err := memgraph.NewWUG()
		if err != nil {
			t.Fatalf("failed creating new memory graph: %v", err)
		}

		sg, err := NewStore(WithGraph(g))
		if err != nil {
			t.Fatalf("failed creating new store: %v", err)
		}

		if _, err := sg.Graph(); err != nil {
			t.Errorf("failed to get store graph handle: %v", err)
		}
	})
}

func TestAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		uid := memuid.NewFromString("someUID")
		s := MustNewStore(t, WithUID(uid))

		e := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e); err != nil {
			t.Errorf("failed storing entity %s: %v", e.UID(), err)
		}
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		uid := memuid.NewFromString("someUID")
		s := MustNewStore(t, WithUID(uid))

		e := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e); err != nil {
			t.Errorf("failed storing entity %s: %v", e.UID(), err)
		}

		// Add the same entity without Upsert
		if err := s.Add(context.Background(), e); !errors.Is(err, store.ErrAlreadyExists) {
			t.Errorf("expected error: %v, got: %v", graph.ErrDuplicateNode, err)
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		uid := memuid.NewFromString("someUID")
		s := MustNewStore(t, WithUID(uid))

		e := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e); err != nil {
			t.Errorf("failed storing entity %s: %v", e.UID(), err)
		}

		ex := MustTestEntity("fooType", "foo1Name", t, entity.WithDOTID("someDOTID"))

		if err := s.Add(context.Background(), ex, store.WithUpsert()); err != nil {
			t.Errorf("failed storing entity %s: %v", e.UID(), err)
		}
	})
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		e := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e); err != nil {
			t.Fatalf("failed storing entity %s: %v", e.UID(), err)
		}

		res, err := s.Get(context.Background(), e.UID())
		if err != nil {
			t.Errorf("failed getting entity %s: %v", e.UID(), err)
		}

		if !reflect.DeepEqual(res.UID(), e.UID()) {
			t.Errorf("expected entity: %s, got: %s", e.UID(), res.UID())
		}
	})

	t.Run("ErrEntityNotFound", func(t *testing.T) {
		s := MustNewStore(t)

		if _, err := s.Get(context.Background(), memuid.New()); !errors.Is(err, store.ErrEntityNotFound) {
			t.Errorf("expected error: %v, got: %v", store.ErrEntityNotFound, err)
		}
	})
}

func TestDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		e := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e); err != nil {
			t.Fatalf("failed storing entity %s: %v", e.UID(), err)
		}

		if err := s.Delete(context.Background(), e.UID()); err != nil {
			t.Fatalf("failed deleting entity %s: %v", e.UID(), err)
		}

		if _, err := s.Get(context.Background(), e.UID()); !errors.Is(err, store.ErrEntityNotFound) {
			t.Errorf("expected %v: got: %v", store.ErrEntityNotFound, err)
		}
	})
}

func TestLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		e1 := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e1); err != nil {
			t.Fatalf("failed storing entity %s: %v", e1.UID(), err)
		}

		e2 := MustTestEntity("fooType", "foo2Name", t)

		if err := s.Add(context.Background(), e2); err != nil {
			t.Fatalf("failed storing entity %s: %v", e2.UID(), err)
		}

		if err := s.Link(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Errorf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
		}
	})
}

func TestUnlink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		e1 := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e1); err != nil {
			t.Fatalf("failed storing entity %s: %v", e1.UID(), err)
		}

		e2 := MustTestEntity("fooType", "foo2Name", t)

		if err := s.Add(context.Background(), e2); err != nil {
			t.Fatalf("failed storing entity %s: %v", e2.UID(), err)
		}

		if err := s.Link(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Fatalf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
		}

		if err := s.Unlink(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Errorf("failed unlinking %v to %v: %v", e1.UID(), e2.UID(), err)
		}
	})
}

func TestBulkAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		ents := MustMakeEntities(5, t)

		if err := s.BulkAdd(context.Background(), ents); err != nil {
			t.Errorf("failed storing entities: %v", err)
		}
	})
}

func TestBulkGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		ents := MustMakeEntities(5, t)

		if err := s.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		uids := make([]uuid.UID, len(ents))
		for i, e := range ents {
			uids[i] = e.UID()
		}

		sents, err := s.BulkGet(context.Background(), uids)
		if err != nil {
			t.Fatalf("failed getting entities: %v", err)
		}

		if len(sents) != len(ents) {
			t.Errorf("expected %d entities, got: %d", len(sents), len(ents))
		}
	})
}

func TestBulkDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		ents := MustMakeEntities(5, t)

		if err := s.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		uids := make([]uuid.UID, len(ents))

		for i, e := range ents {
			uids[i] = e.UID()
		}

		if err := s.BulkDelete(context.Background(), uids); err != nil {
			t.Fatalf("failed deleting entities: %v", err)
		}

		for _, uid := range uids {
			if _, err := s.Get(context.Background(), uid); !errors.Is(err, store.ErrEntityNotFound) {
				t.Errorf("expected %v: got: %v", store.ErrEntityNotFound, err)
			}
		}
	})
}

func TestBulkLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		ents := MustMakeEntities(5, t)

		if err := s.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		e := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e); err != nil {
			t.Fatalf("failed storing entity %s: %v", e.UID(), err)
		}

		uids := make([]uuid.UID, len(ents))

		for i, e := range ents {
			uids[i] = e.UID()
		}

		if err := s.BulkLink(context.Background(), e.UID(), uids); err != nil {
			t.Errorf("failed bulk-linking %v: %v", e.UID(), err)
		}
	})
}

func TestBulkUnlink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		s := MustNewStore(t)

		ents := MustMakeEntities(5, t)

		if err := s.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		e := MustTestEntity("fooType", "foo1Name", t)

		if err := s.Add(context.Background(), e); err != nil {
			t.Fatalf("failed storing entity %s: %v", e.UID(), err)
		}

		uids := make([]uuid.UID, len(ents))

		for i, e := range ents {
			uids[i] = e.UID()
		}

		if err := s.BulkLink(context.Background(), e.UID(), uids); err != nil {
			t.Fatalf("failed bulk-linking %v: %v", e.UID(), err)
		}

		if err := s.BulkUnlink(context.Background(), e.UID(), uids); err != nil {
			t.Errorf("failed bulk-unlinking %v: %v", e.UID(), err)
		}
	})
}
