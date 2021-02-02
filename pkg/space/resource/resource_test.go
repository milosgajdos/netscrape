package resource

import (
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
)

const (
	name    = "ResName"
	group   = "ResGroup"
	version = "ResVersion"
	kind    = "ResKind"
	ns      = false
)

func TestNew(t *testing.T) {
	r, err := New(name, group, version, kind, ns)
	if err != nil {
		t.Fatalf("failed creating new resource: %v", err)
	}

	if n := r.Name(); n != name {
		t.Errorf("expected name: %s, got: %s", name, n)
	}

	if g := r.Group(); g != group {
		t.Errorf("expected group: %s, got: %s", group, g)
	}

	if v := r.Version(); v != version {
		t.Errorf("expected version: %s, got: %s", version, v)
	}

	if k := r.Kind(); k != kind {
		t.Errorf("expected kind: %s, got: %s", kind, k)
	}

	if n := r.Namespaced(); n != ns {
		t.Errorf("expected namespaced: %v, got: %v", ns, n)
	}
}

func TestNewWithOptions(t *testing.T) {
	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	r, err := New(name, group, version, kind, ns, WithAttrs(a))
	if err != nil {
		t.Fatalf("failed creating new resource: %v", err)
	}

	if val := r.Attrs().Get(k); val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
