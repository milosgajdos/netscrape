package entity

import (
	"context"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
)

const (
	entUID     = "testEntUID"
	entType    = "testEntType"
	resUID     = "testResUID"
	resName    = "testResName"
	resType    = "testResType"
	resGroup   = "testResGroup"
	resVersion = "testResVersion"
	resKind    = "testResKind"
	resNsd     = true
	objUID     = "testObjID"
	objType    = "testObjType"
	objName    = "testObjName"
	objNs      = "testObjNs"
	testDOTID  = "dotID"
)

func MustSet(ctx context.Context, a attrs.Attrs, k, v string, t *testing.T) {
	if err := a.Set(ctx, k, v); err != nil {
		t.Fatalf("failed to set val %s for key %s: %v", k, v, err)
	}
}
