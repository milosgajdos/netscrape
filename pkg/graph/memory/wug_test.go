package memory

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	wugObjPath = "testdata/wug/objects.yaml"
)

func TestWUGAddGetRemoveNode(t *testing.T) {
	g, err := NewWUG("test", graph.Options{})
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := newTestObject(nodeID, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object: %v", err)
	}

	n, err := g.NewNode(o, graph.NodeOptions{})
	if err != nil {
		t.Errorf("failed creating new graph node: %v", err)
	}

	// Add a new node
	if err := g.AddNode(n); err != nil {
		t.Errorf("failed adding node: %v", err)
	}

	nodes, err := g.Nodes()
	if err != nil {
		t.Fatalf("failed getting nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	// adding the same nodes twice should not change the node count
	if err := g.AddNode(n); err != nil {
		t.Errorf("failed adding node: %v", err)
	}

	expCount = 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	// Get the node with given uid
	node, err := g.Node(n.UID())
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

	if _, err := g.Node(guid); err != graph.ErrNodeNotFound {
		t.Errorf("expected error %v, got: %#v", graph.ErrNodeNotFound, err)
	}

	// Remove the node with given uid
	if err := g.RemoveNode(n.UID()); err != nil {
		t.Errorf("failed to remove node: %v", err)
	}

	nodes, err = g.Nodes()
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

	if err := g.RemoveNode(guid); err != nil {
		t.Errorf("failed to remove node: %v", err)
	}
}

func TestWUGLinkGetRemoveEdge(t *testing.T) {
	g, err := NewWUG("test", graph.Options{})
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	node1UID := "foo1UID"
	node1Name := "foo1Name"

	o1, err := newTestObject(node1UID, node1Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object %q: %v", node1UID, err)
	}

	n1, err := g.NewNode(o1, graph.NodeOptions{})
	if err != nil {
		t.Errorf("failed creating new node: %v", err)
	}

	if err := g.AddNode(n1); err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	node2UID := "foo2UID"
	node2Name := "foo2Name"

	o2, err := newTestObject(node2UID, node2Name, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object %q: %v", node2UID, err)
	}

	n2, err := g.NewNode(o2, graph.NodeOptions{})
	if err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	if err := g.AddNode(n2); err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	ox, err := newTestObject("nonExUID", "nonExName", nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object %q: %v", node2UID, err)
	}

	nodeX := &Node{
		Object: ox,
		id:     123334444,
	}

	// Link nodes with a node which does not exist in the graph
	if _, err := g.Link(n1.UID(), nodeX.UID(), graph.LinkOptions{}); !errors.Is(err, graph.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", graph.ErrNodeNotFound, err)
	}

	if _, err := g.Link(nodeX.UID(), n2.UID(), graph.LinkOptions{}); !errors.Is(err, graph.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", graph.ErrNodeNotFound, err)
	}

	edges, err := g.Edges()
	if err != nil {
		t.Errorf("failed getting graph edges: %v", err)
	}

	expCount := 0
	if len(edges) != expCount {
		t.Errorf("expected: %d edges, got: %d", expCount, len(edges))
	}

	edge, err := g.Link(n1.UID(), n2.UID(), graph.LinkOptions{Weight: graph.DefaultWeight})
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", n1.UID(), n2.UID(), err)
	}

	if w := edge.Weight(); big.NewFloat(w).Cmp(big.NewFloat(graph.DefaultWeight)) != 0 {
		t.Errorf("expected non-negative weight")
	}

	edges, err = g.Edges()
	if err != nil {
		t.Errorf("failed getting graph edges: %v", err)
	}

	expCount = 1
	if len(edges) != expCount {
		t.Errorf("no edges found in graph")
	}

	// linking already linked nodes must return the same edge/line as returned previously
	exEdge, err := g.Link(n1.UID(), n2.UID(), graph.LinkOptions{})
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", n1.UID(), n2.UID(), err)
	}

	if !reflect.DeepEqual(exEdge, edge) {
		t.Errorf("expected edge %#v, got: %#v", exEdge, edge)
	}

	e, err := g.Edge(n1.UID(), n2.UID())
	if err != nil {
		t.Errorf("failed getting edge between %s and %s: %v", n1.UID(), n2.UID(), err)
	}

	if !reflect.DeepEqual(e, edge) {
		t.Errorf("expected edge %#v, got: %#v", exEdge, edge)
	}

	// remove edge between previously linked nodes which are still present in the graph
	if err := g.RemoveLink(n1.UID(), n2.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", n1.UID(), n2.UID(), err)
	}

	if _, err := g.Edge(n1.UID(), n2.UID()); err != nil && !errors.Is(err, graph.ErrEdgeNotExist) {
		t.Errorf("expected error: %v, got: %v", graph.ErrEdgeNotExist, err)
	}

	// remoe edge between non-existen nodes should return nil
	if err := g.RemoveLink(nodeX.UID(), n1.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", nodeX.UID(), n1.UID(), err)
	}

	if err := g.RemoveLink(n1.UID(), nodeX.UID()); err != nil {
		t.Errorf("failed removing edge between %s and %s: %v", nodeX.UID(), n1.UID(), err)
	}
}

func TestWUGSubGraph(t *testing.T) {
	g, err := makeTestGraph(wugObjPath)
	if err != nil {
		t.Fatalf("failed to create new memory graph: %v", err)
	}

	guid, err := uuid.NewFromString("garbage")
	if err != nil {
		t.Fatalf("error creating new uid: %v", err)
	}

	// subgraph of non-existent node should return error
	if _, err := g.SubGraph(guid, 10); err != graph.ErrNodeNotFound {
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
		sg, err := g.SubGraph(uid, tc.depth)
		if err != nil {
			t.Errorf("failed to get subgraph of node %s: %v", uid, err)
			continue
		}

		storeNodes, err := sg.Nodes()
		if err != nil {
			t.Errorf("failed to fetch subgraph nodes: %v", err)
			continue
		}

		if len(storeNodes) != tc.exp {
			t.Errorf("expected subgraph nodes: %d, got: %d", tc.exp, len(storeNodes))
		}
	}
}

func TestWUGQueryEdge(t *testing.T) {
	g, err := makeTestGraph(wugObjPath)
	if err != nil {
		t.Fatalf("failed to create test graph: %v", err)
	}

	q := base.Build().Add(query.Entity(query.Edge))

	qedges, err := g.Query(q)
	if err != nil {
		t.Errorf("failed to query edges: %v", err)
	}

	edges, err := g.Edges()
	if err != nil {
		t.Fatalf("failed to fetch edges: %v", err)
	}

	if len(qedges) != len(edges) {
		t.Errorf("expected edges: %d, got: %d", len(edges), len(qedges))
	}

	q = base.Build().Add(query.Entity(query.Node))

	nodes, err := g.Query(q)
	if err != nil {
		t.Errorf("failed to query nodes: %v", err)
	}

	relations := make(map[string]bool)

	for _, n := range nodes {
		for _, l := range n.(graph.Node).Links() {
			if r, ok := l.Metadata().Get("relation").(string); ok {
				relations[r] = true
			}
		}
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}

	for r, ok := range relations {
		if ok {
			a.Set("relation", r)

			q = base.Build().
				Add(query.Entity(query.Edge)).
				Add(query.Attrs(a), query.HasAttrsFunc(a))

			edges, err := g.Query(q)
			if err != nil {
				t.Errorf("failed querying edges with attributes %v: %v", a, err)
			}

			for _, edge := range edges {
				for _, k := range a.Keys() {
					v := a.Get(k)
					if val := edge.Attrs().Get(k); val != v {
						t.Errorf("expected attributes: %v:%v, got: %v:%v", k, v, k, val)
					}
				}
			}
		}
	}
}

func TestWUGQueryNode(t *testing.T) {
	g, err := makeTestGraph(wugObjPath)
	if err != nil {
		t.Fatalf("failed to create test graph: %v", err)
	}

	q := base.Build().Add(query.Entity(query.Node))

	qnodes, err := g.Query(q)
	if err != nil {
		t.Errorf("failed to query all nodes: %v", err)
	}

	nodes, err := g.Nodes()
	if err != nil {
		t.Fatalf("failed to fetch nodes: %v", err)
	}

	if len(qnodes) != len(nodes) {
		t.Errorf("expected nodes: %d, got: %d", len(nodes), len(qnodes))
	}

	namespaces := make([]string, len(nodes))
	kinds := make([]string, len(nodes))
	names := make([]string, len(nodes))

	for i, n := range nodes {
		namespaces[i] = n.Namespace()
		kinds[i] = n.Resource().Kind()
		names[i] = n.Name()
	}

	q = base.Build().Add(query.Entity(query.Node))

	for _, ns := range namespaces {
		q = q.Add(query.Namespace(ns), query.StringEqFunc(ns))

		nodes, err := g.Query(q)
		if err != nil {
			t.Errorf("error getting namespace %s nodes: %v", ns, err)
			continue
		}

		for _, n := range nodes {
			if nodeNS := n.(graph.Node).Namespace(); nodeNS != ns {
				t.Errorf("expected: namespace %s, got: %s", ns, nodeNS)
			}
		}

		for _, kind := range kinds {
			q = q.Add(query.Kind(kind), query.StringEqFunc(kind))

			nodes, err := g.Query(q)
			if err != nil {
				t.Errorf("error getting nodes: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, n := range nodes {
				o := n.(graph.Node)
				if o.Namespace() != ns || o.Resource().Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, o.Namespace(), o.Resource().Kind())
				}
			}
		}
	}
}

func TestWUGQuery(t *testing.T) {
	g, err := makeTestGraph(wugObjPath)
	if err != nil {
		t.Fatalf("failed to create new memory graph: %v", err)
	}

	// NOTE: EntityVal is an enum/iota starting with 0 with only two values: Node and Edge
	// Any other number higher than 1 is considered a non-existent Entity
	q := base.Build().Add(query.Entity(query.EntityVal(10000)), query.IsAnyFunc)

	if _, err := g.Query(q); !errors.Is(err, graph.ErrUnknownEntity) {
		t.Errorf("expected: %v, got: %v", graph.ErrUnknownEntity, err)
	}
}

func TestWUGDOT(t *testing.T) {
	id := "testID"

	g, err := NewWUG(id, graph.Options{})
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
