package ingester

import (
	"context"

	"github.com/milosgajdos/netscrape/pkg/broker"
)

// Ingester ingests messages to the broker.
type Ingester interface {
	// Ingest ingests data to the broker on the given topic.
	Ingest(ctx context.Context, b broker.Broker, topic string, msgType broker.Type, data interface{}, opts ...Option) error
}
