package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/graph"
	gm "github.com/milosgajdos/netscrape/pkg/graph/memory"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/store"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	testID = "testID"
)

func TestNew(t *testing.T) {
	m, err := NewStore(testID, nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if id := m.ID(); id != testID {
		t.Fatalf("expected id: %s, got: %s", id, testID)
	}

	if _, err = m.Graph(context.TODO()); err != nil {
		t.Fatalf("failed to get graph handle: %v", err)
	}
}

func TestAddDelete(t *testing.T) {
	m, err := NewStore(testID, nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	node1ID := 1
	node1UID := "foo1UID"
	node1Name := "foo1Name"

	o, err := newTestObject(node1UID, node1Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object %q: %v", node1UID, err)
	}

	n1, err := gm.NewNode(int64(node1ID), o, graph.NodeOptions{})
	if err != nil {
		t.Fatalf("failed creating new node: %v", err)
	}

	if err := m.Add(context.TODO(), n1, store.AddOptions{}); err != nil {
		t.Errorf("failed storing node %s: %v", n1.UID(), err)
	}

	node2ID := 2
	node2UID := "foo2UID"
	node2Name := "foo2Name"

	o2, err := newTestObject(node2UID, node2Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object %q: %v", node1UID, err)
	}

	n2, err := gm.NewNode(int64(node2ID), o2, graph.NodeOptions{})
	if err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	if err := m.Add(context.TODO(), n2, store.AddOptions{}); err != nil {
		t.Errorf("failed storing node %s: %v", n2.UID(), err)
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

	uid, err := uuid.NewFromString("garbage")
	if err != nil {
		t.Fatalf("error creating new uid: %v", err)
	}

	entX, err := entity.NewWithUID(uid)
	if err != nil {
		t.Fatalf("failed creating entity: %v", err)
	}

	if err := m.Add(context.TODO(), entX, store.AddOptions{}); !errors.Is(err, store.ErrUnknownEntity) {
		t.Errorf("expected: %v, got: %v", store.ErrUnknownEntity, err)
	}

	edge, err := gm.NewEdge(n1, n2, graph.DefaultWeight)
	if err != nil {
		t.Errorf("failed creating edge: %v", err)
	}

	if err := m.Add(context.TODO(), edge, store.AddOptions{}); err != nil {
		t.Errorf("failed storing edge %s: %v", edge.UID(), err)
	}

	edges, err := g.Edges(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store edges: %v", err)
	}

	expCount = 1
	if edgeCount := len(edges); edgeCount != expCount {
		t.Errorf("expected edges: %d, got: %d", expCount, edgeCount)
	}

	if err := m.Delete(context.TODO(), edge, store.DelOptions{}); err != nil {
		t.Errorf("failed deleting edge %s: %v", edge.UID(), err)
	}

	edges, err = g.Edges(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store edges: %v", err)
	}

	expCount = 0
	if edgeCount := len(edges); edgeCount != expCount {
		t.Errorf("expected edges: %d, got: %d", expCount, edgeCount)
	}

	if err := m.Delete(context.TODO(), n2, store.DelOptions{}); err != nil {
		t.Errorf("failed storing node %s: %v", n2.UID(), err)
	}

	nodes, err = g.Nodes(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount = 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	if err := m.Delete(context.TODO(), entX, store.DelOptions{}); !errors.Is(err, store.ErrUnknownEntity) {
		t.Errorf("expected: %v, got: %v", store.ErrUnknownEntity, err)
	}
}

func TestQuery(t *testing.T) {
	m, err := NewStore(testID, nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	node1ID := 1
	node1UID := "foo1UID"
	node1Name := "foo1Name"

	o, err := newTestObject(node1UID, node1Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object %q: %v", node1UID, err)
	}

	n1, err := gm.NewNode(int64(node1ID), o, graph.NodeOptions{})
	if err != nil {
		t.Fatalf("failed creating new node: %v", err)
	}

	if err := m.Add(context.TODO(), n1, store.AddOptions{}); err != nil {
		t.Errorf("failed storing node %s: %v", n1.UID(), err)
	}

	q := base.Build().Add(query.Entity(query.Node))

	qnodes, err := m.Query(context.TODO(), q)
	if err != nil {
		t.Errorf("failed to query nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(qnodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	q = base.Build().Add(query.Entity(query.EntityVal(10000)), query.IsAnyFunc)

	if _, err := m.Query(context.TODO(), q); !errors.Is(err, graph.ErrUnknownEntity) {
		t.Errorf("expected: %v, got: %v", graph.ErrUnknownEntity, err)
	}
}
