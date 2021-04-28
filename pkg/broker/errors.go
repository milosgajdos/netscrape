package broker

import "errors"

var (
	// ErrNotImplemented is returned when requesting unimplemented functionality.
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrNotConnected is returned when publishing on broker when disconnected.
	ErrNotConnected = errors.New("ErrNotConnected")
	// ErrSubscriptionInactive is returned when attempting to receive on closed subscriber
	ErrSubscriptionInactive = errors.New("ErrSubscriptionInactive")
	// ErrTopicNotExist is returned when subscribing to a topic that does not exist
	ErrTopicNotExist = errors.New("ErrTopicNotExist")
	// ErrTimeout is returned when publish or subscribe operations timed out
	ErrTimeout = errors.New("ErrTimeout")
)
