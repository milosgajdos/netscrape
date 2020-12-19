package memory

import (
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// Store is in-memory store.
type Store struct {
	// id is store id
	id string
	// g is the store graph
	g memory.Graph
}

// NewStore creates a new in-memory store and returns it.
// If g is nil, the store creates *memory.WUG with default options.
func NewStore(id string, g memory.Graph) (*Store, error) {
	if g == nil {
		var err error
		g, err = memory.NewWUG(id+"-graph", graph.Options{})
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		id: id,
		g:  g,
	}, nil
}

// ID returns store ID.
func (m Store) ID() string {
	return m.id
}

// Graph returns graph handle.
func (m *Store) Graph() graph.Graph {
	return m.g
}

// Add stores e in memory store.
func (m *Store) Add(e store.Entity, opts store.AddOptions) error {
	switch v := e.(type) {
	case graph.Node:
		return m.g.AddNode(v)
	case graph.Edge:
		from := v.FromNode().UID()
		to := v.ToNode().UID()

		if _, err := m.g.Link(from, to, graph.LinkOptions{Attrs: opts.Attrs}); err != nil {
			return err
		}
		return nil
	}

	return store.ErrUnknownEntity
}

// Delete deletes e from memory store.
func (m *Store) Delete(e store.Entity, opts store.DelOptions) error {
	switch v := e.(type) {
	case graph.Node:
		return m.g.RemoveNode(v.UID())
	case graph.Edge:
		return m.g.RemoveLink(v.FromNode().UID(), v.ToNode().UID())
	}

	return store.ErrUnknownEntity
}

// Query queries the store and returns the results
func (m Store) Query(q query.Query) ([]store.Entity, error) {
	g, ok := m.g.(graph.Querier)
	if !ok {
		return nil, fmt.Errorf("graph error: %w", graph.ErrUnsupported)
	}

	qents, err := g.Query(q)
	if err != nil {
		return nil, err
	}

	results := make([]store.Entity, len(qents))

	for i, e := range qents {
		results[i] = e.(store.Entity)
	}

	return results, nil
}