package query

// MatchFunc runs when query matcher matches a value.
type MatchFunc func(interface{}) bool

// Query is an experimental query interface for querying
// Space and Store. This is merely an experiment alas it is
// used in the codebase for querying Space and Topology.
type Query interface {
	// Add adds a new predicate with MatchFunc to Query
	Add(Predicate, ...MatchFunc) Query
	// TODO: consider returning one Matcher for all
	// Predicates and allow a specific one be picked from it.
	// Matcher returns Matcher for given Predicate
	Matcher(PredKind) Matcher
	// Reset resets Query
	Reset() Query
	// String implements fmt.Stringer
	String() string
}

// Predicate is query predicate
type Predicate interface {
	// Kind returns predicate kind
	Kind() PredKind
	// Value returns predicate value
	Value() interface{}
	// String implements fmt.Stringer
	String() string
}

// Matcher returns query matcher
type Matcher interface {
	// Match the given value
	Match(interface{}) bool
	// Predicate returns predicate
	Predicate() Predicate
}
