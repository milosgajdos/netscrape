package memory

import (
	"context"
	"sync"
	"time"

	"github.com/milosgajdos/netscrape/pkg/broker"
	memuid "github.com/milosgajdos/netscrape/pkg/uuid/memory"
)

const (
	// DefaultSize is the default broker capacity.
	DefaultSize = 100
	// Defaultimeout is default timeout for both publish and subscribe ops.
	DefaultTimeout = 5 * time.Second
)

// sub is subscription.
type sub struct {
	id    string
	topic string
}

// queue is topic queue
type queue struct {
	msg  chan broker.Message
	exit chan struct{}
}

// Memory is in-memory broker.
type Memory struct {
	sync.RWMutex
	connected  bool
	size       int
	pubTimeout time.Duration
	subs       map[string]map[string]*Subscriber
	topics     map[string]queue
	exit       chan struct{}
	ctl        chan sub
}

// New crates a new in-memory broker and returns it.
func New(opts ...broker.Option) (*Memory, error) {
	bopts := broker.Options{}
	for _, apply := range opts {
		apply(&bopts)
	}

	size := bopts.Cap
	if size < 0 {
		size = DefaultSize
	}

	pubTimeout := bopts.PubTimeout
	if pubTimeout == 0 {
		pubTimeout = DefaultTimeout
	}

	return &Memory{
		connected:  false,
		size:       size,
		pubTimeout: pubTimeout,
		subs:       make(map[string]map[string]*Subscriber),
		topics:     make(map[string]queue),
	}, nil
}

// run starts the broker.
func (m *Memory) run(ctx context.Context) {
	defer m.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.exit:
			return
		case sub := <-m.ctl:
			m.Lock()
			if _, ok := m.subs[sub.topic]; ok {
				delete(m.subs[sub.topic], sub.id)
			}
			m.Unlock()
		}
	}
}

// Open opens a new broker session.
func (m *Memory) Open(ctx context.Context, opts ...broker.Option) error {
	m.RLock()
	if m.connected {
		m.RUnlock()
		return nil
	}
	m.RUnlock()

	m.Lock()
	defer m.Unlock()

	m.connected = true
	m.exit = make(chan struct{})
	m.ctl = make(chan sub, DefaultSize)

	// TODO: we should make sure Close() stops this goroutine
	// Try adding it to a sync.WaitGroup
	go m.run(ctx)

	return nil
}

func (m *Memory) publishMsg(ctx context.Context, q queue, msg broker.Message, timeout time.Duration) error {
	select {
	case <-ctx.Done():
		return m.Close()
	case <-m.exit:
		return nil
	case <-time.After(timeout):
		return broker.ErrTimeout
	case q.msg <- msg:
		return nil
	}
}

// Pub publishes m on the given topic.
// If the topic does not exist, it is automatically created.
// NOTE: Pub is a blocking call!
func (m *Memory) Pub(ctx context.Context, topic string, msg broker.Message, opts ...broker.Option) error {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return broker.ErrNotConnected
	}
	m.RUnlock()

	popts := broker.Options{}
	for _, apply := range opts {
		apply(&popts)
	}

	pubTimeout := popts.PubTimeout
	if pubTimeout == 0 {
		pubTimeout = m.pubTimeout
	}

	m.Lock()
	q, ok := m.topics[topic]
	if !ok {
		q = queue{
			msg:  make(chan broker.Message, m.size),
			exit: m.exit,
		}
		m.topics[topic] = q
	}
	m.Unlock()

	return m.publishMsg(ctx, q, msg, pubTimeout)
}

// Sub subscribes to the given topic.
func (m *Memory) Sub(ctx context.Context, topic string, opts ...broker.Option) (broker.Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, broker.ErrNotConnected
	}
	m.RUnlock()

	uid := memuid.New()

	m.Lock()
	defer m.Unlock()

	q, ok := m.topics[topic]
	if !ok {
		q = queue{
			msg:  make(chan broker.Message, m.size),
			exit: m.exit,
		}
		m.topics[topic] = q
	}

	sub := &Subscriber{
		id:     uid.String(),
		topic:  topic,
		active: true,
		queue:  q,
		ctl:    m.ctl,
		exit:   make(chan struct{}),
	}

	if m.subs[topic][sub.id] == nil {
		m.subs[topic] = make(map[string]*Subscriber)
	}
	m.subs[topic][sub.id] = sub

	return sub, nil
}

// Close closes broker session.
func (m *Memory) Close() error {
	m.Lock()
	defer m.Unlock()

	if !m.connected {
		return nil
	}

	close(m.exit)

	for topic, subs := range m.subs {
		for id := range subs {
			delete(m.subs[topic], id)
		}
		delete(m.subs, topic)
		delete(m.topics, topic)
	}

	m.connected = false

	return nil
}
