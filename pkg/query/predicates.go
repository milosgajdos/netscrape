package query

import (
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

type predicate struct {
	kind  PredKind
	value interface{}
}

// NewPredicate creates a new predicate and returns it.
func NewPredicate(k PredKind, v interface{}) Predicate {
	return predicate{
		kind:  k,
		value: v,
	}
}

func (p predicate) Kind() PredKind     { return p.kind }
func (p predicate) Value() interface{} { return p.value }
func (p predicate) String() string     { return fmt.Sprint(p.value) }

// Helper functions that return query predicates
func UID(v uuid.UID) Predicate               { return predicate{kind: PUID, value: v} }
func Name(v string) Predicate                { return predicate{kind: PName, value: v} }
func Group(v string) Predicate               { return predicate{kind: PGroup, value: v} }
func Version(v string) Predicate             { return predicate{kind: PVersion, value: v} }
func Kind(v string) Predicate                { return predicate{kind: PKind, value: v} }
func Namespace(v string) Predicate           { return predicate{kind: PNamespace, value: v} }
func Weight(v float64) Predicate             { return predicate{kind: PWeight, value: v} }
func Entity(v EntityVal) Predicate           { return predicate{kind: PEntity, value: v} }
func Attrs(v attrs.Attrs) Predicate          { return predicate{kind: PAttrs, value: v} }
func Metadata(v metadata.Metadata) Predicate { return predicate{kind: PMetadata, value: v} }
