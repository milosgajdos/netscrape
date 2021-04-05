package resource

import (
	"context"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

const (
	testUID     = "resUID"
	testType    = "resType"
	testName    = "resName"
	testGroup   = "resGroup"
	testVersion = "resVersion"
	testKind    = "resKind"
	testNs      = false
	testDOTID   = "dotID"
)

func MustSet(ctx context.Context, a attrs.Attrs, k, v string, t *testing.T) {
	if err := a.Set(ctx, k, v); err != nil {
		t.Fatalf("failed to set val %s for key %s: %v", k, v, err)
	}
}

func TestNew(t *testing.T) {
	r, err := New(testType, testName, testGroup, testVersion, testKind, testNs)
	if err != nil {
		t.Fatalf("failed creating new resource: %v", err)
	}

	if n := r.Name(); n != testName {
		t.Errorf("expected name: %s, got: %s", testName, n)
	}

	if g := r.Group(); g != testGroup {
		t.Errorf("expected group: %s, got: %s", testGroup, g)
	}

	if v := r.Version(); v != testVersion {
		t.Errorf("expected version: %s, got: %s", testVersion, v)
	}

	if k := r.Kind(); k != testKind {
		t.Errorf("expected kind: %s, got: %s", testKind, k)
	}

	if n := r.Namespaced(); n != testNs {
		t.Errorf("expected namespaced: %v, got: %v", testNs, n)
	}
}

func TestNewWithOptions(t *testing.T) {
	a := memattrs.New()
	k, v := "foo", "bar"
	MustSet(context.Background(), a, k, v, t)

	uid := memuid.NewFromString(testUID)

	r, err := New(testType, testName, testGroup, testVersion, testKind, testNs, WithUID(uid), WithDOTID(testDOTID), WithAttrs(a))
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

	if u := r.UID().String(); u != testUID {
		t.Errorf("expected resource uid: %s, got: %s", testUID, u)
	}

	if d := r.DOTID(); d != testDOTID {
		t.Errorf("expected dotid: %s, got: %s", testDOTID, d)
	}
}
