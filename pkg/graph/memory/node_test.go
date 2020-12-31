package memory

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/graph"
)

func TestNode(t *testing.T) {
	if _, err := newTestObject(nodeID, nodeName, nodeNs, nil); err != nil {
		t.Fatalf("failed to create object: %v", err)
	}

	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := newTestObject(nodeID, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object: %v", err)
	}

	dotid, err := graph.DOTID(o)
	if err != nil {
		t.Fatalf("failed to build DOTID: %v", err)
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}
	a.Set("nodename", nodeName)

	node, err := NewNode(nodeGID, o, graph.NodeOptions{Attrs: a})
	if err != nil {
		t.Fatalf("failed to create new node from Space object: %v", err)
	}

	if id := node.ID(); id != nodeGID {
		t.Errorf("expected ID: %d, got: %d", nodeGID, id)
	}

	if nodeObj := node.Object; !reflect.DeepEqual(nodeObj, o) {
		t.Errorf("invalid space.Object for node: %s", node.UID())
	}

	if dotID := node.DOTID(); dotID != dotid {
		t.Errorf("expected DOTID: %s, got: %s", dotid, dotID)
	}

	// NOTE: by default we will get the following attributes:
	// * dotid
	// * name
	// We add "nodename" attribute which leaves us with 3 attributes altogether.
	if dotAttrs := node.Attributes(); len(dotAttrs) != 3 {
		t.Errorf("expected %d attributes, got: %d", 3, len(dotAttrs))
	}

	newDOTID := "DOTID"
	node.SetDOTID(newDOTID)

	if dotID := node.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}
}

func TestNodeWithDOTID(t *testing.T) {
	r, err := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	o, err := newTestObject(nodeID, nodeName, nodeNs, r)
	if err != nil {
		t.Fatalf("failed to create object: %v", err)
	}

	a, err := attrs.New()
	if err != nil {
		t.Fatalf("failed to create attrs: %v", err)
	}
	a.Set("name", nodeName)

	node, err := NewNodeWithDOTID(nodeGID, o, nodeName, graph.NodeOptions{Attrs: a})
	if err != nil {
		t.Errorf("failed to create new node: %v", err)
		return
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
