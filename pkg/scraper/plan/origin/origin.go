package origin

import "net/url"

// Origin is the origin of Space.
type Origin struct {
	url *url.URL
}

// New returns new Origin.
// s must be a valid URL, otherwise New returns error.
func New(s string) (*Origin, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	return &Origin{
		url: u,
	}, nil
}

// URL returns Origin URL.
func (s Origin) URL() *url.URL {
	return s.url
}
