package memory

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
	"github.com/milosgajdos/netscrape/pkg/internal"
	"github.com/milosgajdos/netscrape/pkg/space"
)

const (
	nodeGID = 123
)

func TestNode(t *testing.T) {
	r, err := internal.NewTestResource(nodeResType, nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	e, err := internal.NewTestEntity(nodeID, nodeType, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}
	a.Set("nodename", nodeName)

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

	// NOTE: by default we will get the following attributes:
	// * graph.DOTIDAttr
	// * graph.NameAttr
	// We added "nodename" attribute above which leaves us with 3 attributes altogether.
	if dotAttrs := n.Attributes(); len(dotAttrs) != 3 {
		t.Errorf("expected %d attributes, got: %d", 3, len(dotAttrs))
	}

	newDOTID := "DOTID"
	n.SetDOTID(newDOTID)

	if dotID := n.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}

	if count := len(n.Attrs().Keys()); count == 0 {
		t.Fatalf("expected node attributes got: %d", count)
	}
}

func TestNodeWithDOTID(t *testing.T) {
	r, err := internal.NewTestResource(nodeResType, nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := internal.NewTestEntity(nodeID, nodeType, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}
	a.Set("name", nodeName)

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

	// NOTE: we expect the node to have 2 attributes:
	// * name: set when the node was created with given attrs options
	// * dotid: set via node.SetDOTID
	exp := 2
	if dotAttrs := node.Attributes(); len(dotAttrs) != exp {
		t.Errorf("expected attributes: %d, got: %d", exp, len(dotAttrs))
	}
}
