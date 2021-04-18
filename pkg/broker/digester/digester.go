package digester

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/broker"
)

// Digester digests broker messages.
type Digester interface {
	// Digest handles messages received from the broker using the provided handler.
	Digest(context.Context, broker.Subscriber, ...Option) error
}
