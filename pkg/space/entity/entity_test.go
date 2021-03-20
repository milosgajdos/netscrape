package entity

import (
	"reflect"
	"testing"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
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

func TestNewPartial(t *testing.T) {
	e, err := NewPartial()
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.Type() != PartialType {
		t.Errorf("expected type: %s, got: %s", PartialType, e.Type())
	}

	if e.Name() != PartialName {
		t.Errorf("expected name: %s, got: %s", PartialName, e.Name())
	}

	if e.Namespace() != PartialNs {
		t.Errorf("expected namespace: %s, got: %s", PartialNs, e.Namespace())
	}

	if r := e.Resource(); r != nil {
		t.Errorf("expected nil resource, got: %v", r)
	}

	uid, err := uuid.NewFromString("partialUID")
	if err != nil {
		t.Fatalf("failed to created uid: %v", err)
	}

	e, err = NewPartial(WithUID(uid))
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if u := e.UID().String(); u != uid.String() {
		t.Errorf("expected uid: %v, got: %v", uid.String(), u)
	}
}

func TestNew(t *testing.T) {
	r, err := resource.New(resType, resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	e, err := New(entType, entName, entNs, r)
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.Type() != entType {
		t.Errorf("expected type: %s, got: %s", entType, e.Type())
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

	a, err := memattrs.New()
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
