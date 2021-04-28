package broker

import "time"

// Sink is broker sink strategy.
type Sink int

const (
	// FanIn delivers message to at most one subscriber.
	FanIn Sink = iota
	// FanOut fans out message to all subscribers.
	FanOut
)

// Options configure broker.
type Options struct {
	// Cap configures broker capacity.
	Cap int
	// Sink configures publish strategy.
	Sink Sink
	// PubTimeout configures publish timeout.
	PubTimeout time.Duration
	// RecvTimeout configures receive timeout.
	RecvTimeout time.Duration
}

// Option is functional broker option.
type Option func(*Options)

// WithCap sets Cap option
func WithCap(c int) Option {
	return func(o *Options) {
		o.Cap = c
	}
}

// WithSink configures Sink options
func WithSink(s Sink) Option {
	return func(o *Options) {
		o.Sink = s
	}
}

// WithPubTimeout configures PubTimeout option
func WithPubTimeout(p time.Duration) Option {
	return func(o *Options) {
		o.PubTimeout = p
	}
}

// WithSubTimeout configures RecvTimeout option
func WithSubTimeout(s time.Duration) Option {
	return func(o *Options) {
		o.RecvTimeout = s
	}
}
