package memory

import (
	"context"
	"math/big"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/internal"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
)

const (
	node1DOTID = "node1ID"
	node2DOTID = "node2ID"
	edgeUID    = "testID"
	weight     = graph.DefaultWeight
)

func TestEdge(t *testing.T) {
	e, err := internal.NewTestObject()
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n1, err := NewNode(1, e, graph.WithDOTID(node1DOTID))
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	e2, err := internal.NewTestObject()
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	n2, err := NewNode(2, e2, graph.WithDOTID(node2DOTID))
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	// pass nil attributes
	if _, err := NewEdge(n1, n2, graph.WithDOTID(edgeUID), graph.WithWeight(weight)); err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}

	keys, err := e.Attrs().Keys(context.Background())
	if err != nil {
		t.Fatalf("failed to get attr keys: %v", err)
	}

	if count := len(keys); count != 0 {
		t.Errorf("expected 0 attributes, got: %d", count)
	}

	a := memattrs.New()

	edge, err := NewEdge(n1, n2, graph.WithDOTID(edgeUID), graph.WithWeight(weight), graph.WithAttrs(a))
	if err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}

	if uid := e.UID(); uid == nil {
		t.Fatalf("expected uid, got: %v", uid)
	}

	fromNode, err := edge.FromNode()
	if err != nil {
		t.Fatalf("failed to get %v FromNode: %v", e.UID(), err)
	}

	if uid := fromNode.UID(); uid != n1.UID() {
		t.Errorf("expected ID: %s, got: %s", n1.UID(), uid)
	}

	toNode, err := edge.ToNode()
	if err != nil {
		t.Fatalf("failed to get %v ToNode: %v", e.UID(), err)
	}

	if uid := toNode.UID(); uid != n2.UID() {
		t.Errorf("expected ID: %s, got: %s", n2.UID(), uid)
	}

	if uid := edge.DOTID(); uid != edgeUID {
		t.Errorf("expected DOTID: %s, got: %s", edgeUID, uid)
	}

	if w := edge.Weight(); big.NewFloat(w).Cmp(big.NewFloat(weight)) != 0 {
		t.Errorf("expected weight %f, got: %f", weight, w)
	}

	re := edge.ReversedEdge()

	if re.From().ID() != edge.To().ID() {
		t.Errorf("expected from ID: %d, got: %d", edge.To().ID(), re.From().ID())
	}

	if re.To().ID() != edge.From().ID() {
		t.Errorf("expected to UID: %d, got: %d", edge.From().ID(), re.To().ID())
	}

	newDOTID := "DOTID"
	edge.SetDOTID(newDOTID)

	if dotID := edge.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}

	if dotAttrs := edge.Attributes(); len(dotAttrs) != len(a.Attributes()) {
		t.Errorf("expected attributes: %d, got: %d", len(a.Attributes()), len(dotAttrs))
	}
}
