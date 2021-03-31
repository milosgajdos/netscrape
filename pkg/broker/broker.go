package broker

import "context"

// Broker provides a message broker.
type Broker interface {
	// Pub publishes messages on the given topic.
	Pub(ctx context.Context, topic string, m Message, opts ...Option) error
	// Sub creates a new subscriber to the given topic.
	Sub(ctx context.Context, topic string, opts ...Option) (Subscriber, error)
}

// Handler processes messages published on topic.
type Handler func(context.Context, Message) error

// Subscriber processes messages from broker.
type Subscriber interface {
	// ID returns subscription ID.
	ID(context.Context) (string, error)
	// Topic returns subscription topic.
	Topic(context.Context, ...Option) (string, error)
	// Unsubscribe unsubscribes from the topic.
	Unsubscribe(context.Context, ...Option) error
	// Receive processes received messages with handler.
	Receive(context.Context, Handler, ...Option) error
}

// Message is broker message.
type Message struct {
	// UID is unique message ID.
	UID string
	// Data contains message payload.
	Data []byte
	// Attrs are message attributes.
	Attrs map[string]string
}
