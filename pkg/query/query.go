package query

// Query is an experimental query interface.
// This is *truly a wild experiment*, alas it is
// used for querying space.Plan and space.Top.
type Query interface {
	// Add a new predicate with MatchFunc to Query
	Add(Predicate, ...MatchFunc) Query
	// Matcher returns Matcher for given Type
	Matcher(Type) Matcher
	// Reset resets Query
	Reset() Query
	// String implements fmt.Stringer
	String() string
}

// MatchFunc to execute when matching a value.
type MatchFunc func(interface{}) bool

// Predicate is query predicate
type Predicate interface {
	// Type returns predicate type
	Type() Type
	// Value returns predicate value
	Value() interface{}
	// String implements fmt.Stringer
	String() string
}

// Matcher returns query matcher
type Matcher interface {
	// Match given value
	Match(interface{}) bool
	// Predicate returns matcher Predicate
	Predicate() Predicate
}
