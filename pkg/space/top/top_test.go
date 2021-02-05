package top

import (
	"context"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
	"github.com/milosgajdos/netscrape/pkg/space/plan"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	objPath = "../testdata/undirected/entities.yaml"
	resPath = "../testdata/undirected/resources.yaml"
)

func newTop(resPath, objPath string) (*Top, error) {
	src := "file:///" + resPath

	a, err := plan.NewMock(src)
	if err != nil {
		return nil, err
	}

	top, err := NewMock(a, objPath)
	if err != nil {
		return nil, err
	}

	return top, nil
}

func TestEntities(t *testing.T) {
	top, err := newTop(resPath, objPath)
	if err != nil {
		t.Errorf("failed to create mock Top: %v", err)
		return
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
	top, err := newTop(resPath, objPath)
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

func TestTopGet(t *testing.T) {
	top, err := newTop(resPath, objPath)
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
