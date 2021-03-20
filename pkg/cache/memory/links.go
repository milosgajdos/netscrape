package memory

import (
	"context"
	"sync"

	"github.com/milosgajdos/netscrape/pkg/cache"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Links in is-memory space.Link cache.
type Links struct {
	// from indexes from->to links.
	from map[uuid.UID]map[uuid.UID]space.Link
	// to indexes to<-from links.
	to map[uuid.UID]map[uuid.UID]space.Link
	// mu synchronizes access to Links
	mu *sync.RWMutex
}

// NewLinks creates a new Links cache and returns it
func NewLinks() (*Links, error) {
	return &Links{
		from: make(map[uuid.UID]map[uuid.UID]space.Link),
		to:   make(map[uuid.UID]map[uuid.UID]space.Link),
		mu:   &sync.RWMutex{},
	}, nil
}

func (c *Links) put(ctx context.Context, link space.Link, opts ...cache.Option) error {
	copts := cache.Options{}
	for _, apply := range opts {
		apply(&copts)
	}

	f, t := link.From(), link.To()

	if c.from[f][t] == nil {
		c.from[f] = make(map[uuid.UID]space.Link)
	}

	if c.to[t][f] == nil {
		c.to[t] = make(map[uuid.UID]space.Link)
	}

	if copts.Upsert {
		c.from[f][t] = link
		c.to[t][f] = link
		return nil
	}

	if _, ok := c.from[f][t]; !ok {
		c.from[f][t] = link
	}

	if _, ok := c.to[t][f]; !ok {
		c.to[t][f] = link
	}
	return nil
}

// Put stores link in the cache.
func (c *Links) Put(ctx context.Context, link space.Link, opts ...cache.Option) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.put(ctx, link, opts...)
}

func (c *Links) get(ctx context.Context, uid uuid.UID, index map[uuid.UID]space.Link, opts ...cache.Option) ([]space.Link, error) {
	lx := make([]space.Link, len(index))

	i := 0
	for _, l := range index {
		lx[i] = l
		i++
	}

	return lx, nil
}

// GetFrom returns all links from the given uid.
func (c *Links) GetFrom(ctx context.Context, uid uuid.UID, opts ...cache.Option) ([]space.Link, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if from, ok := c.from[uid]; ok {
		return c.get(ctx, uid, from, opts...)
	}

	return []space.Link{}, nil
}

// GetTo returns all link to the given uid.
func (c *Links) GetTo(ctx context.Context, uid uuid.UID, opts ...cache.Option) ([]space.Link, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if to, ok := c.to[uid]; ok {
		return c.get(ctx, uid, to, opts...)
	}

	return []space.Link{}, nil
}

func (c *Links) delete(ctx context.Context, uid uuid.UID, opts ...cache.Option) error {
	delete(c.to, uid)
	delete(c.from, uid)
	return nil
}

// Delete removes all links to or from the given UID.
func (c *Links) Delete(ctx context.Context, uid uuid.UID, opts ...cache.Option) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.delete(ctx, uid, opts...)
}

// Clear clears cache.
func (c *Links) Clear(ctx context.Context, opts ...cache.Option) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.from = make(map[uuid.UID]map[uuid.UID]space.Link)
	c.to = make(map[uuid.UID]map[uuid.UID]space.Link)

	return nil
}

// BulkPut puts all links key-ed by UID into cache.
func (c *Links) BulkPut(ctx context.Context, links []space.Link, opts ...cache.Option) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, l := range links {
		if err := c.put(ctx, l, opts...); err != nil {
			return err
		}
	}
	return nil
}

// BulkGetFrom returns all links from the given UIDs.
func (c *Links) BulkGetFrom(ctx context.Context, uids []uuid.UID, opts ...cache.Option) (map[uuid.UID][]space.Link, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m := make(map[uuid.UID][]space.Link)

	for _, uid := range uids {
		lx := []space.Link{}
		from, ok := c.from[uid]
		if !ok {
			continue
		}
		links, err := c.get(ctx, uid, from, opts...)
		if err != nil {
			return nil, err
		}
		lx = append(lx, links...)
		m[uid] = lx
	}
	return m, nil
}

// BulkGetTo returns all links to the given UIDs.
func (c *Links) BulkGetTo(ctx context.Context, uids []uuid.UID, opts ...cache.Option) (map[uuid.UID][]space.Link, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m := make(map[uuid.UID][]space.Link)

	for _, uid := range uids {
		lx := []space.Link{}
		to, ok := c.to[uid]
		if !ok {
			continue
		}
		links, err := c.get(ctx, uid, to, opts...)
		if err != nil {
			return nil, err
		}
		lx = append(lx, links...)
		m[uid] = lx
	}
	return m, nil
}

// BulkDelete removes Entities with given uid from topology.
func (c *Links) BulkDelete(ctx context.Context, uids []uuid.UID, opts ...cache.Option) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, uid := range uids {
		if err := c.delete(ctx, uid, opts...); err != nil {
			return err
		}
	}
	return nil
}
