package memory

import (
	"context"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/cache"
	"github.com/milosgajdos/netscrape/pkg/space/link"
	"github.com/milosgajdos/netscrape/pkg/uuid"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

func MustNewLinksCache(t *testing.T) *Links {
	c, err := NewLinks()
	if err != nil {
		t.Fatalf("failed to create links cache: %v", err)
	}
	return c
}

func MustLink(from, to uuid.UID, t *testing.T) cache.Link {
	l, err := link.New(from, to)
	if err != nil {
		t.Fatalf("failed to create link from %s to %s : %v", from, to, err)
	}
	return l
}

func MustLinks(count int, t *testing.T) []cache.Link {
	links := make([]cache.Link, count)
	for i := 0; i < count; i++ {
		from := memuid.New()
		to := memuid.New()
		links[i] = MustLink(from, to, t)
	}
	return links
}

func TestPut(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		c := MustNewLinksCache(t)
		from := memuid.New()
		to := memuid.New()
		l := MustLink(from, to, t)

		if err := c.Put(context.Background(), l); err != nil {
			t.Errorf("failed adding link: %v", err)
		}

		// adding the same link without Upsert returns nil
		if err := c.Put(context.Background(), l); err != nil {
			t.Errorf("failed adding link: %v", err)
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		c := MustNewLinksCache(t)
		from := memuid.New()
		to := memuid.New()
		l := MustLink(from, to, t)

		if err := c.Put(context.Background(), l); err != nil {
			t.Errorf("failed adding link: %v", err)
		}

		if err := c.Put(context.Background(), l, cache.WithUpsert()); err != nil {
			t.Errorf("failed adding link: %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("From", func(t *testing.T) {
		c := MustNewLinksCache(t)
		from := memuid.New()
		to := memuid.New()
		l := MustLink(from, to, t)

		if err := c.Put(context.Background(), l); err != nil {
			t.Fatalf("failed adding link: %v", err)
		}

		links, err := c.GetFrom(context.Background(), from)
		if err != nil {
			t.Fatalf("failed getting links from %s: %v", from, err)
		}

		exp := 1
		if cl := len(links); cl != exp {
			t.Errorf("expected links: %d, got: %d", exp, cl)
		}

		if f := links[0].From(); !reflect.DeepEqual(f, from) {
			t.Errorf("expected uid: %s, got: %s", from, f)
		}
	})

	t.Run("To", func(t *testing.T) {
		c := MustNewLinksCache(t)
		from := memuid.New()
		to := memuid.New()
		l := MustLink(from, to, t)

		if err := c.Put(context.Background(), l); err != nil {
			t.Fatalf("failed adding link: %v", err)
		}

		links, err := c.GetTo(context.Background(), to)
		if err != nil {
			t.Fatalf("failed getting links from %s: %v", from, err)
		}

		exp := 1
		if e := len(links); e != exp {
			t.Errorf("expected links: %d, got: %d", exp, e)
		}

		if tl := links[0].To(); !reflect.DeepEqual(tl, to) {
			t.Errorf("expected uid: %s, got: %s", to, tl)
		}
	})
}

func TestDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		c := MustNewLinksCache(t)
		from := memuid.New()
		to := memuid.New()
		l := MustLink(from, to, t)

		if err := c.Put(context.Background(), l); err != nil {
			t.Fatalf("failed adding link: %v", err)
		}

		if err := c.Delete(context.Background(), from); err != nil {
			t.Fatalf("failed deleting links for %s: %v", from, err)
		}

		links, err := c.GetFrom(context.Background(), from)
		if err != nil {
			t.Errorf("unexpected error getting link for %v: %v", from, err)
		}

		exp := 0
		if e := len(links); e != exp {
			t.Errorf("expected links: %d, got: %d", exp, e)
		}
	})
}

func TestClear(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("From", func(t *testing.T) {
		c := MustNewLinksCache(t)
		from := memuid.New()
		to := memuid.New()
		l := MustLink(from, to, t)

		if err := c.Put(context.Background(), l); err != nil {
			t.Fatalf("failed adding link: %v", err)
		}

		if err := c.Clear(context.Background()); err != nil {
			t.Fatalf("failed to celar cache: %v", err)
		}

		links, err := c.GetFrom(context.Background(), from)
		if err != nil {
			t.Errorf("unexpected error getting link for %v: %v", from, err)
		}

		exp := 0
		if e := len(links); e != exp {
			t.Errorf("expected links: %d, got: %d", exp, e)
		}

		links, err = c.GetTo(context.Background(), to)
		if err != nil {
			t.Errorf("unexpected error getting link for %v: %v", from, err)
		}

		if e := len(links); e != exp {
			t.Errorf("expected links: %d, got: %d", exp, e)
		}
	})
}

func TestBulkAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		c := MustNewLinksCache(t)
		links := MustLinks(5, t)

		if err := c.BulkPut(context.Background(), links); err != nil {
			t.Errorf("failed storing links in cache: %v", err)
		}
	})
}

func TestBulkGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("From", func(t *testing.T) {
		c := MustNewLinksCache(t)
		links := MustLinks(5, t)

		if err := c.BulkPut(context.Background(), links); err != nil {
			t.Errorf("failed storing links in cache: %v", err)
		}

		var uids []uuid.UID

		for _, l := range links {
			uids = append(uids, l.From(), l.To())
		}

		m, err := c.BulkGetFrom(context.Background(), uids)
		if err != nil {
			t.Fatalf("failed getting the bulk lof links: %v", err)
		}

		for uid, lx := range m {
			links, err := c.GetFrom(context.Background(), uid)
			if err != nil {
				t.Errorf("unexpected error getting link for %v: %v", uid, err)
			}

			if len(lx) != len(links) {
				t.Errorf("expected links: %d, got: %d", len(lx), len(links))
			}
		}
	})

	t.Run("To", func(t *testing.T) {
		c := MustNewLinksCache(t)
		links := MustLinks(5, t)

		if err := c.BulkPut(context.Background(), links); err != nil {
			t.Errorf("failed storing links in cache: %v", err)
		}

		var uids []uuid.UID

		for _, l := range links {
			uids = append(uids, l.From(), l.To())
		}

		m, err := c.BulkGetTo(context.Background(), uids)
		if err != nil {
			t.Fatalf("failed getting the bulk lof links: %v", err)
		}

		for uid, lx := range m {
			links, err := c.GetTo(context.Background(), uid)
			if err != nil {
				t.Errorf("unexpected error getting link for %v: %v", uid, err)
			}

			if len(lx) != len(links) {
				t.Errorf("expected links: %d, got: %d", len(lx), len(links))
			}
		}
	})
}

func TestBulkDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		c := MustNewLinksCache(t)
		links := MustLinks(5, t)

		if err := c.BulkPut(context.Background(), links); err != nil {
			t.Errorf("failed storing links in cache: %v", err)
		}

		var uids []uuid.UID

		for _, l := range links {
			uids = append(uids, l.From(), l.To())
		}

		if err := c.BulkDelete(context.Background(), uids); err != nil {
			t.Fatalf("failed deleting entities: %v", err)
		}

		for _, uid := range uids {
			links, err := c.GetFrom(context.Background(), uid)
			if err != nil {
				t.Errorf("unexpected error getting link for %v: %v", uid, err)
			}

			exp := 0
			if e := len(links); e != exp {
				t.Errorf("expected links: %d, got: %d", exp, e)
			}

			links, err = c.GetTo(context.Background(), uid)
			if err != nil {
				t.Errorf("unexpected error getting link for %v: %v", uid, err)
			}

			if e := len(links); e != exp {
				t.Errorf("expected links: %d, got: %d", exp, e)
			}
		}
	})
}
