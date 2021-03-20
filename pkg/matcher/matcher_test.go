package matcher

import (
	"testing"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	memattrs "github.com/milosgajdos/netscrape/pkg/attrs/memory"
	"github.com/milosgajdos/netscrape/pkg/matcher/property"
)

func MustMatcher(t *testing.T, opts ...Option) *matcher {
	f, err := New(opts...)
	if err != nil {
		t.Fatalf("failed creating filter: %v", err)
	}

	return f
}

func TestFilterNone(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("NoOpts", func(t *testing.T) {
		f := MustMatcher(t)
		got, err := f.Match(property.Name, "foo")
		if err != nil {
			t.Fatalf("failed to match: %v", err)
		}

		exp := false
		if got != exp {
			t.Errorf("expected: %v, got: %v", exp, got)
		}
	})

	t.Run("InvalidProperty", func(t *testing.T) {
		f := MustMatcher(t)
		if _, err := f.Match(property.Property(-1000), "foo"); err != ErrFilterNotFound {
			t.Fatalf("expected error: %v, got: %v", ErrFilterNotFound, err)
		}
	})
}

func TestFilterTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			types []string
			val   string
			exp   bool
		}{
			{[]string{}, "foo", false},
			{[]string{"foo"}, "foo", true},
			{[]string{"bar"}, "foo", false},
			{[]string{"foo", "bar"}, "foo", true},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithTypes(Types(tc.types...)))
			got, err := f.Match(property.Type, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}

func TestFilterNames(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			names []string
			val   string
			exp   bool
		}{
			{[]string{}, "foo", false},
			{[]string{"foo", "bar"}, "foo", true},
			{[]string{"foo"}, "bar", false},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithNames(Names(tc.names...)))
			got, err := f.Match(property.Name, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}

func TestFilterGroups(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			groups []string
			val    string
			exp    bool
		}{
			{[]string{}, "foo", false},
			{[]string{"foo", "bar"}, "foo", true},
			{[]string{"foo"}, "bar", false},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithGroups(Groups(tc.groups...)))
			got, err := f.Match(property.Group, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}

func TestFilterVersions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			versions []string
			val      string
			exp      bool
		}{
			{[]string{}, "foo", false},
			{[]string{"foo", "bar"}, "foo", true},
			{[]string{"foo"}, "bar", false},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithVersions(Versions(tc.versions...)))
			got, err := f.Match(property.Version, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}

func TestFilterKinds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			kinds []string
			val   string
			exp   bool
		}{
			{[]string{}, "foo", false},
			{[]string{"foo", "bar"}, "foo", true},
			{[]string{"foo"}, "bar", false},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithKinds(Kinds(tc.kinds...)))
			got, err := f.Match(property.Kind, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}

func TestFilterNamespaces(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			ns  []string
			val string
			exp bool
		}{
			{[]string{}, "foo", false},
			{[]string{"foo", "bar"}, "foo", true},
			{[]string{"foo"}, "bar", false},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithNamespaces(Namespaces(tc.ns...)))
			got, err := f.Match(property.Namespace, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}

func TestFilterWeights(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			weights []float64
			val     float64
			exp     bool
		}{
			{[]float64{2.0}, 1.0, false},
			{[]float64{1.0}, 1.0, true},
			{[]float64{1.0, 2.0}, 1.0, true},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithWeights(Weights(tc.weights...)))
			got, err := f.Match(property.Weight, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}

func TestFilterAttrs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		a := memattrs.NewFromMap(map[string]string{"foo": "bar"})
		ax := memattrs.NewFromMap(map[string]string{"fooX": "barX", "foo": "bar"})
		ay := memattrs.NewFromMap(map[string]string{"bar": "foo", "car": "dac"})
		az := memattrs.NewFromMap(map[string]string{})

		testCases := []struct {
			a   attrs.Attrs
			val attrs.Attrs
			exp bool
		}{
			{ax, a, true},
			{ay, a, false},
			{az, a, false},
		}

		for _, tc := range testCases {
			f := MustMatcher(t, WithAttrs(Attrs(tc.a)))
			got, err := f.Match(property.Attrs, tc.val)
			if err != nil {
				t.Fatalf("failed to match: %v", err)
			}

			if got != tc.exp {
				t.Errorf("expected: %v, got: %v", tc.exp, got)
			}
		}
	})
}
