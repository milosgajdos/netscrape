package memory

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	wugEntPath = "testdata/wug/entities.yaml"
)

func TestWUGAddGetRemoveNode(t *testing.T) {
	g, err := NewWUG()
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	if uid := g.UID(); uid == nil {
		t.Errorf("expected uid, got: %v", uid)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := newTestEntity(nodeID, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n, err := g.NewNode(context.TODO(), o)
	if err != nil {
		t.Errorf("failed creating new graph node: %v", err)
	}

	// Add a new node
	if err := g.AddNode(context.TODO(), n); err != nil {
		t.Errorf("failed adding node: %v", err)
	}

	nodes, err := g.Nodes(context.TODO())
	if err != nil {
		t.Fatalf("failed getting nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	// adding the same nodes twice should not change the node count
	if err := g.AddNode(context.TODO(), n); err != nil {
		t.Errorf("failed adding node: %v", err)
	}

	expCount = 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	// Get the node with given uid
	node, err := g.Node(context.TODO(), n.UID())
	if err != nil {
		t.Errorf("failed to get %s node: %v", n.UID(), err)
	}

	if !reflect.DeepEqual(n, node) {
		t.Errorf("expected node %#v, got: %#v", node, n)
	}

	guid, err := uuid.NewFromString("garbage")
	if err != nil {
		t.Fatalf("error creating new uid: %v", err)
	}

	if _, err := g.Node(context.TODO(), guid); err != graph.ErrNodeNotFound {
		t.Errorf("expected error %v, got: %#v", graph.ErrNodeNotFound, err)
	}

	// Remove the node with given uid
	if err := g.RemoveNode(context.TODO(), n.UID()); err != nil {
		t.Errorf("failed to remove node: %v", err)
	}

	nodes, err = g.Nodes(context.TODO())
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount = 0
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	guid, err = uuid.NewFromString("garbage")
	if err != nil {
		t.Fatalf("error creating new uid: %v", err)
	}

	if err := g.RemoveNode(context.TODO(), guid); err != nil {
		t.Errorf("failed to remove node: %v", err)
	}
}

func TestWUGLinkGetRemoveEdge(t *testing.T) {
	g, err := NewWUG()
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	node1UID := "foo1UID"
	node1Name := "foo1Name"

	o1, err := newTestEntity(node1UID, node1Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", node1UID, err)
	}

	n1, err := g.NewNode(context.TODO(), o1)
	if err != nil {
		t.Errorf("failed creating new node: %v", err)
	}

	if err := g.AddNode(context.TODO(), n1); err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	node2UID := "foo2UID"
	node2Name := "foo2Name"

	o2, err := newTestEntity(node2UID, node2Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", node2UID, err)
	}

	n2, err := g.NewNode(context.TODO(), o2)
	if err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	if err := g.AddNode(context.TODO(), n2); err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	ox, err := newTestEntity("nonExUID", "nonExName", nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity %q: %v", node2UID, err)
	}

	nodeX, err := NewNode(123334444, ox)
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	// Link nodes with a node which does not exist in the graph
	if _, err := g.Link(context.TODO(), n1.UID(), nodeX.UID()); !errors.Is(err, graph.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", graph.ErrNodeNotFound, err)
	}

	if _, err := g.Link(context.TODO(), nodeX.UID(), n2.UID()); !errors.Is(err, graph.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", graph.ErrNodeNotFound, err)
	}

	edges, err := g.Edges(context.TODO())
	if err != nil {
		t.Errorf("failed getting graph edges: %v", err)
	}

	expCount := 0
	if len(edges) != expCount {
		t.Errorf("expected: %d edges, got: %d", expCount, len(edges))
	}

	edge, err := g.Link(context.TODO(), n1.UID(), n2.UID(), graph.WithWeight(graph.DefaultWeight))
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", n1.UID(), n2.UID(), err)
	}

	nodesFrom, err := g.From(context.TODO(), n1.UID())
	if err != nil {
		t.Errorf("failed to get nodes from %s: %v", n1.UID(), err)
	}

	expCount = 1
	if count := len(nodesFrom); count != expCount {
		t.Errorf("expected: %d nodes, got: %d", expCount, count)
	}

	if len(nodesFrom) == 1 {
		if nodesFrom[0].UID().Value() != n2.UID().Value() {
			t.Errorf("expected node link to %s from %s", n2.UID().Value(), n1.UID().Value())
		}
	}

	if w := edge.Weight(); big.NewFloat(w).Cmp(big.NewFloat(graph.DefaultWeight)) != 0 {
		t.Errorf("expected non-negative weight")
	}

	edges, err = g.Edges(context.TODO())
	if err != nil {
		t.Errorf("failed getting graph edges: %v", err)
	}

	expCount = 1
	if len(edges) != expCount {
		t.Errorf("no edges found in graph")
	}

	// linking already linked nodes must return the same edge/line as returned previously
	exEdge, err := g.Link(context.TODO(), n1.UID(), n2.UID())
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", n1.UID(), n2.UID(), err)
	}

	if !reflect.DeepEqual(exEdge, edge) {
		t.Errorf("expected edge %#v, got: %#v", exEdge, edge)
	}

	e, err := g.Edge(context.TODO(), n1.UID(), n2.UID())
	if err != nil {
		t.Errorf("failed getting edge between %s and %s: %v", n1.UID(), n2.UID(), err)
	}

	if !reflect.DeepEqual(e, edge) {
		t.Errorf("expected edge %#v, got: %#v", exEdge, edge)
	}

	// remove edge between previously linked nodes which are still present in the graph
	if err := g.Unlink(context.TODO(), n1.UID(), n2.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", n1.UID(), n2.UID(), err)
	}

	if _, err := g.Edge(context.TODO(), n1.UID(), n2.UID()); err != nil && !errors.Is(err, graph.ErrEdgeNotExist) {
		t.Errorf("expected error: %v, got: %v", graph.ErrEdgeNotExist, err)
	}

	// remoe edge between non-existen nodes should return nil
	if err := g.Unlink(context.TODO(), nodeX.UID(), n1.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", nodeX.UID(), n1.UID(), err)
	}

	if err := g.Unlink(context.TODO(), n1.UID(), nodeX.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", nodeX.UID(), n1.UID(), err)
	}
}

func TestWUGSubGraph(t *testing.T) {
	g, err := makeTestGraph(wugEntPath)
	if err != nil {
		t.Fatalf("failed to create new memory graph: %v", err)
	}

	guid, err := uuid.NewFromString("garbage")
	if err != nil {
		t.Fatalf("error creating new uid: %v", err)
	}

	// subgraph of non-existent node should return error
	if _, err := g.SubGraph(context.TODO(), guid, 10); err != graph.ErrNodeNotFound {
		t.Errorf("expected: %v, got: %v", graph.ErrNodeNotFound, err)
	}

	// NOTE: we are hardcoding the test value here
	// since we know that this node's neighbourhood
	suid := "fooGroup/v1/fooKind/fooNs/foo1"
	uid, err := uuid.NewFromString(suid)
	if err != nil {
		t.Fatalf("error creating new uid: %v", err)
	}

	testCases := []struct {
		depth int
		exp   int
	}{
		{0, 1},   // return node
		{1, 5},   // return node + adjacent nodes
		{100, 8}, // return all nodes reachable from node
	}

	for _, tc := range testCases {
		sg, err := g.SubGraph(context.TODO(), uid, tc.depth)
		if err != nil {
			t.Errorf("failed to get subgraph of node %s: %v", uid, err)
			continue
		}

		storeNodes, err := sg.Nodes(context.TODO())
		if err != nil {
			t.Errorf("failed to fetch subgraph nodes: %v", err)
			continue
		}

		if len(storeNodes) != tc.exp {
			t.Errorf("expected subgraph nodes: %d, got: %d", tc.exp, len(storeNodes))
		}
	}
}

func TestWUGQuery(t *testing.T) {
	g, err := makeTestGraph(wugEntPath)
	if err != nil {
		t.Fatalf("failed to create test graph: %v", err)
	}

	q := base.Build().Add(predicate.Entity(query.Node))

	qnodes, err := g.Query(context.TODO(), q)
	if err != nil {
		t.Errorf("failed to query all nodes: %v", err)
	}

	nodes, err := g.Nodes(context.TODO())
	if err != nil {
		t.Fatalf("failed to fetch nodes: %v", err)
	}

	if len(qnodes) != len(nodes) {
		t.Errorf("expected nodes: %d, got: %d", len(nodes), len(qnodes))
	}

	uids := make([]uuid.UID, len(nodes))

	for i, n := range nodes {
		uids[i] = n.UID()
	}

	q = base.Build().Add(predicate.Entity(query.Node))

	for _, uid := range uids {
		q = q.Add(predicate.UID(uid))

		nodes, err := g.Query(context.TODO(), q)
		if err != nil {
			t.Errorf("error getting node %s: %v", uid, err)
			continue
		}

		for _, node := range nodes {
			if u := node.UID().Value(); u != uid.Value() {
				t.Errorf("expected uid: %s, got: %s", uid, u)
				continue
			}
		}
	}
}

func TestWUGQueryUnknown(t *testing.T) {
	g, err := makeTestGraph(wugEntPath)
	if err != nil {
		t.Fatalf("failed to create new memory graph: %v", err)
	}

	// NOTE: EntityVal is an enum/iota starting with 0 with only two values: Node and Edge
	// Any other number higher than 1 is considered a non-existent Entity
	q := base.Build().Add(predicate.Entity(query.EntityVal(10000)), base.IsAnyFunc)

	if _, err := g.Query(context.TODO(), q); !errors.Is(err, graph.ErrUnsupported) {
		t.Errorf("expected: %v, got: %v", graph.ErrUnsupported, err)
	}
}

func TestWUGDOT(t *testing.T) {
	id := "testID"

	g, err := NewWUG(graph.WithDOTID(id))
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	if dotID := g.DOTID(); dotID != id {
		t.Errorf("expected DOTID: %s, got: %s", id, dotID)
	}

	dot, err := g.DOT()
	if err != nil {
		t.Errorf("failed to get DOT graph: %v", err)
	}

	if len(dot) == 0 {
		t.Errorf("expected non-empty DOT graph string")
	}
}
