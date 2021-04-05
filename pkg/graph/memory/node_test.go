package memory

import (
	"context"
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/internal"
	"github.com/milosgajdos/netscrape/pkg/space"

	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
)

const (
	nodeGID = 123
)

func TestNode(t *testing.T) {
	r, err := internal.NewTestResource(nodeResType, nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	e, err := internal.NewTestEntity(nodeType, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	a := memattrs.New()
	MustSet(context.Background(), a, "nodename", nodeName, t)

	n, err := NewNode(nodeGID, e, graph.WithAttrs(a))
	if err != nil {
		t.Fatalf("failed to create new node from Space entity: %v", err)
	}

	if id := n.ID(); id != nodeGID {
		t.Errorf("expected ID: %d, got: %d", nodeGID, id)
	}

	if nodeEnt := n.Entity.(space.Entity); !reflect.DeepEqual(nodeEnt, e) {
		t.Errorf("invalid graph.Entity for node: %s", n.UID())
	}

	if dotEnt, ok := e.(graph.DOTEntity); ok {
		dotid := dotEnt.DOTID()
		if dotID := n.DOTID(); dotID != dotid {
			t.Errorf("expected DOTID: %s, got: %s", dotid, dotID)
		}

	}

	// NOTE: we set the "nodename" attribute above
	exp := 1
	if dotAttrs := n.Attributes(); len(dotAttrs) != exp {
		t.Errorf("expected %d attributes, got: %d", exp, len(dotAttrs))
	}

	newDOTID := "DOTID"
	n.SetDOTID(newDOTID)

	if dotID := n.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}

	keys, err := n.Attrs().Keys(context.Background())
	if err != nil {
		t.Fatalf("failed to get attr keys: %v", err)
	}

	if count := len(keys); count == 0 {
		t.Fatalf("expected node attributes got: %d", count)
	}
}

func TestNodeWithDOTID(t *testing.T) {
	r, err := internal.NewTestResource(nodeResType, nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := internal.NewTestEntity(nodeType, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	a := memattrs.New()
	MustSet(context.Background(), a, "name", nodeName, t)

	node, err := NewNode(nodeGID, o, graph.WithDOTID(nodeName), graph.WithAttrs(a))
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	if id := node.ID(); id != nodeGID {
		t.Errorf("expected ID: %d, got: %d", nodeGID, id)
	}

	if dotID := node.DOTID(); dotID != nodeName {
		t.Errorf("expected DOTID: %s, got: %s", nodeName, dotID)
	}

	newDOTID := "DOTID"
	node.SetDOTID(newDOTID)

	if dotID := node.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}

	// NOTE: we set the "name" attribute
	exp := 1
	if dotAttrs := node.Attributes(); len(dotAttrs) != exp {
		t.Errorf("expected attributes: %d, got: %d", exp, len(dotAttrs))
	}
}
