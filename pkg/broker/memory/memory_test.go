package memory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/milosgajdos/netscrape/pkg/broker"
)

func MustBroker(t *testing.T, opts ...broker.Option) *Memory {
	b, err := New(opts...)
	if err != nil {
		t.Fatalf("failed creating broker: %v", err)
	}

	return b
}

func TestOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("Open", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		if !b.connected {
			t.Errorf("expected broker session to be connected")
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("DoubleOpen", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		if !b.connected {
			t.Errorf("expected broker session to be connected")
		}

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})
}

func TestPub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("PublishOnNewTopic", func(t *testing.T) {
		b := MustBroker(t, broker.WithCap(1))

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "foo"
		msg := broker.Message{
			UID:  "fooID",
			Data: []byte(`foo data`),
		}

		if err := b.Pub(context.Background(), topic, msg); err != nil {
			t.Fatalf("failed to publish message: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("PublishOnExistingTopic", func(t *testing.T) {
		// NOTE: we need broker with buffer for 2 messages
		// so the two publish calls below don't block
		b := MustBroker(t, broker.WithCap(2))

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "foo"
		msg := broker.Message{
			UID:  "fooID",
			Data: []byte(`foo data`),
		}

		// NOTE: first publish on non-existent topic creates the topic
		if err := b.Pub(context.Background(), topic, msg); err != nil {
			t.Fatalf("failed to publish message: %v", err)
		}

		if err := b.Pub(context.Background(), topic, msg); err != nil {
			t.Fatalf("failed to publish message: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("PublishTimeout", func(t *testing.T) {
		// NOTE: we want zero sized broker
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "foo"
		msg := broker.Message{
			UID:  "fooID",
			Data: []byte(`foo data`),
		}

		pt := 100 * time.Millisecond

		if err := b.Pub(context.Background(), topic, msg, broker.WithPubTimeout(pt)); !errors.Is(err, broker.ErrTimeout) {
			t.Fatalf("expected error: %v, got: %v", broker.ErrTimeout, err)
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("PublishConcurrentlyOnExistingTopic", func(t *testing.T) {
		count := 2

		// NOTE: we need broker with buffer for 2 messages
		// so the two publish calls below don't block
		b := MustBroker(t, broker.WithCap(count))

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "foo"
		msg := broker.Message{
			UID:  "fooID",
			Data: []byte(`foo data`),
		}

		errChan := make(chan error, count)

		for i := 0; i < count; i++ {
			go func(i int) {
				errChan <- b.Pub(context.Background(), topic, msg)
			}(i)
		}

		for i := 0; i < count; i++ {
			if err := <-errChan; err != nil {
				t.Errorf("failed publishing msg: %v", err)
			}
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("PublishNotConnected", func(t *testing.T) {
		b := MustBroker(t)

		topic := "foo"
		msg := broker.Message{
			UID:  "fooID",
			Data: []byte(`foo data`),
		}

		if err := b.Pub(context.Background(), topic, msg); !errors.Is(err, broker.ErrNotConnected) {
			t.Fatalf("expected error: %v, got: %v", broker.ErrNotConnected, err)
		}
	})
}

func TestSub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("SubscribeToNewTopic", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		topic := "foo"

		// first subscriber
		sub1, err := b.Sub(context.Background(), topic)
		if err != nil {
			t.Fatalf("failed to subscribe to topic %s: %v", topic, err)
		}

		tp, err := sub1.Topic(context.Background())
		if err != nil {
			t.Fatalf("failed to read the topic: %v", err)
		}

		if tp != topic {
			t.Errorf("expected topic: %s, tot: %s", topic, tp)
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("SubscribeToExistingTopic", func(t *testing.T) {
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

		tp, err := sub.Topic(context.Background())
		if err != nil {
			t.Fatalf("failed to read the topic: %v", err)
		}

		if tp != topic {
			t.Errorf("expected topic: %s, tot: %s", topic, tp)
		}

		if err := b.Close(); err != nil {
			t.Errorf("failed to close broker session: %v", err)
		}
	})

	t.Run("SubscribeNotConnected", func(t *testing.T) {
		b := MustBroker(t)

		topic := "foo"

		if _, err := b.Sub(context.Background(), topic); !errors.Is(err, broker.ErrNotConnected) {
			t.Errorf("expected error: %v, got: %v", broker.ErrNotConnected, err)
		}
	})
}

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Fatalf("failed closing session: %v", err)
		}
	})

	t.Run("DoubleClose", func(t *testing.T) {
		b := MustBroker(t)

		if err := b.Open(context.Background()); err != nil {
			t.Fatalf("failed to open broker session: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Fatalf("failed closing session: %v", err)
		}

		if err := b.Close(); err != nil {
			t.Fatalf("failed closing session: %v", err)
		}
	})
}
