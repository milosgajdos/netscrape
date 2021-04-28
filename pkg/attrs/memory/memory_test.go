package memory

import (
	"context"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"gonum.org/v1/gonum/graph/encoding"
)

func MustSet(ctx context.Context, a attrs.Attrs, k, v string, t *testing.T) {
	if err := a.Set(ctx, k, v); err != nil {
		t.Fatalf("failed to set val %s for key %s: %v", k, v, err)
	}
}

func TestAttributes(t *testing.T) {
	a := New()
	exp := 0

	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	keys, err := a.Keys(context.Background())
	if err != nil {
		t.Fatalf("failed getting attr keys: %v", err)
	}
	if count := len(keys); count != exp {
		t.Errorf("expected %d keys, got: %d", exp, count)
	}
}

func TestNewCopyFrom(t *testing.T) {
	a := New()

	k, v := "testKey", "testVal"
	MustSet(context.Background(), a, k, v, t)

	a2, err := NewCopyFrom(context.Background(), a)
	if err != nil {
		t.Fatalf("failed copying attrs from map: %v", err)
	}

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
		val, err := a.Get(context.Background(), k)
		if err != nil {
			t.Fatalf("failed to get %s: %v", k, err)
		}

		if val != tv {
			t.Errorf("expected %s for key %s, got: %s", v, k, val)
		}
	}
}

func TestGetAttribute(t *testing.T) {
	a := New()

	exp := ""
	k := "foo"

	val, err := a.Get(context.Background(), k)
	if err != nil {
		t.Fatalf("failed to get %s: %v", k, err)
	}

	if val != exp {
		t.Errorf("expected %s for key %s, got: %s", exp, k, val)
	}
}

func TestSetAttribute(t *testing.T) {
	a := New()

	attr := encoding.Attribute{
		Key:   "foo",
		Value: "bar",
	}

	if err := a.Set(context.Background(), attr.Key, attr.Value); err != nil {
		t.Fatalf("failed setting val %s, for key %s: %v", attr.Value, attr.Key, err)
	}

	val, err := a.Get(context.Background(), attr.Key)
	if err != nil {
		t.Fatalf("failed to get %s: %v", attr.Key, err)
	}
	if val != attr.Value {
		t.Errorf("expected: %s, got: %s", attr.Value, val)
	}

	exp := 1

	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	keys, err := a.Keys(context.Background())
	if err != nil {
		t.Fatalf("failed to get keys: %v", err)
	}

	if count := len(keys); count != exp {
		t.Errorf("expected %d keys, got: %d", exp, count)
	}
}
