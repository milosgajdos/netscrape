package memory

import (
	"context"
	"sync"
	"time"

	"github.com/milosgajdos/netscrape/pkg/broker"
)

// Subscriber is broker subscriber.
type Subscriber struct {
	sync.RWMutex
	id     string
	topic  string
	active bool
	queue  queue
	ctl    chan<- sub
	exit   chan struct{}
}

// ID returns subscriber ID
func (s *Subscriber) ID(ctx context.Context) (string, error) {
	return s.id, nil
}

// Topic returns subscription topic.
func (s *Subscriber) Topic(ctx context.Context, opts ...broker.Option) (string, error) {
	return s.topic, nil
}

// Unsubscribe unsubscribes from the topic.
func (s *Subscriber) Unsubscribe(ctx context.Context, opts ...broker.Option) error {
	s.RLock()
	defer s.RUnlock()

	if s.active {
		close(s.exit)
		s.active = false
		// TODO(milosgajdos): should this really be async?
		go func(id, topic string) {
			select {
			case s.ctl <- sub{id: id, topic: topic}:
			case <-s.queue.exit:
			}
		}(s.id, s.topic)
	}

	return nil
}

// Receive processes received messages with handler
// NOTE: Receive is a blocking call!
func (s *Subscriber) Receive(ctx context.Context, h broker.Handler, opts ...broker.Option) error {
	ropts := broker.Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	recvTimeout := ropts.RecvTimeout
	if recvTimeout == 0 {
		recvTimeout = DefaultTimeout
	}

	s.Lock()
	if !s.active {
		s.Unlock()
		return broker.ErrSubscriptionInactive
	}
	s.Unlock()

	select {
	case <-ctx.Done():
		return nil
	case <-time.After(recvTimeout):
		return broker.ErrTimeout
	case <-s.queue.exit:
		return nil
	case <-s.exit:
		return nil
	case msg := <-s.queue.msg:
		if err := h(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}
