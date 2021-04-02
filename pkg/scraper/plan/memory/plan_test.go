package memory

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/scraper/plan"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	resPath = "../testdata/resources.yaml"
)

func MustNewPlan(src string, t *testing.T) *Plan {
	p, err := NewMock(src)
	if err != nil {
		t.Fatalf("failed to create mock Plan: %v", err)
	}
	return p
}

func MustTestResource(t, n, g, v, k string, test *testing.T) space.Resource {
	r, err := resource.New(t, n, g, v, k, false)
	if err != nil {
		test.Fatalf("failed to create resource: %v", err)
	}
	return r
}

func TestOrigin(t *testing.T) {
	src := "file:///" + resPath
	p := MustNewPlan(src, t)

	o, err := p.Origin(context.Background())
	if err != nil {
		t.Fatalf("failed to get space origin: %v", err)
	}

	if !strings.EqualFold(src, o.URL().String()) {
		t.Errorf("expected: %s, got: %s", src, o.URL().String())
	}
}

func TestAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		src := "file:///" + resPath
		p := MustNewPlan(src, t)
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
		src := "file:///" + resPath
		p := MustNewPlan(src, t)

		rx, err := p.GetAll(context.Background())
		if err != nil {
			t.Fatalf("failed getting all resource: %v", err)
		}

		exp := 12
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
		src := "file:///" + resPath
		p := MustNewPlan(src, t)
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
		src := "file:///" + resPath
		p := MustNewPlan(src, t)

		uid, err := uuid.New()
		if err != nil {
			t.Fatalf("failed to generate uid: %v", err)
		}
		if _, err := p.Get(context.Background(), uid); !errors.Is(err, plan.ErrResourceNotFound) {
			t.Errorf("expected error: %v, got: %v", plan.ErrResourceNotFound, err)
		}
	})
}

func TestDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		src := "file:///" + resPath
		p := MustNewPlan(src, t)
		r := MustTestResource("fooType", "fooName", "fooGroup", "fooVersion", "fooKind", t)

		if err := p.Add(context.Background(), r); err != nil {
			t.Fatalf("failed adding resource %s: %v", r.UID(), err)
		}

		if err := p.Delete(context.Background(), r.UID()); err != nil {
			t.Fatalf("failed removing resource %s: %v", r.UID(), err)
		}

		if _, err := p.Get(context.Background(), r.UID()); !errors.Is(err, plan.ErrResourceNotFound) {
			t.Errorf("expected %v: got: %v", plan.ErrResourceNotFound, err)
		}
	})
}
