package memory

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

func TestNew(t *testing.T) {
	m, err := NewStore()
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if _, err = m.Graph(); err != nil {
		t.Fatalf("failed to get graph handle: %v", err)
	}
}

func TestAddGetDelete(t *testing.T) {
	m, err := NewStore()
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

	res, err := m.Get(context.Background(), e1.UID())
	if err != nil {
		t.Errorf("failed getting node %s: %v", e1.UID(), err)
	}

	if !reflect.DeepEqual(res.UID(), e1.UID()) {
		t.Errorf("expected entity with uid: %s, got: %s", e1.UID(), res.UID())
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

	g, err := m.Graph()
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

	if err := m.Delete(context.TODO(), e2.UID()); err != nil {
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

	uid, err := uuid.New()
	if err != nil {
		t.Fatalf("failed to generate uid: %v", err)
	}

	if _, err := m.Get(context.Background(), uid); !errors.Is(err, store.ErrEntityNotFound) {
		t.Errorf("expected error: %v, got: %v", store.ErrEntityNotFound, err)
	}
}

func TestLink(t *testing.T) {
	m, err := NewStore()
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

	g, err := m.Graph()
	if err != nil {
		t.Fatalf("failed to get graph handle: %v", err)
	}

	nodes, err := g.From(context.TODO(), e1.UID())
	if err != nil {
		t.Fatalf("failed to get store links: %v", err)
	}

	expCount := 1
	if count := len(nodes); count != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, count)
	}

	if err := m.Unlink(context.TODO(), e1.UID(), e2.UID()); err != nil {
		t.Errorf("failed linking %v to %v: %v", e1.UID(), e2.UID(), err)
	}

	nodes, err = g.From(context.TODO(), e1.UID())
	if err != nil {
		t.Fatalf("failed to get store links: %v", err)
	}

	expCount = 0
	if count := len(nodes); count != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, count)
	}
}
