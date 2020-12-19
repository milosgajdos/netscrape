package plan

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/origin"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/space/types"
	"github.com/milosgajdos/netscrape/pkg/space/utils"
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

	p, err := utils.GetFilePathFromUrl(s.URL(), false)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var resources []types.Resource
	if err := yaml.Unmarshal(data, &resources); err != nil {
		return nil, err
	}

	for _, r := range resources {
		m, err := resource.New(
			r.Name,
			r.Kind,
			r.Group,
			r.Version,
			r.Namespaced,
		)
		if err != nil {
			return nil, err
		}

		if err := a.Add(m, space.AddOptions{}); err != nil {
			return nil, err
		}
	}

	return a, nil
}