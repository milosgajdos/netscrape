package plan

import (
	"strings"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
)

const (
	resPath = "../testdata/undirected/resources.yaml"
)

func TestSource(t *testing.T) {
	src := "file:///" + resPath

	space, err := NewMock(src)
	if err != nil {
		t.Fatalf("failed to create mock Space: %v", err)
	}

	if s := space.Origin(); !strings.EqualFold(src, s.URL().String()) {
		t.Errorf("expected: %s, got: %s", src, s.URL().String())
	}
}

func TestResources(t *testing.T) {
	src := "file:///" + resPath

	space, err := NewMock(src)
	if err != nil {
		t.Fatalf("failed to create mock Space: %v", err)
	}

	resources := space.Resources()
	if len(resources) == 0 {
		t.Errorf("no resources found")
	}
}

func TestSpaceGet(t *testing.T) {
	src := "file:///" + resPath

	space, err := NewMock(src)
	if err != nil {
		t.Errorf("failed to create mock Space: %v", err)
		return
	}

	groups := make([]string, len(space.Resources()))
	versions := make([]string, len(space.Resources()))
	kinds := make([]string, len(space.Resources()))
	names := make([]string, len(space.Resources()))

	for i, r := range space.Resources() {
		groups[i] = r.Group()
		versions[i] = r.Version()
		kinds[i] = r.Kind()
		names[i] = r.Name()
	}

	for _, group := range groups {
		q := base.Build().Add(query.Group(group), query.StringEqFunc(group))

		resources, err := space.Get(q)
		if err != nil {
			t.Errorf("error querying group %s: %v", group, err)
		}

		for _, r := range resources {
			if r.Group() != group {
				t.Errorf("expected: %s, got: %s", group, r.Group())
			}
		}

		for _, version := range versions {
			q = q.Add(query.Version(version), query.StringEqFunc(version))

			resources, err := space.Get(q)
			if err != nil {
				t.Errorf("error querying g/v %s/%s: %v", group, version, err)
			}

			for _, res := range resources {
				if res.Version() != version || res.Group() != group {
					t.Errorf("expected: %s/%s, got: %s/%s", group, version, res.Group(), res.Version())
				}
			}

			for _, kind := range kinds {
				q = q.Add(query.Kind(kind), query.StringEqFunc(kind))

				resources, err := space.Get(q)
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
					q = q.Add(query.Name(name), query.StringEqFunc(name))

					resources, err := space.Get(q)
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
