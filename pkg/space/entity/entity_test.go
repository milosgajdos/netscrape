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
	entUID     = "testID"
	entName    = "testName"
	entNs      = "testNs"
	toUID      = "toUID"
)

func newTestResource(name, group, version, kind string, namespaced bool, opts ...resource.Option) (space.Resource, error) {
	return resource.New(name, group, version, kind, namespaced, opts...)
}

func TestNew(t *testing.T) {
	r, err := newTestResource(resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	e, err := New(entName, entNs, r)
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.Name() != entName {
		t.Errorf("expected name: %s, got: %s", entName, e.Name())
	}

	if e.Namespace() != entNs {
		t.Errorf("expected namespace: %s, got: %s", entNs, e.Namespace())
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

	uid, err := uuid.NewFromString(entUID)
	if err != nil {
		t.Errorf("failed to create new uid: %v", err)
	}

	e, err := New(entName, entNs, r, WithUID(uid))
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.UID().Value() != entUID {
		t.Errorf("expected entity uid: %s, got: %s", entUID, e.UID().Value())
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	e, err = New(entName, entNs, r, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new entity: %v", err)
	}

	if val := e.Attrs().Get(k); val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
