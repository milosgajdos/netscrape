package matcher

import (
	"math/big"

	"github.com/milosgajdos/netscrape/pkg/attrs"
)

type MatchFunc func(interface{}) bool

// StringEqFunc returns MatchFunc that checks
// the equality of a string to all strings in sx.
func StringEqFunc(sx ...string) MatchFunc {
	return func(v interface{}) bool {
		vs, ok := v.(string)
		if !ok {
			return false
		}
		for _, s := range sx {
			if s == vs {
				return true
			}
		}
		return false
	}
}

// FloatEqFunc returns MatchFunc which checks
// the equality of an arbitrary float to f1
func FloatEqFunc(fx ...float64) MatchFunc {
	return func(v interface{}) bool {
		vf, ok := v.(float64)
		if !ok {
			return false
		}
		for _, f := range fx {
			if big.NewFloat(f).Cmp(big.NewFloat(vf)) == 0 {
				return true
			}
		}
		return false
	}
}

// HasAttrsFunc returns MatchFunc which checks
// if a contains k/v of an arbitrary attrs.Attrs
func HasAttrsFunc(a attrs.Attrs) MatchFunc {
	return func(v interface{}) bool {
		va, ok := v.(attrs.Attrs)
		if !ok {
			return false
		}
		for _, k := range va.Keys() {
			if val := a.Get(k); val != va.Get(k) {
				return false
			}
		}
		return true
	}
}
