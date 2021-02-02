package link

import (
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

func createLinkEnds() (uuid.UID, uuid.UID, error) {
	from, err := uuid.New()
	if err != nil {
		return nil, nil, err
	}

	to, err := uuid.New()
	if err != nil {
		return nil, nil, err
	}

	return from, to, nil
}

func TestNew(t *testing.T) {
	from, to, err := createLinkEnds()
	if err != nil {
		t.Fatalf("failed created link ends: %v", err)
	}

	l, err := New(from, to)
	if err != nil {
		t.Errorf("failed to create new link: %v", err)
	}

	if l.To().Value() != to.Value() {
		t.Errorf("expeted to uid: %v, got: %v", to.Value(), l.To().Value())
	}

	if l.From().Value() != from.Value() {
		t.Errorf("expeted from uid: %v, got: %v", from.Value(), l.From().Value())
	}

	if c := len(l.Attrs().Keys()); c != 0 {
		t.Errorf("expected 0 attrs, got: %d", c)
	}
}

func TestNewWithOptions(t *testing.T) {
	from, to, err := createLinkEnds()
	if err != nil {
		t.Fatalf("failed created link ends: %v", err)
	}

	linkUID := "fooUID"
	luid, err := uuid.NewFromString(linkUID)
	if err != nil {
		t.Errorf("failed to create new uid: %v", err)
	}

	l, err := New(from, to, WithUID(luid))
	if err != nil {
		t.Errorf("failed to create new link: %v", err)
	}

	if l.UID().Value() != linkUID {
		t.Errorf("expected link uid: %s, got: %s", linkUID, l.UID().Value())
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create new attrs: %v", err)
	}
	k, v := "foo", "bar"
	a.Set(k, v)

	l, err = New(from, to, WithAttrs(a))
	if err != nil {
		t.Errorf("failed to create new link: %v", err)
	}

	if val := l.Attrs().Get(k); val != v {
		t.Errorf("expected attrs val: %s, for key: %s, got: %s", v, k, val)
	}
}
