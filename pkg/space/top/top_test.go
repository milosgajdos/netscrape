package top

import (
	"context"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/space/plan"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	objPath = "../testdata/undirected/objects.yaml"
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

func TestObjects(t *testing.T) {
	top, err := newTop(resPath, objPath)
	if err != nil {
		t.Errorf("failed to create mock Top: %v", err)
		return
	}

	objects, err := top.Objects(context.TODO())
	if err != nil {
		t.Fatalf("failed to get objects: %v", err)
	}

	if len(objects) == 0 {
		t.Errorf("no objects found")
	}
}

func TestGetUID(t *testing.T) {
	top, err := newTop(resPath, objPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	objects, err := top.Objects(context.TODO())
	if err != nil {
		t.Fatalf("failed to get objects: %v", err)
	}

	uids := make([]uuid.UID, len(objects))

	for i, o := range objects {
		uids[i] = o.UID()
	}

	for _, uid := range uids {
		q := base.Build().Add(query.UID(uid), query.UUIDEqFunc(uid))

		objects, err := top.Get(context.TODO(), q)

		if err != nil {
			t.Errorf("error getting object: %s: %v", uid, err)
			continue
		}

		if len(objects) != 1 {
			t.Errorf("expected single %s object, got: %d", uid, len(objects))
			continue
		}

		if objects[0].UID().Value() != uid.Value() {
			t.Errorf("expected: %s, got: %s", uid, objects[0].UID())
		}
	}
}

func TestTopGet(t *testing.T) {
	top, err := newTop(resPath, objPath)
	if err != nil {
		t.Fatalf("failed to create mock Top: %v", err)
	}

	q := base.Build()

	objects, err := top.Get(context.TODO(), q)
	if err != nil {
		t.Errorf("error querying objects: %v", err)
	}

	allObjects, err := top.Objects(context.TODO())
	if err != nil {
		t.Fatalf("failed to get all topology objects: %v", err)
	}

	if len(objects) != len(allObjects) {
		t.Errorf("expected: %d, got: %d", len(objects), len(allObjects))

	}

	namespaces := make([]string, len(allObjects))
	kinds := make([]string, len(allObjects))
	names := make([]string, len(allObjects))

	for i, o := range allObjects {
		namespaces[i] = o.Namespace()
		kinds[i] = o.Resource().Kind()
		names[i] = o.Name()
	}

	for _, ns := range namespaces {
		q := base.Build().Add(query.Namespace(ns), query.StringEqFunc(ns))

		objects, err := top.Get(context.TODO(), q)
		if err != nil {
			t.Errorf("error getting namespace %s objects: %v", ns, err)
			continue
		}

		for _, o := range objects {
			if o.Namespace() != ns {
				t.Errorf("expected: %s, got: %s", ns, o.Namespace())
			}
		}

		for _, kind := range kinds {
			q = q.Add(query.Kind(kind), query.StringEqFunc(kind))

			objects, err = top.Get(context.TODO(), q)
			if err != nil {
				t.Errorf("error getting objects: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, o := range objects {
				if o.Namespace() != ns || o.Resource().Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, o.Namespace(), o.Resource().Kind())
				}
			}

			for _, name := range names {
				q = q.Add(query.Name(name), query.StringEqFunc(name))

				objects, err = top.Get(context.TODO(), q)
				if err != nil {
					t.Errorf("error getting objects: %s/%s/%s: %v", ns, kind, name, err)
					continue
				}

				for _, o := range objects {
					if o.Namespace() != ns || o.Resource().Kind() != kind || o.Name() != name {
						t.Errorf("expected: %s/%s/%s, got: %s/%s/%s", ns, kind, name,
							o.Namespace(), o.Resource().Kind(), o.Name())
					}
				}
			}
		}
	}
}
