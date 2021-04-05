package entity

import (
	"context"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space/resource"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
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

func MustSet(ctx context.Context, a attrs.Attrs, k, v string, t *testing.T) {
	if err := a.Set(ctx, k, v); err != nil {
		t.Fatalf("failed to set val %s for key %s: %v", k, v, err)
	}
}

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

	uid := memuid.NewFromString("partialUID")

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

	uid := memuid.NewFromString(entUID)

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

	a := memattrs.New()

	k, v := "foo", "bar"
	MustSet(context.Background(), a, k, v, t)

	e, err = New(entType, entName, entNs, r, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new entity: %v", err)
	}

	val, err := e.Attrs().Get(context.Background(), k)
	if err != nil {
		t.Fatalf("failed to get val for key %s: %v", k, err)
	}

	if val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
