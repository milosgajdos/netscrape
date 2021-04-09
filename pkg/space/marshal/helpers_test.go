package marshal

import (
	"testing"

	"github.com/milosgajdos/netscrape/pkg/internal"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/link"
)

func MustEntity(t *testing.T, opts ...entity.Option) space.Entity {
	e, err := internal.NewTestEntity(opts...)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	return e
}

func MustResource(t *testing.T, opts ...entity.Option) space.Resource {
	r, err := internal.NewTestResource(opts...)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}
	return r
}

func MustObject(t *testing.T, opts ...entity.Option) space.Object {
	o, err := internal.NewTestObject(opts...)
	if err != nil {
		t.Fatalf("failed to create object: %v", err)
	}
	return o
}

func MustLink(t *testing.T, opts ...link.Option) space.Link {
	l, err := internal.NewTestLink(opts...)
	if err != nil {
		t.Fatalf("failed to create link: %v", err)
	}
	return l
}
