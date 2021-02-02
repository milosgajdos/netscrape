package object

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
	objUID     = "testID"
	objName    = "testName"
	objNs      = "testNs"
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

	o, err := New(objName, objNs, r)
	if err != nil {
		t.Fatalf("failed creating new object: %v", err)
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
}

func TestNewWithOptions(t *testing.T) {
	r, err := newTestResource(resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	uid, err := uuid.NewFromString(objUID)
	if err != nil {
		t.Errorf("failed to create new uid: %v", err)
	}

	o, err := New(objName, objNs, r, WithUID(uid))
	if err != nil {
		t.Fatalf("failed creating new object: %v", err)
	}

	if o.UID().Value() != objUID {
		t.Errorf("expected object uid: %s, got: %s", objUID, o.UID().Value())
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	o, err = New(objName, objNs, r, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new object: %v", err)
	}

	if val := o.Attrs().Get(k); val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}

func TestObjectLink(t *testing.T) {
	r, err := newTestResource(resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed creating test resource: %v", err)
	}

	o, err := New(objName, objNs, r)
	if err != nil {
		t.Fatalf("failed creating new object: %v", err)
	}

	to, err := uuid.NewFromString(toUID)
	if err != nil {
		t.Fatalf("failed creating to UID: %v", err)
	}

	if err := o.Link(to); err != nil {
		t.Fatalf("failed adding link to %s: %v", to.Value(), err)
	}

	if len(o.Links()) != 1 {
		t.Errorf("expected 1 link, got: %d", len(o.Links()))
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	if err := o.Link(to, space.WithAttrs(a), space.WithMerge(true)); err != nil {
		t.Fatalf("failed adding link to %s: %v", to.Value(), err)
	}

	if len(o.Links()) != 1 {
		t.Errorf("expected 1 link, got: %d", len(o.Links()))
	}

	for _, l := range o.Links() {
		if l.To().Value() == to.Value() {
			if val := l.Attrs().Get(k); val != v {
				t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
			}
		}
	}
}
