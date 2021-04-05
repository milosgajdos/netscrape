package memory

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/internal"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

const (
	wdgEntPath = "testdata/wdg/entities.yaml"
)

func TestWDGAddGetRemoveNode(t *testing.T) {
	g, err := NewWDG()
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	if uid := g.UID(); uid == nil {
		t.Errorf("expected uid, got: %v", uid)
	}

	r, err := internal.NewTestResource(nodeResType, nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := internal.NewTestEntity(nodeType, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n, err := g.NewNode(context.Background(), o)
	if err != nil {
		t.Errorf("failed creating new graph node: %v", err)
	}

	if err := g.AddNode(context.Background(), n); err != nil {
		t.Errorf("failed adding node: %v", err)
	}

	nodes, err := g.Nodes(context.Background())
	if err != nil {
		t.Fatalf("failed getting nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	// adding the same node without upsert should return error.
	if err := g.AddNode(context.Background(), n); !errors.Is(err, graph.ErrDuplicateNode) {
		t.Errorf("expected error: %v, got: %v", graph.ErrDuplicateNode, err)
	}

	expCount = 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	node, err := g.Node(context.Background(), n.UID())
	if err != nil {
		t.Errorf("failed to get %s node: %v", n.UID(), err)
	}

	if !reflect.DeepEqual(n, node) {
		t.Errorf("expected node %#v, got: %#v", node, n)
	}

	guid := memuid.NewFromString("garbage")

	if _, err := g.Node(context.Background(), guid); err != graph.ErrNodeNotFound {
		t.Errorf("expected error %v, got: %#v", graph.ErrNodeNotFound, err)
	}

	if err := g.RemoveNode(context.Background(), n.UID()); err != nil {
		t.Errorf("failed to remove node: %v", err)
	}

	nodes, err = g.Nodes(context.Background())
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount = 0
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	guid = memuid.NewFromString("garbage")

	if err := g.RemoveNode(context.Background(), guid); err != nil {
		t.Errorf("failed to remove node: %v", err)
	}
}

func TestWDGLinkGetRemoveEdge(t *testing.T) {
	g, err := NewWDG()
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	r, err := internal.NewTestResource(nodeResType, nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o1, err := internal.NewTestEntity(nodeType, "foo1Name", nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n1, err := g.NewNode(context.Background(), o1)
	if err != nil {
		t.Errorf("failed creating new node: %v", err)
	}

	if err := g.AddNode(context.Background(), n1); err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	o2, err := internal.NewTestEntity(nodeType, "foo2Name", nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n2, err := g.NewNode(context.Background(), o2)
	if err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	if err := g.AddNode(context.Background(), n2); err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	ox, err := internal.NewTestEntity(nodeType, "nonExName", nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	nodeX, err := NewNode(123334444, ox)
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	// Link nodes with a node which does not exist in the graph
	if _, err := g.Link(context.Background(), n1.UID(), nodeX.UID()); !errors.Is(err, graph.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", graph.ErrNodeNotFound, err)
	}

	if _, err := g.Link(context.Background(), nodeX.UID(), n2.UID()); !errors.Is(err, graph.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", graph.ErrNodeNotFound, err)
	}

	edges, err := g.Edges(context.Background())
	if err != nil {
		t.Errorf("failed getting graph edges: %v", err)
	}

	expCount := 0
	if len(edges) != expCount {
		t.Errorf("expected: %d edges, got: %d", expCount, len(edges))
	}

	edge, err := g.Link(context.Background(), n1.UID(), n2.UID(), graph.WithWeight(graph.DefaultWeight))
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", n1.UID(), n2.UID(), err)
	}

	nodesFrom, err := g.From(context.Background(), n1.UID())
	if err != nil {
		t.Errorf("failed to get nodes from %s: %v", n1.UID(), err)
	}

	expCount = 1
	if count := len(nodesFrom); count != expCount {
		t.Errorf("expected: %d nodes, got: %d", expCount, count)
	}

	if len(nodesFrom) == 1 {
		if nodesFrom[0].UID().String() != n2.UID().String() {
			t.Errorf("expected node link to %s from %s", n2.UID().String(), n1.UID().String())
		}
	}

	if w := edge.Weight(); big.NewFloat(w).Cmp(big.NewFloat(graph.DefaultWeight)) != 0 {
		t.Errorf("expected non-negative weight")
	}

	edges, err = g.Edges(context.Background())
	if err != nil {
		t.Errorf("failed getting graph edges: %v", err)
	}

	expCount = 1
	if len(edges) != expCount {
		t.Errorf("no edges found in graph")
	}

	// linking already linked nodes must return the same edge/line as returned previously
	exEdge, err := g.Link(context.Background(), n1.UID(), n2.UID())
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", n1.UID(), n2.UID(), err)
	}

	if !reflect.DeepEqual(exEdge, edge) {
		t.Errorf("expected edge %#v, got: %#v", exEdge, edge)
	}

	e, err := g.Edge(context.Background(), n1.UID(), n2.UID())
	if err != nil {
		t.Errorf("failed getting edge between %s and %s: %v", n1.UID(), n2.UID(), err)
	}

	if !reflect.DeepEqual(e, edge) {
		t.Errorf("expected edge %#v, got: %#v", exEdge, edge)
	}

	// remove edge between previously linked nodes which are still present in the graph
	if err := g.Unlink(context.Background(), n1.UID(), n2.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", n1.UID(), n2.UID(), err)
	}

	if _, err := g.Edge(context.Background(), n1.UID(), n2.UID()); err != nil && !errors.Is(err, graph.ErrEdgeNotExist) {
		t.Errorf("expected error: %v, got: %v", graph.ErrEdgeNotExist, err)
	}

	// remoe edge between non-existent nodes should return nil
	if err := g.Unlink(context.Background(), nodeX.UID(), n1.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", nodeX.UID(), n1.UID(), err)
	}

	if err := g.Unlink(context.Background(), n1.UID(), nodeX.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", nodeX.UID(), n1.UID(), err)
	}
}

func TestWDGSubGraph(t *testing.T) {
	g, err := makeTestGraph(wdgEntPath)
	if err != nil {
		t.Fatalf("failed to create new memory graph: %v", err)
	}

	guid := memuid.NewFromString("garbage")

	// subgraph of non-existent node should return error
	if _, err := g.SubGraph(context.Background(), guid, 10); err != graph.ErrNodeNotFound {
		t.Errorf("expected: %v, got: %v", graph.ErrNodeNotFound, err)
	}

	// NOTE: we are hardcoding the test value here
	// since we know that this node's neighbourhood
	suid := "fooGroup/v1/fooKind/fooNs/foo1"
	uid := memuid.NewFromString(suid)

	testCases := []struct {
		depth int
		exp   int
	}{
		{0, 1},   // return node
		{1, 5},   // return node + adjacent nodes
		{100, 8}, // return all nodes reachable from node
	}

	for _, tc := range testCases {
		sg, err := g.SubGraph(context.Background(), uid, tc.depth)
		if err != nil {
			t.Errorf("failed to get subgraph of node %s: %v", uid, err)
			continue
		}

		storeNodes, err := sg.Nodes(context.Background())
		if err != nil {
			t.Errorf("failed to fetch subgraph nodes: %v", err)
			continue
		}

		if len(storeNodes) != tc.exp {
			t.Errorf("expected subgraph nodes: %d, got: %d", tc.exp, len(storeNodes))
		}
	}
}

func TestWDGDOT(t *testing.T) {
	id := "testID"

	g, err := NewWDG(graph.WithDOTID(id))
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
