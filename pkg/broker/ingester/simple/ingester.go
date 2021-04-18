package simple

import (
	"context"
	"encoding/json"

	"github.com/milosgajdos/netscrape/pkg/broker"
	"github.com/milosgajdos/netscrape/pkg/broker/ingester"
)

// Ingester digests data from broker.
type Ingester struct {
	opts ingester.Options
}

// NewIngester creates a new ingester and returns it
func NewIngester(opts ...ingester.Option) (*Ingester, error) {
	ropts := ingester.Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	return &Ingester{
		opts: ropts,
	}, nil
}

// Ingest ingests messages to the broker marshaled with the given marshaler.
func (in *Ingester) Ingest(ctx context.Context, b broker.Broker, topic string, msgType broker.Type, data interface{}, opts ...ingester.Option) error {
	ropts := ingester.Options{}
	for _, apply := range opts {
		apply(&ropts)
	}

	msg := &broker.Message{
		Type: msgType,
	}

	var err error

	m := ropts.Marshaler
	if m == nil {
		msg.Data, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Pub(ctx, topic, *msg)
	}

	msg.Data, err = m.Marshal(data)
	if err != nil {
		return err
	}

	return b.Pub(ctx, topic, *msg)
}
