package memory

import (
	"math/big"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/graph"
)

const (
	node1DOTID = "node1ID"
	node2DOTID = "node2ID"
	edgeUID    = "testID"
	weight     = graph.DefaultWeight
)

func TestEdge(t *testing.T) {
	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := newTestObject(nodeID, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object: %v", err)
	}

	n1, err := NewNodeWithDOTID(1, o, node1DOTID, graph.NodeOptions{})
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	o2, err := newTestObject(nodeID, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object: %v", err)
	}

	n2, err := NewNodeWithDOTID(2, o2, node2DOTID, graph.NodeOptions{})
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}

	e, err := NewEdgeWithDOTID(edgeUID, n1, n2, weight, entity.Attrs(a))
	if err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}

	if uid := e.FromNode().UID(); uid != n1.UID() {
		t.Errorf("expected ID: %s, got: %s", n1.UID(), uid)
	}

	if uid := e.ToNode().UID(); uid != n2.UID() {
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
