package base

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// Matcher is matches predicate query
type Matcher struct {
	p     query.Predicate
	funcs []query.MatchFunc
}

func newMatcher(p query.Predicate, funcs ...query.MatchFunc) *Matcher {
	if len(funcs) == 0 {
		funcs = append(funcs, getDefaultMatchFunc(p))
	}

	return &Matcher{
		p:     p,
		funcs: funcs,
	}
}

// getDefaultMatchFunc returns default query.MatchFunc for the given predicate type
func getDefaultMatchFunc(p query.Predicate) query.MatchFunc {
	switch p.Type() {
	case query.UID:
		return UUIDEqFunc(p.Value().(uuid.UID))
	case query.Name, query.Group, query.Version, query.Kind, query.Namespace:
		return StringEqFunc(p.Value().(string))
	case query.Weight:
		return FloatEqFunc(p.Value().(float64))
	case query.Entity:
		return EntityEqFunc(p.Value().(query.EntityVal))
	case query.Attrs:
		return HasAttrsFunc(p.Value().(attrs.Attrs))
	}

	return IsAnyFunc
}

// Match returns true if the val matches the matcher's predicate value
// and all the query.MatchFuncs return true.
func (m Matcher) Match(val interface{}) bool {
	// NOTE: we first check if val is set to query.Any
	// which is a wildcard for matching any predicate
	var any bool
	if v, ok := val.(query.WildCard); ok {
		any = (v == query.Any)
	}

	match := true
	for _, fn := range m.funcs {
		match = (match && fn(val)) || any
	}

	return match
}

// Predicate returns query predicate
func (m Matcher) Predicate() query.Predicate { return m.p }
