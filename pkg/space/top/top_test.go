package top

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	entPath = "../testdata/undirected/entities.yaml"
)

func MustNewTop(src string, t *testing.T) *Top {
	p, err := NewMock(src)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}
	return p
}

func MustTestEntity(uid, name string, t *testing.T) space.Entity {
	r, err := newTestResource(resName, resGroup, resVersion, resKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	e, err := newTestEntity(uid, name, entNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", uid, err)
	}

	return e
}

func MustMakeEntities(count int, t *testing.T) []space.Entity {
	ents := make([]space.Entity, count)

	for i := 0; i < count; i++ {
		uid := fmt.Sprintf("uid%d", i)
		name := fmt.Sprintf("name%d", i)

		ents[i] = MustTestEntity(uid, name, t)
	}

	return ents
}

func TestAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)
		e := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e); err != nil {
			t.Errorf("failed adding entity %s: %v", e.UID(), err)
		}
	})
}

func TestGetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		rx, err := p.GetAll(context.Background())
		if err != nil {
			t.Fatalf("failed getting all entity: %v", err)
		}

		exp := 10
		if c := len(rx); c != exp {
			t.Errorf("expected entitys: %d, got: %d", exp, c)
		}
	})
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)
		e := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e); err != nil {
			t.Fatalf("failed adding entity %s: %v", e.UID(), err)
		}

		res, err := p.Get(context.Background(), e.UID())
		if err != nil {
			t.Fatalf("failed getting entity %s: %v", e.UID(), err)
		}

		if !reflect.DeepEqual(res.UID(), e.UID()) {
			t.Errorf("expected entity: %s, got: %s", e.UID(), res.UID())
		}
	})

	t.Run("ErrEntityNotFound", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		uid, err := uuid.New()
		if err != nil {
			t.Fatalf("failed to generate uid: %v", err)
		}
		if _, err := p.Get(context.Background(), uid); !errors.Is(err, space.ErrEntityNotFound) {
			t.Errorf("expected error: %v, got: %v", space.ErrEntityNotFound, err)
		}
	})
}

func TestDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)
		e := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e); err != nil {
			t.Fatalf("failed adding entity %s: %v", e.UID(), err)
		}

		if err := p.Delete(context.Background(), e.UID()); err != nil {
			t.Fatalf("failed removing entity %s: %v", e.UID(), err)
		}

		if _, err := p.Get(context.Background(), e.UID()); !errors.Is(err, space.ErrEntityNotFound) {
			t.Errorf("expected %v: got: %v", space.ErrEntityNotFound, err)
		}
	})
}

func TestLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		e1 := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e1); err != nil {
			t.Fatalf("failed storing entity %s: %v", e1.UID(), err)
		}

		e2 := MustTestEntity("foo2UID", "foo2Name", t)

		if err := p.Add(context.Background(), e2); err != nil {
			t.Fatalf("failed storing entity %s: %v", e2.UID(), err)
		}

		if err := p.Link(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Fatalf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
		}

		a, err := attrs.New()
		if err != nil {
			t.Fatalf("failed to create new attrs: %v", err)
		}
		k, v := "foo", "bar"
		a.Set(k, v)

		if err := p.Link(context.Background(), e1.UID(), e2.UID(), space.WithAttrs(a), space.WithMerge()); err != nil {
			t.Errorf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
		}
	})
}

func TestLinks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		e1 := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e1); err != nil {
			t.Fatalf("failed storing entity %s: %v", e1.UID(), err)
		}

		e2 := MustTestEntity("foo2UID", "foo2Name", t)

		if err := p.Add(context.Background(), e2); err != nil {
			t.Fatalf("failed storing entity %s: %v", e2.UID(), err)
		}

		if err := p.Link(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Fatalf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
		}

		links, err := p.Links(context.Background(), e1.UID())
		if err != nil {
			t.Fatalf("failed gettings %v links: %v", e1.UID(), err)
		}
		expCount := 1
		if count := len(links); count != expCount {
			t.Fatalf("expected links: %d, got: %d", expCount, count)
		}
	})

	t.Run("ErrNoLinksFound", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		uid, err := uuid.New()
		if err != nil {
			t.Fatalf("failed to generate uid: %v", err)
		}
		if _, err := p.Links(context.Background(), uid); !errors.Is(err, space.ErrNoLinksFound) {
			t.Errorf("expected error: %v, got: %v", space.ErrNoLinksFound, err)
		}
	})
}

func TestUnlink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		e1 := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e1); err != nil {
			t.Fatalf("failed storing entity %s: %v", e1.UID(), err)
		}

		e2 := MustTestEntity("foo2UID", "foo2Name", t)

		if err := p.Add(context.Background(), e2); err != nil {
			t.Fatalf("failed storing entity %s: %v", e2.UID(), err)
		}

		if err := p.Link(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Fatalf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
		}

		if err := p.Unlink(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Errorf("failed unlinking %v to %v: %v", e1.UID(), e2.UID(), err)
		}
	})

	t.Run("UnlinkNonExisten", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		e1 := MustTestEntity("foo1UID", "foo1Name", t)
		e2 := MustTestEntity("foo2UID", "foo2Name", t)

		if err := p.Unlink(context.Background(), e1.UID(), e2.UID()); err != nil {
			t.Errorf("failed unlinking %v to %v: %v", e1.UID(), e2.UID(), err)
		}
	})
}

func TestBulkAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		ents := MustMakeEntities(5, t)

		if err := p.BulkAdd(context.Background(), ents); err != nil {
			t.Errorf("failed storing entities: %v", err)
		}
	})
}

func TestBulkGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		ents := MustMakeEntities(5, t)

		if err := p.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		uids := make([]uuid.UID, len(ents))
		for i, e := range ents {
			uids[i] = e.UID()
		}

		sents, err := p.BulkGet(context.Background(), uids)
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
		p := MustNewTop(entPath, t)

		ents := MustMakeEntities(5, t)

		if err := p.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		uids := make([]uuid.UID, len(ents))

		for i, e := range ents {
			uids[i] = e.UID()
		}

		if err := p.BulkDelete(context.Background(), uids); err != nil {
			t.Fatalf("failed deleting entities: %v", err)
		}

		for _, uid := range uids {
			if _, err := p.Get(context.Background(), uid); !errors.Is(err, space.ErrEntityNotFound) {
				t.Errorf("expected %v: got: %v", space.ErrEntityNotFound, err)
			}
		}
	})
}

func TestBulkLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		ents := MustMakeEntities(5, t)

		if err := p.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		e := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e); err != nil {
			t.Fatalf("failed storing entity %s: %v", e.UID(), err)
		}

		uids := make([]uuid.UID, len(ents))

		for i, e := range ents {
			uids[i] = e.UID()
		}

		if err := p.BulkLink(context.Background(), e.UID(), uids); err != nil {
			t.Errorf("failed bulk-linking %v: %v", e.UID(), err)
		}
	})
}

func TestBulkUnlink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewTop(entPath, t)

		ents := MustMakeEntities(5, t)

		if err := p.BulkAdd(context.Background(), ents); err != nil {
			t.Fatalf("failed storing entities: %v", err)
		}

		e := MustTestEntity("foo1UID", "foo1Name", t)

		if err := p.Add(context.Background(), e); err != nil {
			t.Fatalf("failed storing entity %s: %v", e.UID(), err)
		}

		uids := make([]uuid.UID, len(ents))

		for i, e := range ents {
			uids[i] = e.UID()
		}

		if err := p.BulkLink(context.Background(), e.UID(), uids); err != nil {
			t.Fatalf("failed bulk-linking %v: %v", e.UID(), err)
		}

		if err := p.BulkUnlink(context.Background(), e.UID(), uids); err != nil {
			t.Errorf("failed bulk-unlinking %v: %v", e.UID(), err)
		}
	})
}
