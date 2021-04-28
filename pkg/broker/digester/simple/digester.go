package simple

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/broker/digester"
	"github.com/milosgajdos/netscrape/pkg/broker/handlers"
)

// Digester digests data from broker.
type Digester struct {
	opts digester.Options
}

// NewDigester creates a new digester and returns it
func NewDigester(opts ...digester.Option) (*Digester, error) {
	ropts := digester.Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	return &Digester{
		opts: ropts,
	}, nil
}

// handle processes rmessages received by sub with handler h.
func (d *Digester) handle(ctx context.Context, sub broker.Subscriber, h broker.Handler) error {
	for {
		if err := sub.Receive(ctx, h); err != nil {
			return err
		}
	}
}

// Digest reads data from broker via sub and handles it via an optional handler.
// If no handler has been given, the message payload is copied to stdout.
func (d *Digester) Digest(ctx context.Context, sub broker.Subscriber, opts ...digester.Option) error {
	dopts := digester.Options{}
	for _, apply := range opts {
		apply(&dopts)
	}

	h := dopts.Handler
	if h == nil {
		h = handlers.DumpData
	}

	return d.handle(ctx, sub, h)
}
