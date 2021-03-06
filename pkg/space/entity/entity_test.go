package entity

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	resName    = "nodeResName"
	resGroup   = "nodeResGroup"
	resVersion = "nodeResVersion"
	resKind    = "nodeResKind"
	objUID     = "testID"
	objName    = "testName"
	objNs      = "testNs"
	dotid      = "dotID"
)

func newTestResource(name, group, version, kind string, namespaced bool, opts ...resource.Option) (space.Resource, error) {
	return resource.New(name, group, version, kind, namespaced, opts...)
}

func TestNew(t *testing.T) {
	r, err := newTestResource(resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	e, err := New(objName, objNs, r)
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.Name() != objName {
		t.Errorf("expected name: %s, got: %s", objName, e.Name())
	}

	if e.Namespace() != objNs {
		t.Errorf("expected namespace: %s, got: %s", objNs, e.Namespace())
	}

	if !reflect.DeepEqual(e.Resource(), r) {
		t.Errorf("expected resource: %v, got: %v", r, e.Resource())
	}
}

func TestNewWithOptions(t *testing.T) {
	r, err := newTestResource(resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	uid, err := uuid.NewFromString(objUID)
	if err != nil {
		t.Errorf("failed to create new uid: %v", err)
	}

	e, err := New(objName, objNs, r, WithUID(uid), WithDOTID(dotid))
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.UID().Value() != objUID {
		t.Errorf("expected entity uid: %s, got: %s", objUID, e.UID().Value())
	}

	if d := e.DOTID(); d != dotid {
		t.Errorf("expected dotid: %s, got: %s", dotid, d)
	}

	dotid2 := "dotid2"
	e.SetDOTID(dotid2)

	if d := e.DOTID(); d != dotid2 {
		t.Errorf("expected dotid: %s, got: %s", dotid2, d)
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	e, err = New(objName, objNs, r, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new entity: %v", err)
	}

	if val := e.Attrs().Get(k); val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
