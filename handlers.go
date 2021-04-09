package netscrape

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/store"
)

// storeHandler returns broker handler that stores data decoded from received messages in s.
func storeHandler(ctx context.Context, s store.Store, m broker.Marshaler) broker.Handler {
	return func(ctx context.Context, msg broker.Message) error {
		switch msg.Type {
		case broker.Entity:
			var e space.Entity
			if err := m.Unmarshal(msg.Data, &e); err != nil {
				return err
			}
			return s.Add(ctx, e, store.WithUpsert(), store.WithAttrs(e.Attrs()))
		case broker.Object:
			var o space.Object
			if err := m.Unmarshal(msg.Data, &o); err != nil {
				return err
			}
			return s.Add(ctx, o, store.WithUpsert(), store.WithAttrs(o.Attrs()))
		case broker.Link:
			var l space.Link
			if err := m.Unmarshal(msg.Data, &l); err != nil {
				return err
			}

			err := s.Link(ctx, l.From(), l.To(), store.WithAttrs(l.Attrs()))
			if errors.Is(err, store.ErrEntityNotFound) {
				// NOTE: store.ErrEntityNotFound means either from or to entity do not exist in store
				// We will attempt to create partial entities and then attempt to link them again.
				from, err := entity.New(entity.Partial, entity.WithUID(l.From()))
				if err != nil {
					return err
				}

				if err := s.Add(ctx, from); err != nil && !errors.Is(err, store.ErrAlreadyExists) {
					return err
				}

				to, err := entity.New(entity.Partial, entity.WithUID(l.To()))
				if err != nil {
					return err
				}

				if err := s.Add(ctx, to); err != nil && !errors.Is(err, store.ErrAlreadyExists) {
					return err
				}
				return s.Link(ctx, l.From(), l.To(), store.WithAttrs(l.Attrs()))
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
