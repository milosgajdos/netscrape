package base

import (
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/query"
)

// Query is a base query
type Query struct {
	matchers map[query.PredKind]*Matcher
}

// Build creates a new Query and returns it
func Build() query.Query {
	q := &Query{
		matchers: make(map[query.PredKind]*Matcher),
	}

	return q.MatchAny()
}

// Add adds new predicates to query.
func (q *Query) Add(p query.Predicate, funcs ...query.MatchFunc) query.Query {
	q.matchers[p.Kind()] = newMatcher(p, funcs...)
	return q
}

// Matcher returns the Matcher for the given predicate kind
func (q *Query) Matcher(k query.PredKind) query.Matcher {
	return q.matchers[k]
}

// Reset resets the query.
func (q *Query) Reset() query.Query {
	return Build()
}

// String implements fmt.Stringer
func (q *Query) String() string {
	var result string
	for k, m := range q.matchers {
		result += fmt.Sprintf("%s: %s\n", k, m.Predicate().Value())
	}

	return result
}

// MatchAny returns query which matches any predicate
func (q *Query) MatchAny() query.Query {
	for _, k := range []query.PredKind{
		query.PUID,
		query.PName,
		query.PGroup,
		query.PVersion,
		query.PKind,
		query.PNamespace,
		query.PWeight,
		query.PEntity,
		query.PAttrs,
		query.PMetadata,
	} {
		q.matchers[k] = newMatcher(query.NewPredicate(k, query.Any), query.IsAnyFunc)
	}

	return q
}
