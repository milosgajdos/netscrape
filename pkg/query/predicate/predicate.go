package predicate

import (
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

type predicate struct {
	typ query.Type
	val interface{}
}

// New creates new predicate and returns it.
func New(t query.Type, v interface{}) predicate {
	return predicate{
		typ: t,
		val: v,
	}
}

// Type returns predicate type
func (p predicate) Type() query.Type { return p.typ }

// Value returns predicate Value
func (p predicate) Value() interface{} { return p.val }

// String implements fmt.Stringer
func (p predicate) String() string { return fmt.Sprint(p.val) }

// Helper functions that return query predicates.
func UID(v uuid.UID) query.Predicate           { return New(query.UID, v) }
func Name(v string) query.Predicate            { return New(query.Name, v) }
func Group(v string) query.Predicate           { return New(query.Group, v) }
func Version(v string) query.Predicate         { return New(query.Version, v) }
func Kind(v string) query.Predicate            { return New(query.Kind, v) }
func Namespace(v string) query.Predicate       { return New(query.Namespace, v) }
func Weight(v float64) query.Predicate         { return New(query.Weight, v) }
func Entity(v query.EntityVal) query.Predicate { return New(query.Entity, v) }
func Attrs(v attrs.Attrs) query.Predicate      { return New(query.Attrs, v) }
