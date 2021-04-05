package simple

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/internal"
	"github.com/milosgajdos/netscrape/pkg/scraper"
	"github.com/milosgajdos/netscrape/pkg/space"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func MustNewSimple(t *testing.T) *Simple {
	p, err := NewSimple()
	if err != nil {
		t.Fatalf("failed to create mock Plan: %v", err)
	}
	return p
}

func MustTestResource(t, n, g, v, k string, test *testing.T) space.Resource {
	r, err := internal.NewTestResource(t, n, g, v, k, false)
	if err != nil {
		test.Fatalf("failed to create resource: %v", err)
	}
	return r
}

func TestAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewSimple(t)
		r := MustTestResource("fooType", "fooName", "fooGroup", "fooVersion", "fooKind", t)

		if err := p.Add(context.Background(), r); err != nil {
			t.Errorf("failed adding resource %s: %v", r.UID(), err)
		}
	})
}

func TestGetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewSimple(t)

		rx, err := p.GetAll(context.Background())
		if err != nil {
			t.Fatalf("failed getting all resource: %v", err)
		}

		exp := 0
		if c := len(rx); c != exp {
			t.Errorf("expected resources: %d, got: %d", exp, c)
		}
	})
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewSimple(t)
		r := MustTestResource("fooType", "fooName", "fooGroup", "fooVersion", "fooKind", t)

		if err := p.Add(context.Background(), r); err != nil {
			t.Fatalf("failed adding resource %s: %v", r.UID(), err)
		}

		res, err := p.Get(context.Background(), r.UID())
		if err != nil {
			t.Fatalf("failed getting resource %s: %v", r.UID(), err)
		}

		if !reflect.DeepEqual(res.UID(), r.UID()) {
			t.Errorf("expected entity: %s, got: %s", r.UID(), res.UID())
		}
	})

	t.Run("ErrResourceNotFound", func(t *testing.T) {
		p := MustNewSimple(t)

		if _, err := p.Get(context.Background(), memuid.New()); !errors.Is(err, scraper.ErrResourceNotFound) {
			t.Errorf("expected error: %v, got: %v", scraper.ErrResourceNotFound, err)
		}
	})
}

func TestDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		p := MustNewSimple(t)
		r := MustTestResource("fooType", "fooName", "fooGroup", "fooVersion", "fooKind", t)

		if err := p.Add(context.Background(), r); err != nil {
			t.Fatalf("failed adding resource %s: %v", r.UID(), err)
		}

		if err := p.Delete(context.Background(), r.UID()); err != nil {
			t.Fatalf("failed removing resource %s: %v", r.UID(), err)
		}

		if _, err := p.Get(context.Background(), r.UID()); !errors.Is(err, scraper.ErrResourceNotFound) {
			t.Errorf("expected %v: got: %v", scraper.ErrResourceNotFound, err)
		}
	})
}
