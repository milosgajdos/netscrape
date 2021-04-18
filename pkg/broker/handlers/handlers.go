package handlers

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/milosgajdos/netscrape/pkg/broker"
)

// DumpData dumps message data to stdout.
func DumpData(ctx context.Context, m broker.Message) error {
	if _, err := io.Copy(os.Stdout, bytes.NewReader(m.Data)); err != nil {
		return err
	}
	return nil
}
