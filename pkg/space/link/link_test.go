package link

import (
	"context"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func MustSet(ctx context.Context, a attrs.Attrs, k, v string, t *testing.T) {
	if err := a.Set(ctx, k, v); err != nil {
		t.Fatalf("failed to set val %s for key %s: %v", k, v, err)
	}
}

func TestNew(t *testing.T) {
	from, to := memuid.New(), memuid.New()

	l, err := New(from, to)
	if err != nil {
		t.Errorf("failed to create new link: %v", err)
	}

	if l.To().String() != to.String() {
		t.Errorf("expeted to uid: %v, got: %v", to.String(), l.To().String())
	}

	if l.From().String() != from.String() {
		t.Errorf("expeted from uid: %v, got: %v", from.String(), l.From().String())
	}

	keys, err := l.Attrs().Keys(context.Background())
	if err != nil {
		t.Fatalf("failed to get keys: %v", err)
	}

	if c := len(keys); c != 0 {
		t.Errorf("expected 0 attrs, got: %d", c)
	}
}

func TestNewWithOptions(t *testing.T) {
	from, to := memuid.New(), memuid.New()

	linkUID := "fooUID"
	luid := memuid.NewFromString(linkUID)

	l, err := New(from, to, WithUID(luid))
	if err != nil {
		t.Errorf("failed to create new link: %v", err)
	}

	if l.UID().String() != linkUID {
		t.Errorf("expected link uid: %s, got: %s", linkUID, l.UID().String())
	}

	a := memattrs.New()
	k, v := "foo", "bar"
	MustSet(context.Background(), a, k, v, t)

	l, err = New(from, to, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new link: %v", err)
	}

	val, err := l.Attrs().Get(context.Background(), k)
	if err != nil {
		t.Fatalf("failed to get val for key %s: %v", k, err)
	}

	if val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
