package resource

import (
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	testUID     = "resUID"
	testName    = "ResName"
	testGroup   = "ResGroup"
	testVersion = "ResVersion"
	testKind    = "ResKind"
	testNs      = false
	testDOTID   = "dotID"
)

func TestNew(t *testing.T) {
	r, err := New(testName, testGroup, testVersion, testKind, testNs)
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
	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	uid, err := uuid.NewFromString(testUID)
	if err != nil {
		t.Errorf("failed to create new uid: %v", err)
	}

	r, err := New(testName, testGroup, testVersion, testKind, testNs, WithUID(uid), WithDOTID(testDOTID), WithAttrs(a))
	if err != nil {
		t.Fatalf("failed creating new resource: %v", err)
	}

	if val := r.Attrs().Get(k); val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}

	if u := r.UID().Value(); u != testUID {
		t.Errorf("expected resource uid: %s, got: %s", testUID, u)
	}

	if d := r.DOTID(); d != testDOTID {
		t.Errorf("expected dotid: %s, got: %s", testDOTID, d)
	}
}
