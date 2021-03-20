package memory

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph/encoding"
)

func TestAttributes(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}

	exp := 0

	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	keys := a.Keys()
	if count := len(keys); count != exp {
		t.Errorf("expected %d keys, got: %d", exp, count)
	}
}

func TestNewCopyFrom(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}

	k, v := "testKey", "testVal"
	a.Set(k, v)

	a2 := NewCopyFrom(a)

	if !reflect.DeepEqual(a, a2) {
		t.Errorf("expected %v, got: %v", a, a2)
	}
}

func TestNewFromMap(t *testing.T) {
	tk, tv := "testKey", "testVal"
	m := map[string]string{
		tk: tv,
	}

	a := NewFromMap(m)

	for k, v := range m {
		if val := a.Get(k); val != tv {
			t.Errorf("expected %s for key %s, got: %s", v, k, val)
		}
	}
}

func TestGetAttribute(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}

	exp := ""
	if val := a.Get("foo"); val != exp {
		t.Errorf("expected: %s, got: %s", exp, val)
	}
}

func TestSetAttribute(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}

	attr := encoding.Attribute{
		Key:   "foo",
		Value: "bar",
	}

	a.Set(attr.Key, attr.Value)

	if val := a.Get(attr.Key); val != attr.Value {
		t.Errorf("expected: %s, got: %s", attr.Value, val)
	}

	exp := 1

	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	keys := a.Keys()

	if count := len(keys); count != exp {
		t.Errorf("expected %d keys, got: %d", exp, count)
	}
}
