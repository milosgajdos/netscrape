package netscrape

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/marshal"
	"github.com/milosgajdos/netscrape/pkg/store"
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

// storeHandler returns broker handler that stores data decoded from received messages in s.
func storeHandler(ctx context.Context, s store.Store) broker.Handler {
	return func(ctx context.Context, m broker.Message) error {
		switch m.Type {
		case broker.Entity:
			var e marshal.Entity
			if err := marshal.Unmarshal(marshal.JSON, m.Data, &e); err != nil {
				return err
			}

			ent, err := marshal.ToSpaceEntity(e)
			if err != nil {
				return err
			}

			if err := s.Add(ctx, ent, store.WithUpsert(), store.WithAttrs(ent.Attrs())); err != nil {
				return err
			}
		case broker.Link:
			var l marshal.Link
			if err := marshal.Unmarshal(marshal.JSON, m.Data, &l); err != nil {
				return err
			}

			lnk, err := marshal.ToSpaceLink(l)
			if err != nil {
				return err
			}

			err = s.Link(ctx, lnk.From(), lnk.To(), store.WithAttrs(lnk.Attrs()))
			if errors.Is(err, store.ErrEntityNotFound) {
				// NOTE: store.ErrEntityNotFound means either from or to entity do not exist in store
				// We will attempt to create partial entities and then attempt to link them again.
				fromEnt, err := entity.NewPartial(entity.WithUID(lnk.From()))
				if err != nil {
					return err
				}

				if err := s.Add(ctx, fromEnt); err != nil && !errors.Is(err, store.ErrAlreadyExists) {
					return err
				}

				toEnt, err := entity.NewPartial(entity.WithUID(lnk.To()))
				if err != nil {
					return err
				}

				if err := s.Add(ctx, toEnt); err != nil && !errors.Is(err, store.ErrAlreadyExists) {
					return err
				}

				return s.Link(ctx, lnk.From(), lnk.To(), store.WithAttrs(lnk.Attrs()))
			}

			if err != nil {
				return err
			}
		}
		return nil
	}
}

// dumpHandler dumps message to stdout.
func dumpHandler(ctx context.Context, m broker.Message) error {
	if _, err := io.Copy(os.Stdout, bytes.NewReader(m.Data)); err != nil {
		return err
	}
	return nil
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
func (in *Ingester) Ingest(ctx context.Context, sub broker.Subscriber, opts ...Option) error {
	ropts := Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	if s := ropts.Store; s != nil {
		return in.handle(ctx, sub, storeHandler(ctx, s))
	}

	return in.handle(ctx, sub, dumpHandler)
}
