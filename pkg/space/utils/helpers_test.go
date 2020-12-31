package utils

import (
	"net/url"
	"testing"
)

func TestGetFilePathFromUrl(t *testing.T) {
	tests := []struct {
		u   *url.URL
		abs bool
		exp string
		err bool
	}{
		{&url.URL{Scheme: "foo", Path: "/bar"}, false, "", true},
		{&url.URL{Scheme: "file", Path: "", Opaque: ""}, false, "", true},
		{&url.URL{Scheme: "file", Path: "", Opaque: "/bar"}, false, "bar", false},
		{&url.URL{Scheme: "file", Path: "/bar"}, true, "/bar", false},
	}

	for _, test := range tests {
		p, err := GetFilePathFromUrl(test.u, test.abs)
		if test.err {
			if err == nil {
				t.Errorf("expected error, got: %v", err)
				continue
			}
		}

		if p != test.exp {
			t.Errorf("expected: %v, got: %v", test.exp, p)
		}
	}
}
