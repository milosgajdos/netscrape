package memory

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/milosgajdos/netscrape/pkg/broker"
)

func TestSubscribe(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("Subscribe", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "fooTopic"

		sub1, err := b.Sub(context.Background(), topic)
		if err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		sub2, err := b.Sub(context.Background(), topic)
		if err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		sub1ID, err := sub1.ID(context.Background())
		if err != nil {
			t.Fatalf("failed to read sub ID: %v", err)
		}

		sub2ID, err := sub2.ID(context.Background())
		if err != nil {
			t.Fatalf("failed to read sub ID: %v", err)
		}

		if sub1ID == sub2ID {
			t.Fatal("subscriber IDs not unique")
		}

		if err := b.Close(); err != nil {
			t.Fatalf("failed closing session: %v", err)
		}
	})
}

func TestUnsubscribe(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("Unsubscribe", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "fooTopic"

		sub1, err := b.Sub(context.Background(), topic)
		if err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		if _, err := b.Sub(context.Background(), topic); err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		if err := sub1.Unsubscribe(context.Background()); err != nil {
			t.Errorf("failed to unsubscribe: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Fatalf("failed closing session: %v", err)
		}
	})

	t.Run("DoubleUnsubscribe", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "fooTopic"

		sub, err := b.Sub(context.Background(), topic)
		if err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		if err := sub.Unsubscribe(context.Background()); err != nil {
			t.Errorf("failed to unsubscribe: %v", err)
		}

		if err := sub.Unsubscribe(context.Background()); err != nil {
			t.Errorf("failed to unsubscribe: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Fatalf("failed closing session: %v", err)
		}
	})
}

func TestReceive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("ActiveReceive", func(t *testing.T) {
		b := MustBroker(t, broker.WithCap(1))

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "fooTopic"

		msg := broker.Message{
			UID:  "fooID",
			Data: []byte(`foo data`),
		}

		if err := b.Pub(context.Background(), topic, msg); err != nil {
			t.Fatalf("failed to publish message: %v", err)
		}

		sub, err := b.Sub(context.Background(), topic)
		if err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		errChan := make(chan error)

		h := func(ctx context.Context, m broker.Message) error {
			if m.UID != msg.UID {
				return fmt.Errorf("expected msg ID: %s, got: %s", msg.UID, m.UID)
			}
			return nil
		}

		go func() {
			errChan <- sub.Receive(context.Background(), h)
		}()

		if err := <-errChan; err != nil {
			t.Errorf("failed receiveing message: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("InactiveReceive", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "fooTopic"

		sub, err := b.Sub(context.Background(), topic)
		if err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		if err := sub.Unsubscribe(context.Background()); err != nil {
			t.Errorf("failed to unsubscribe: %v", err)
		}

		if err := sub.Receive(context.Background(), nil); !errors.Is(err, broker.ErrSubscriptionInactive) {
			t.Errorf("expected error: %v, got: %v", broker.ErrSubscriptionInactive, err)
		}

		if err := b.Close(); err != nil {
			t.Fatalf("failed closing session: %v", err)
		}
	})
}
