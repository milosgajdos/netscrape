package base

import (
	"fmt"

	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/predicate"
)

// Query is a base query
type Query struct {
	matchers map[query.Type]*Matcher
}

// Build creates a new Query and returns it
func Build() query.Query {
	q := &Query{
		matchers: make(map[query.Type]*Matcher),
	}

	return q.MatchAny()
}

// Add adds new predicates to query.
func (q *Query) Add(p query.Predicate, funcs ...query.MatchFunc) query.Query {
	q.matchers[p.Type()] = newMatcher(p, funcs...)
	return q
}

// Matcher returns the Matcher for the given predicate kind
func (q *Query) Matcher(t query.Type) query.Matcher {
	return q.matchers[t]
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
	for _, k := range []query.Type{
		query.UID,
		query.Name,
		query.Group,
		query.Version,
		query.Kind,
		query.Namespace,
		query.Weight,
		query.Entity,
		query.Attrs,
	} {
		q.matchers[k] = newMatcher(predicate.New(k, query.Any), IsAnyFunc)
	}

	return q
}
