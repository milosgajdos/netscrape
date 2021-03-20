package memory

import (
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"gonum.org/v1/gonum/graph/encoding"
)

// DOTAttrs returns attrs as []encoding.Attribute
func DOTAttrs(a attrs.Attrs) []encoding.Attribute {
	keys := a.Keys()
	attrs := make([]encoding.Attribute, len(keys))

	i := 0
	for _, k := range keys {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: a.Get(k),
		}
		i++
	}

	return attrs
}
