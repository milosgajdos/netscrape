package entity

import (
	"context"
	"reflect"
	"testing"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func TestNewObject(t *testing.T) {
	r, err := NewResource(resType, resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	o, err := NewObject(objType, objName, objNs, r)
	if err != nil {
		t.Fatalf("failed creating new object: %v", err)
	}

	if o.Type() != objType {
		t.Errorf("expected type: %s, got: %s", objType, o.Type())
	}

	if o.Name() != objName {
		t.Errorf("expected name: %s, got: %s", objName, o.Name())
	}

	if o.Namespace() != objNs {
		t.Errorf("expected namespace: %s, got: %s", objNs, o.Namespace())
	}

	if !reflect.DeepEqual(o.Resource(), r) {
		t.Errorf("expected resource: %v, got: %v", r, o.Resource())
	}

	if _, err = NewObject(objType, objName, objNs, nil); err != nil {
		t.Fatalf("failed creating new object: %v", err)
	}
}

func TestNewObjectWithOptions(t *testing.T) {
	r, err := NewResource(resType, resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	uid := memuid.NewFromString(objUID)

	o, err := NewObject(objType, objName, objNs, r, WithUID(uid), WithDOTID(testDOTID))
	if err != nil {
		t.Fatalf("failed creating new object: %v", err)
	}

	if o.UID().String() != objUID {
		t.Errorf("expected object uid: %s, got: %s", objUID, o.UID().String())
	}

	if d := o.DOTID(); d != testDOTID {
		t.Errorf("expected dotid: %s, got: %s", testDOTID, d)
	}

	dotid2 := "dotid2"
	o.SetDOTID(dotid2)

	if d := o.DOTID(); d != dotid2 {
		t.Errorf("expected dotid: %s, got: %s", dotid2, d)
	}

	a := memattrs.New()

	k, v := "foo", "bar"
	MustSet(context.Background(), a, k, v, t)

	o, err = NewObject(objType, objName, objNs, r, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new object: %v", err)
	}

	val, err := o.Attrs().Get(context.Background(), k)
	if err != nil {
		t.Fatalf("failed to get val for key %s: %v", k, err)
	}

	if val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
