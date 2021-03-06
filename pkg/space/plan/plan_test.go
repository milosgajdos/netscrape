package plan

import (
	"context"
	"strings"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
)

const (
	resPath = "../testdata/undirected/resources.yaml"
)

func TestSource(t *testing.T) {
	src := "file:///" + resPath

	space, err := NewMock(src)
	if err != nil {
		t.Fatalf("failed to create mock Plan: %v", err)
	}

	s, err := space.Origin(context.Background())
	if err != nil {
		t.Fatalf("failed to get space origin: %v", err)
	}

	if !strings.EqualFold(src, s.URL().String()) {
		t.Errorf("expected: %s, got: %s", src, s.URL().String())
	}
}

func TestResources(t *testing.T) {
	src := "file:///" + resPath

	space, err := NewMock(src)
	if err != nil {
		t.Fatalf("failed to create mock Plan: %v", err)
	}

	resources, err := space.Resources(context.Background())
	if err != nil {
		t.Fatalf("failed to get plan resources: %v", err)
	}

	if len(resources) == 0 {
		t.Errorf("no resources found")
	}
}

func TestPlanGet(t *testing.T) {
	src := "file:///" + resPath

	plan, err := NewMock(src)
	if err != nil {
		t.Errorf("failed to create mock Plan: %v", err)
		return
	}

	resources, err := plan.Resources(context.Background())
	if err != nil {
		t.Fatalf("failed to get plan resources: %v", err)
	}
	count := len(resources)

	groups := make([]string, count)
	versions := make([]string, count)
	kinds := make([]string, count)
	names := make([]string, count)

	for i, r := range resources {
		groups[i] = r.Group()
		versions[i] = r.Version()
		kinds[i] = r.Kind()
		names[i] = r.Name()
	}

	for _, group := range groups {
		q := base.Build().Add(predicate.Group(group))

		resources, err := plan.Get(context.Background(), q)
		if err != nil {
			t.Errorf("error querying group %s: %v", group, err)
		}

		for _, r := range resources {
			if r.Group() != group {
				t.Errorf("expected: %s, got: %s", group, r.Group())
			}
		}

		for _, version := range versions {
			q = q.Add(predicate.Version(version))

			resources, err := plan.Get(context.Background(), q)
			if err != nil {
				t.Errorf("error querying g/v %s/%s: %v", group, version, err)
			}

			for _, res := range resources {
				if res.Version() != version || res.Group() != group {
					t.Errorf("expected: %s/%s, got: %s/%s", group, version, res.Group(), res.Version())
				}
			}

			for _, kind := range kinds {
				q = q.Add(predicate.Kind(kind))

				resources, err := plan.Get(context.Background(), q)
				if err != nil {
					t.Errorf("error querying g/v/k: %s/%s/%s: %v", group, version, kind, err)
				}

				for _, res := range resources {
					if res.Kind() != kind || res.Version() != version || res.Group() != group {
						t.Errorf("expected: %s/%s/%s, got: %s/%s/%s", group, version, kind,
							res.Group(), res.Version(), res.Kind())
					}
				}

				for _, name := range names {
					q = q.Add(predicate.Name(name))

					resources, err := plan.Get(context.Background(), q)
					if err != nil {
						t.Errorf("error querying g/v/k/n: %s/%s/%s/%s: %v", group, version, kind, name, err)
					}

					for _, res := range resources {
						if res.Name() != name || res.Kind() != kind || res.Version() != version || res.Group() != group {
							t.Errorf("expected: %s/%s/%s/%s, got: %s/%s/%s/%s", group, version, kind, name,
								res.Group(), res.Version(), res.Kind(), res.Name())
						}
					}
				}
			}
		}
	}
}

func TestPlanRemove(t *testing.T) {
	src := "file:///" + resPath

	plan, err := NewMock(src)
	if err != nil {
		t.Errorf("failed to create mock Plan: %v", err)
		return
	}

	ctx := context.Background()

	rx, err := plan.Resources(ctx)
	if err != nil {
		t.Fatalf("failed to get plan resources: %v", err)
	}

	for _, r := range rx {
		if err := plan.Remove(ctx, r); err != nil {
			t.Errorf("failed removing resource %s/%s/%s: %v", r.Group(), r.Version(), r.Kind(), err)
		}
	}

	rx, err = plan.Resources(ctx)
	if err != nil {
		t.Fatalf("failed to get plan resources: %v", err)
	}

	if count := len(rx); count != 0 {
		t.Errorf("expected %d resources, got: %d", 0, count)
	}
}
