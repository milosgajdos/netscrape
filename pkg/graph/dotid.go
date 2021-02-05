package graph

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/space"
)

// DOTIDFromEntity returns GraphViz DOT ID for given space.Entity.
// NOTE: the returned DOTIDFromEntity follows the below naming convention:
// resourceGroup/resourceVersion/resourceKind/entityNamespace/entityName
func DOTIDFromEntity(o space.Entity) (string, error) {
	if o.Resource() == nil {
		return "", ErrMissingResource
	}

	return strings.Join([]string{
		o.Resource().Group(),
		o.Resource().Version(),
		o.Resource().Kind(),
		o.Namespace(),
		o.Name()}, "/"), nil
}
