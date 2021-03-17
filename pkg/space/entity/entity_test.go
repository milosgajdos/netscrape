package entity

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
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
)

func TestNew(t *testing.T) {
	r, err := resource.New(resType, resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	e, err := New(entType, entName, entNs, r)
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

	if _, err = New(entType, entName, entNs, nil); err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}
}

func TestNewWithOptions(t *testing.T) {
	r, err := resource.New(resType, resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	uid, err := uuid.NewFromString(entUID)
	if err != nil {
		t.Errorf("failed to create new uid: %v", err)
	}

	e, err := New(entType, entName, entNs, r, WithUID(uid), WithDOTID(entDOTID))
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.UID().String() != entUID {
		t.Errorf("expected entity uid: %s, got: %s", entUID, e.UID().String())
	}

	if d := e.DOTID(); d != entDOTID {
		t.Errorf("expected dotid: %s, got: %s", entDOTID, d)
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

	e, err = New(entType, entName, entNs, r, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new entity: %v", err)
	}

	if val := e.Attrs().Get(k); val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
