package matcher

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/matcher/property"
)

// Matcher matches values for the given property.
type Matcher interface {
	Match(property.Property, interface{}) (bool, error)
}

type matcher struct {
	m map[property.Property]MatchFunc
}

// New creates new filter and returns it.
func New(opts ...Option) (*matcher, error) {
	fopts := Options{}
	for _, apply := range opts {
		apply(&fopts)
	}

	m := map[property.Property]MatchFunc{
		property.Type:      fopts.Types,
		property.Name:      fopts.Names,
		property.Group:     fopts.Groups,
		property.Version:   fopts.Versions,
		property.Kind:      fopts.Kinds,
		property.Namespace: fopts.Namespaces,
		property.Weight:    fopts.Weights,
		property.Attrs:     fopts.Attrs,
	}

	return &matcher{
		m: m,
	}, nil
}

// Matches matches val with the given property filter and returns the match result.
// It returns error if there is no MatchFunc defined for the given property p.
func (f matcher) Match(p property.Property, val interface{}) (bool, error) {
	skip, ok := f.m[p]
	if !ok {
		return false, ErrFilterNotFound
	}

	if skip == nil {
		return false, nil
	}

	return skip(val), nil
}

// Helper functions that return various match funcs for matching filter values.
func Types(v ...entity.Type) MatchFunc { return TypeEqFunc(v...) }
func Names(v ...string) MatchFunc      { return StringEqFunc(v...) }
func Groups(v ...string) MatchFunc     { return StringEqFunc(v...) }
func Versions(v ...string) MatchFunc   { return StringEqFunc(v...) }
func Kinds(v ...string) MatchFunc      { return StringEqFunc(v...) }
func Namespaces(v ...string) MatchFunc { return StringEqFunc(v...) }
func Weights(v ...float64) MatchFunc   { return FloatEqFunc(v...) }
func Attrs(v attrs.Attrs) MatchFunc    { return HasAttrsFunc(v) }
