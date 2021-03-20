package memory

import (
	"math/big"
	"testing"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/internal"
)

const (
	node1DOTID = "node1ID"
	node2DOTID = "node2ID"
	edgeUID    = "testID"
	weight     = graph.DefaultWeight
)

func TestEdge(t *testing.T) {
	r, err := internal.NewTestResource(nodeResType, nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := internal.NewTestEntity(nodeID, nodeType, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n1, err := NewNode(1, o, graph.WithDOTID(node1DOTID))
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	o2, err := internal.NewTestEntity(nodeID, nodeType, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n2, err := NewNode(2, o2, graph.WithDOTID(node2DOTID))
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	// pass nil attributes
	e, err := NewEdge(n1, n2, graph.WithDOTID(edgeUID), graph.WithWeight(weight))
	if err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}

	if count := len(e.Attrs().Keys()); count != 0 {
		t.Errorf("expected 0 attributes, got: %d", count)
	}

	a, err := memattrs.New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}

	e, err = NewEdge(n1, n2, graph.WithDOTID(edgeUID), graph.WithWeight(weight), graph.WithAttrs(a))
	if err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}

	if uid := e.UID(); uid == nil {
		t.Fatalf("expected uid, got: %v", uid)
	}

	fromNode, err := e.FromNode()
	if err != nil {
		t.Fatalf("failed to get %v FromNode: %v", e.UID(), err)
	}

	if uid := fromNode.UID(); uid != n1.UID() {
		t.Errorf("expected ID: %s, got: %s", n1.UID(), uid)
	}

	toNode, err := e.ToNode()
	if err != nil {
		t.Fatalf("failed to get %v ToNode: %v", e.UID(), err)
	}

	if uid := toNode.UID(); uid != n2.UID() {
		t.Errorf("expected ID: %s, got: %s", n2.UID(), uid)
	}

	if uid := e.DOTID(); uid != edgeUID {
		t.Errorf("expected DOTID: %s, got: %s", edgeUID, uid)
	}

	if w := e.Weight(); big.NewFloat(w).Cmp(big.NewFloat(weight)) != 0 {
		t.Errorf("expected weight %f, got: %f", weight, w)
	}

	re := e.ReversedEdge()

	if re.From().ID() != e.To().ID() {
		t.Errorf("expected from ID: %d, got: %d", e.To().ID(), re.From().ID())
	}

	if re.To().ID() != e.From().ID() {
		t.Errorf("expected to UID: %d, got: %d", e.From().ID(), re.To().ID())
	}

	newDOTID := "DOTID"
	e.SetDOTID(newDOTID)

	if dotID := e.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}

	if dotAttrs := e.Attributes(); len(dotAttrs) != len(a.Attributes()) {
		t.Errorf("expected attributes: %d, got: %d", len(a.Attributes()), len(dotAttrs))
	}
}
