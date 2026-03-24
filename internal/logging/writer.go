package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	vein "github.com/igorrius/go-vein"
)

// AsyncWriter subscribes to the go-vein bus for LogEvent messages and writes them
// to the configured sink via a standard slog handler.
type AsyncWriter struct {
	handler slog.Handler
	bus     *vein.Dispatcher
	done    chan struct{}
}

func newAsyncWriter(handler slog.Handler, bus *vein.Dispatcher) *AsyncWriter {
	return &AsyncWriter{
		handler: handler,
		bus:     bus,
		done:    make(chan struct{}),
	}
}

// Start subscribes to LogEvent on the bus and begins the processing goroutine.
// The goroutine runs until ctx is cancelled; it drains remaining queued events before stopping.
func (w *AsyncWriter) Start(ctx context.Context) {
	sub := vein.SubscribeTo[LogEvent](w.bus)
	ch := sub.OnC()

	go func() {
		defer close(w.done)
		defer sub.Unsubscribe()

		for {
			select {
			case event := <-ch:
				w.writeEvent(event)
			case <-ctx.Done():
				for {
					select {
					case event := <-ch:
						w.writeEvent(event)
					default:
						return
					}
				}
			}
		}
	}()
}

func (w *AsyncWriter) writeEvent(event LogEvent) {
	r := slog.NewRecord(event.Time, event.Level, event.Message, 0)
	r.AddAttrs(event.Attrs...)
	if err := w.handler.Handle(context.Background(), r); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: sink write error: %v\n", err)
	}
}
