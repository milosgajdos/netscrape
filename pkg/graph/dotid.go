package graph

import (
	"strings"

	"github.com/milosgajdos/netscrape/pkg/space"
)

// DOTID returns GraphViz DOT ID for given space.Object.
// NOTE: the returned DOTID follows the below naming convention:
// resourceGroup/resourceVersion/resourceKind/objectNamespace/objectName
func DOTID(o space.Object) (string, error) {
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