package entity

import (
	"context"
	"testing"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func TestNewResource(t *testing.T) {
	r, err := NewResource(resType, resName, resGroup, resVersion, resKind, resNsd)
	if err != nil {
		t.Fatalf("failed creating new resource: %v", err)
	}

	if n := r.Name(); n != resName {
		t.Errorf("expected name: %s, got: %s", resName, n)
	}

	if g := r.Group(); g != resGroup {
		t.Errorf("expected group: %s, got: %s", resGroup, g)
	}

	if v := r.Version(); v != resVersion {
		t.Errorf("expected version: %s, got: %s", resVersion, v)
	}

	if k := r.Kind(); k != resKind {
		t.Errorf("expected kind: %s, got: %s", resKind, k)
	}

	if n := r.Namespaced(); n != resNsd {
		t.Errorf("expected namespaced: %v, got: %v", resNsd, n)
	}
}

func TestNewWResourceithOptions(t *testing.T) {
	a := memattrs.New()
	k, v := "foo", "bar"
	MustSet(context.Background(), a, k, v, t)

	uid := memuid.NewFromString(resUID)

	r, err := NewResource(resType, resName, resGroup, resVersion, resKind, resNsd, WithUID(uid), WithDOTID(testDOTID), WithAttrs(a))
	if err != nil {
		t.Fatalf("failed creating new resource: %v", err)
	}

	val, err := r.Attrs().Get(context.Background(), k)
	if err != nil {
		t.Fatalf("failed to get value for key %s: %v", k, err)
	}

	if val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}

	if u := r.UID().String(); u != resUID {
		t.Errorf("expected resource uid: %s, got: %s", resUID, u)
	}

	if d := r.DOTID(); d != testDOTID {
		t.Errorf("expected dotid: %s, got: %s", testDOTID, d)
	}
}
