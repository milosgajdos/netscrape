package memory

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"gonum.org/v1/gonum/graph/encoding"
)

// DOTAttrs returns attrs as []encoding.Attribute
func DOTAttrs(a attrs.Attrs) []encoding.Attribute {
	ctx := context.Background()
	keys, err := a.Keys(ctx)
	if err != nil {
		return nil
	}
	attrs := make([]encoding.Attribute, len(keys))

	i := 0
	for _, k := range keys {
		v, err := a.Get(context.Background(), k)
		if err != nil {
			return nil
		}
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: v,
		}
		i++
	}
	return attrs
}
