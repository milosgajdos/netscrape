package plan

import (
	"context"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/space/marshal"
	"github.com/milosgajdos/netscrape/pkg/space/origin"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
)

// NewMock creates mock Plan from given path and returns it.
func NewMock(path string) (*Plan, error) {
	s, err := origin.New(path)
	if err != nil {
		return nil, err
	}

	a, err := New(s)
	if err != nil {
		return nil, err
	}

	p, err := GetFilePathFromUrl(s.URL(), false)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var resources []marshal.Resource
	if err := yaml.Unmarshal(data, &resources); err != nil {
		return nil, err
	}

	for _, r := range resources {
		m, err := resource.New(
			r.Type,
			r.Name,
			r.Group,
			r.Version,
			r.Kind,
			r.Namespaced,
		)
		if err != nil {
			return nil, err
		}

		if err := a.Add(context.Background(), m); err != nil {
			return nil, err
		}
	}

	return a, nil
}
