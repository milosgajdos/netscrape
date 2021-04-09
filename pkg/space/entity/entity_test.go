package entity

import (
	"context"
	"testing"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func TestNew(t *testing.T) {
	e, err := New(entType)
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.Type() != entType {
		t.Errorf("expected type: %s, got: %s", entType, e.Type())
	}
}

func TestNewWithOptions(t *testing.T) {
	uid := memuid.NewFromString(entUID)

	e, err := New(entType, WithUID(uid), WithDOTID(testDOTID))
	if err != nil {
		t.Fatalf("failed creating new entity: %v", err)
	}

	if e.UID().String() != entUID {
		t.Errorf("expected entity uid: %s, got: %s", entUID, e.UID().String())
	}

	if d := e.DOTID(); d != testDOTID {
		t.Errorf("expected dotid: %s, got: %s", testDOTID, d)
	}

	dotid2 := "dotid2"
	e.SetDOTID(dotid2)

	if d := e.DOTID(); d != dotid2 {
		t.Errorf("expected dotid: %s, got: %s", dotid2, d)
	}

	a := memattrs.New()

	k, v := "foo", "bar"
	MustSet(context.Background(), a, k, v, t)

	e, err = New(entType, WithAttrs(a))
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
