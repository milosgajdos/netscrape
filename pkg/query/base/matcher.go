package base

import (
	"github.com/milosgajdos/netscrape/pkg/query"
)

// Matcher is matches predicate query
type Matcher struct {
	p     query.Predicate
	funcs []query.MatchFunc
}

func newMatcher(p query.Predicate, funcs ...query.MatchFunc) *Matcher {
	return &Matcher{
		p:     p,
		funcs: funcs,
	}
}

// Predicate returns query predicate
func (m Matcher) Predicate() query.Predicate { return m.p }

// Match returns true if the val matches the matcher's predicate value
// and all the query.MatchFuncs return true.
func (m Matcher) Match(val interface{}) bool {
	// NOTE: check if the value is set to query.Any
	// which is a placeholder to match any predicate value
	var any bool
	if v, ok := val.(query.Match); ok {
		any = (v == query.Any)
	}

	match := true
	for _, fn := range m.funcs {
		match = (match && fn(val)) || any
	}

	return match
}
