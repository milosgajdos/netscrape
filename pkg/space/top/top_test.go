package top

import (
	"context"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	entPath = "../testdata/undirected/entities.yaml"
)

func newTop(entPath string) (*Top, error) {
	top, err := NewMock(entPath)
	if err != nil {
		return nil, err
	}

	return top, nil
}

func TestNew(t *testing.T) {
	top, err := New()
	if err != nil {
		t.Fatalf("failed to create space Top: %v", err)
	}

	if top == nil {
		t.Fatalf("expected new Top, got: %v", top)
	}
}

func TestEntities(t *testing.T) {
	top, err := newTop(entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	entities, err := top.Entities(context.Background())
	if err != nil {
		t.Fatalf("failed to get entities: %v", err)
	}

	if len(entities) == 0 {
		t.Errorf("no entities found")
	}
}

func TestGetUID(t *testing.T) {
	top, err := newTop(entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	entities, err := top.Entities(context.Background())
	if err != nil {
		t.Fatalf("failed to get entities: %v", err)
	}

	uids := make([]uuid.UID, len(entities))

	for i, o := range entities {
		uids[i] = o.UID()
	}

	for _, uid := range uids {
		q := base.Build().Add(predicate.UID(uid))

		entities, err := top.Get(context.Background(), q)

		if err != nil {
			t.Errorf("error getting entity: %s: %v", uid, err)
			continue
		}

		if len(entities) != 1 {
			t.Errorf("expected single %s entity, got: %d", uid, len(entities))
			continue
		}

		if entities[0].UID().Value() != uid.Value() {
			t.Errorf("expected: %s, got: %s", uid, entities[0].UID())
		}
	}
}

func TestGet(t *testing.T) {
	top, err := newTop(entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	q := base.Build()

	entities, err := top.Get(context.Background(), q)
	if err != nil {
		t.Errorf("error querying entities: %v", err)
	}

	// empty query should return empty slice
	expCount := 0
	if count := len(entities); count != expCount {
		t.Errorf("expected: %d, got: %d", expCount, count)

	}

	entities, err = top.Entities(context.Background())
	if err != nil {
		t.Fatalf("failed to get all topology entities: %v", err)
	}

	types := make([]entity.Type, len(entities))

	for i, e := range entities {
		types[i] = e.Type()
	}

	for _, ent := range types {
		q := base.Build().Add(predicate.Entity(ent))

		entities, err := top.Get(context.Background(), q)
		if err != nil {
			t.Errorf("error getting entities with type %s : %v", ent, err)
			continue
		}

		for _, e := range entities {
			if n := e.Type(); n != ent {
				t.Errorf("expected: %s, got: %s", ent, n)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	top, err := newTop(entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	r, err := newTestResource(resName, resGroup, resVersion, resKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	e, err := newTestEntity(entUID, entName, entNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	ctx := context.Background()

	if err := top.Add(ctx, e); err != nil {
		t.Fatalf("failed adding entitty %v to top: %v", e.UID().Value(), err)
	}

	q := base.Build().Add(predicate.UID(e.UID()))

	ents, err := top.Get(ctx, q)
	if err != nil {
		t.Fatalf("failed getting entity %v: %v", e.UID(), err)
	}

	expCount := 1
	if count := len(ents); count != expCount {
		t.Fatalf("expected entitites: %d, got: %d", expCount, count)
	}

	if err := top.Remove(ctx, e.UID()); err != nil {
		t.Fatalf("failed removing entitty %v from top: %v", e.UID().Value(), err)
	}

	ents, err = top.Get(ctx, q)
	if err != nil {
		t.Fatalf("failed getting entity %v: %v", e.UID(), err)
	}

	expCount = 0
	if count := len(ents); count != expCount {
		t.Fatalf("expected entitites: %d, got: %d", expCount, count)
	}
}

func TestLinks(t *testing.T) {
	top, err := newTop(entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	r, err := newTestResource(resName, resGroup, resVersion, resKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	e1, err := newTestEntity(entUID, entName, entNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	e2, err := newTestEntity(ent2UID, ent2Name, ent2Ns, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	ctx := context.Background()

	if err := top.Add(ctx, e1); err != nil {
		t.Fatalf("failed adding entitty %v to top: %v", e1.UID().Value(), err)
	}

	if err := top.Add(ctx, e2); err != nil {
		t.Fatalf("failed adding entitty %v to top: %v", e2.UID().Value(), err)
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	if err := top.Link(ctx, e1.UID(), e2.UID(), space.WithAttrs(a), space.WithMerge(true)); err != nil {
		t.Fatalf("failed linking entities: %v", err)
	}

	links, err := top.Links(ctx, e1.UID())
	if err != nil {
		t.Fatalf("failed gettings %v links: %v", e1.UID(), err)
	}

	expCount := 1
	if count := len(links); count != expCount {
		t.Fatalf("expected links: %d, got: %d", expCount, count)
	}
}
