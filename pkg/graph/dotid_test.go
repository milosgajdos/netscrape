package graph

import (
	"strings"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/space/object"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	resName    = "resName"
	resGroup   = "resGroup"
	resVersion = "resVersion"
	resKind    = "resKind"
	objUID     = "testID"
	objName    = "testName"
	objNs      = "testNs"
)

func TestDOTID(t *testing.T) {
	uid, err := uuid.NewFromString(objUID)
	if err != nil {
		t.Fatalf("failed to create new uid from string %q: %v", objUID, err)
	}

	o, err := object.New(uid, objName, objNs, nil)
	if err != nil {
		t.Fatalf("failed to create new object: %v", err)
	}

	if _, err := DOTID(o); err == nil {
		t.Errorf("expected error, got: %v", err)
	}

	r, err := resource.New(resName, resGroup, resVersion, resKind, true)
	if err != nil {
		t.Fatalf("failed to create new resource: %v", err)
	}

	o, err = object.New(uid, objName, objNs, r)
	if err != nil {
		t.Fatalf("failed to create new object: %v", err)
	}

	exp := strings.Join([]string{
		resGroup,
		resVersion,
		resKind,
		objNs,
		objName}, "/")

	dotid, err := DOTID(o)
	if err != nil {
		t.Fatalf("failed to generate DOTID: %v", err)
	}

	if dotid != exp {
		t.Errorf("expected DOTID: %s, got: %s", exp, dotid)
	}
}
