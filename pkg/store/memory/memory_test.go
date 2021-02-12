package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
)

func TestNew(t *testing.T) {
	m, err := New(nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if _, err = m.Graph(context.TODO()); err != nil {
		t.Fatalf("failed to get graph handle: %v", err)
	}
}

func TestAddDelete(t *testing.T) {
	m, err := New(nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	node1UID := "foo1UID"
	node1Name := "foo1Name"

	e1, err := newTestEntity(node1UID, node1Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", node1UID, err)
	}

	if err := m.Add(context.TODO(), e1); err != nil {
		t.Errorf("failed storing node %s: %v", e1.UID(), err)
	}

	node2UID := "foo2UID"
	node2Name := "foo2Name"

	e2, err := newTestEntity(node2UID, node2Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", node1UID, err)
	}

	if err := m.Add(context.TODO(), e2); err != nil {
		t.Errorf("failed storing node %s: %v", e2.UID(), err)
	}

	g, err := m.Graph(context.TODO())
	if err != nil {
		t.Fatalf("failed to get graph handle: %v", err)
	}

	nodes, err := g.Nodes(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount := 2
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	if err := m.Delete(context.TODO(), e2); err != nil {
		t.Errorf("failed deleting node %s: %v", e2.UID(), err)
	}

	nodes, err = g.Nodes(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount = 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}
}

func TestLink(t *testing.T) {
	m, err := New(nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	node1UID := "foo1UID"
	node1Name := "foo1Name"

	e1, err := newTestEntity(node1UID, node1Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", node1UID, err)
	}

	if err := m.Add(context.TODO(), e1); err != nil {
		t.Errorf("failed storing node %s: %v", e1.UID(), err)
	}

	node2UID := "foo2UID"
	node2Name := "foo2Name"

	e2, err := newTestEntity(node2UID, node2Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", node1UID, err)
	}

	if err := m.Add(context.TODO(), e2); err != nil {
		t.Fatalf("failed storing node %s: %v", e2.UID(), err)
	}

	if err := m.Link(context.TODO(), e1.UID(), e2.UID()); err != nil {
		t.Errorf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
	}

	g, err := m.Graph(context.TODO())
	if err != nil {
		t.Fatalf("failed to get graph handle: %v", err)
	}

	links, err := g.Edges(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store links: %v", err)
	}

	expCount := 1
	if linkCount := len(links); linkCount != expCount {
		t.Errorf("expected links: %d, got: %d", expCount, linkCount)
	}

	if err := m.Unlink(context.TODO(), e1.UID(), e2.UID()); err != nil {
		t.Errorf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
	}

	links, err = g.Edges(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store links: %v", err)
	}

	expCount = 0
	if linkCount := len(links); linkCount != expCount {
		t.Errorf("expected links: %d, got: %d", expCount, linkCount)
	}
}

func TestQuery(t *testing.T) {
	m, err := New(nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	nodeUID := "foo1UID"
	nodeName := "foo1Name"

	e, err := newTestEntity(nodeUID, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", nodeUID, err)
	}

	if err := m.Add(context.TODO(), e); err != nil {
		t.Errorf("failed storing node %s: %v", e.UID(), err)
	}

	q := base.Build().Add(predicate.Entity(query.Node))

	qnodes, err := m.Query(context.TODO(), q)
	if err != nil {
		t.Errorf("failed to query store: %v", err)
	}

	expCount := 1
	if nodeCount := len(qnodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	q = base.Build().Add(predicate.Entity(query.EntityVal(10000)), base.IsAnyFunc)

	if _, err := m.Query(context.TODO(), q); !errors.Is(err, graph.ErrUnknownEntity) {
		t.Errorf("expected: %v, got: %v", graph.ErrUnknownEntity, err)
	}
}
