package signal

import (
	"context"
	"os"
	"os/signal"
)

// Handle will create a notification for a terminate signal and cancel the
// returned context in that case, and will stop listening upon the given
// context being canceled.
func Handle(ctx context.Context) context.Context {
	sigCtx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			return
		}
	}()

	return sigCtx
}
