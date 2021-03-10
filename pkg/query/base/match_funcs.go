package base

import (
	"math/big"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/entity"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

// IsAnyFunc always returns true
func IsAnyFunc(v interface{}) bool {
	return true
}

// StringEqFunc returns MatchFunc option which checks
// the equality of an arbitrary string to s1
func StringEqFunc(s1 ...string) query.MatchFunc {
	return func(s2 interface{}) bool {
		for _, s := range s1 {
			if !(s == s2.(string)) {
				return false
			}
		}
		return true
	}
}

// EntityEqFunc returns MatchFunc option which checks
// the equality of an arbitrary entity.Type with e1.
func EntityEqFunc(e1 ...entity.Type) query.MatchFunc {
	return func(e2 interface{}) bool {
		for _, e := range e1 {
			if !(e == e2.(entity.Type)) {
				return false
			}
		}
		return true
	}
}

// FloatEqFunc returns MatchFunc which checks
// the equality of an arbitrary float to f1
func FloatEqFunc(f1 ...float64) query.MatchFunc {
	return func(f2 interface{}) bool {
		for _, f := range f1 {
			if big.NewFloat(f).Cmp(big.NewFloat(f2.(float64))) != 0 {
				return false
			}
		}
		return true
	}
}

// UIDEqFunc returns MatchFunc which checks
// the equality of an arbitrary uid to u1
func UUIDEqFunc(u1 ...uuid.UID) query.MatchFunc {
	return func(u2 interface{}) bool {
		for _, u := range u1 {
			if !(u.Value() == u2.(uuid.UID).Value()) {
				return false
			}
		}
		return true
	}
}

// HasAttrsFunc returns MatchFunc which checks
// if a contains k/v of an arbitrary attrs.Attrs
func HasAttrsFunc(a attrs.Attrs) query.MatchFunc {
	return func(a2 interface{}) bool {
		a2attrs := a2.(attrs.Attrs)
		for _, k := range a2attrs.Keys() {
			if v := a.Get(k); v != a2attrs.Get(k) {
				return false
			}
		}
		return true
	}
}
