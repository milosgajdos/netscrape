package netscrape

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/space/marshal/json"
)

// Ingester ingests data to store.
type Ingester struct {
	opts Options
}

// NewIngester creates a new ingester and returns it
func NewIngester(opts ...Option) (*Ingester, error) {
	ropts := Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	return &Ingester{
		opts: ropts,
	}, nil
}

// handle processes rmessages received by sub with handler h.
func (in *Ingester) handle(ctx context.Context, sub broker.Subscriber, h broker.Handler) error {
	for {
		if err := sub.Receive(ctx, h); err != nil {
			return err
		}
	}
}

// Ingest reads data from broker via sub and optionally stores it in a store.
// If no store is provided it dumps the received data to standard output.
// To store the data in the given store, the broker messages must be unmarshaled.
// If no marshaler is provided by default a JSON marshaler is created.
func (in *Ingester) Ingest(ctx context.Context, sub broker.Subscriber, opts ...Option) error {
	ropts := Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	m := ropts.Marshaler
	if m == nil {
		var err error
		m, err = json.NewMarshaler()
		if err != nil {
			return err
		}
	}

	if s := ropts.Store; s != nil {
		return in.handle(ctx, sub, storeHandler(ctx, s, m))
	}

	return in.handle(ctx, sub, dumpHandler)
}
