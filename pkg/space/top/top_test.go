package top

import (
	"context"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/plan"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	entPath = "../testdata/undirected/entities.yaml"
	resPath = "../testdata/undirected/resources.yaml"
)

func newTop(resPath, entPath string) (*Top, error) {
	src := "file:///" + resPath

	p, err := plan.NewMock(src)
	if err != nil {
		return nil, err
	}

	top, err := NewMock(p, entPath)
	if err != nil {
		return nil, err
	}

	return top, nil
}

func TestNew(t *testing.T) {
	src := "file:///" + resPath

	p, err := plan.NewMock(src)
	if err != nil {
		t.Fatalf("failed to create space Plan: %v", err)
	}

	top, err := New(p)
	if err != nil {
		t.Fatalf("failed to create space Top: %v", err)
	}

	ctx := context.Background()

	plan, err := top.Plan(ctx)
	if err != nil {
		t.Fatalf("failed getting topology plan: %v", err)
	}

	o, err := plan.Origin(ctx)
	if err != nil {
		t.Fatalf("failed getting plan origin: %v", err)
	}

	if origin := o.URL().String(); origin != src {
		t.Errorf("expected origin: %s, got: %s", src, origin)
	}
}

func TestEntities(t *testing.T) {
	top, err := newTop(resPath, entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	entities, err := top.Entities(context.TODO())
	if err != nil {
		t.Fatalf("failed to get entities: %v", err)
	}

	if len(entities) == 0 {
		t.Errorf("no entities found")
	}
}

func TestGetUID(t *testing.T) {
	top, err := newTop(resPath, entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	entities, err := top.Entities(context.TODO())
	if err != nil {
		t.Fatalf("failed to get entities: %v", err)
	}

	uids := make([]uuid.UID, len(entities))

	for i, o := range entities {
		uids[i] = o.UID()
	}

	for _, uid := range uids {
		q := base.Build().Add(predicate.UID(uid))

		entities, err := top.Get(context.TODO(), q)

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
	top, err := newTop(resPath, entPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	q := base.Build()

	entities, err := top.Get(context.TODO(), q)
	if err != nil {
		t.Errorf("error querying entities: %v", err)
	}

	allEntities, err := top.Entities(context.TODO())
	if err != nil {
		t.Fatalf("failed to get all topology entities: %v", err)
	}

	if len(entities) != len(allEntities) {
		t.Errorf("expected: %d, got: %d", len(entities), len(allEntities))

	}

	namespaces := make([]string, len(allEntities))
	kinds := make([]string, len(allEntities))
	names := make([]string, len(allEntities))

	for i, o := range allEntities {
		namespaces[i] = o.Namespace()
		kinds[i] = o.Resource().Kind()
		names[i] = o.Name()
	}

	for _, ns := range namespaces {
		q := base.Build().Add(predicate.Namespace(ns))

		entities, err := top.Get(context.TODO(), q)
		if err != nil {
			t.Errorf("error getting namespace %s entities: %v", ns, err)
			continue
		}

		for _, o := range entities {
			if o.Namespace() != ns {
				t.Errorf("expected: %s, got: %s", ns, o.Namespace())
			}
		}

		for _, kind := range kinds {
			q = q.Add(predicate.Kind(kind))

			entities, err = top.Get(context.TODO(), q)
			if err != nil {
				t.Errorf("error getting entities: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, o := range entities {
				if o.Namespace() != ns || o.Resource().Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, o.Namespace(), o.Resource().Kind())
				}
			}

			for _, name := range names {
				q = q.Add(predicate.Name(name))

				entities, err = top.Get(context.TODO(), q)
				if err != nil {
					t.Errorf("error getting entities: %s/%s/%s: %v", ns, kind, name, err)
					continue
				}

				for _, o := range entities {
					if o.Namespace() != ns || o.Resource().Kind() != kind || o.Name() != name {
						t.Errorf("expected: %s/%s/%s, got: %s/%s/%s", ns, kind, name,
							o.Namespace(), o.Resource().Kind(), o.Name())
					}
				}
			}
		}
	}
}

func TestRemove(t *testing.T) {
	top, err := newTop(resPath, entPath)
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
	top, err := newTop(resPath, entPath)
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
