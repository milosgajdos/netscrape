package memory

import (
	"fmt"
	"net/url"
	"strings"
)

// GetFilePathFromUrl returns a path from the path in u.
// NOTE: this is a mad hack to work around URL shenanigans
func GetFilePathFromUrl(u *url.URL, abs bool) (string, error) {
	if u.Scheme != "file" {
		return "", fmt.Errorf("unexpected URL scheme: %v", u.Scheme)
	}

	p := u.Path

	// if the path is empty, we may have been given opaque URL
	// see: https://golang.org/pkg/net/url/#URL
	if p == "" {
		// check if u is an Opaque URL
		if p = u.Opaque; p == "" {
			return "", fmt.Errorf("empty URL path")
		}
	}

	if !abs {
		return strings.TrimPrefix(p, "/"), nil
	}

	return p, nil
}
