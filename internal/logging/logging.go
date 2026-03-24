package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	vein "github.com/igorrius/go-vein"

	"github.com/igorrius/go-app-tui-template/internal/config"
)

type loggingState struct {
	writer *AsyncWriter
	cancel context.CancelFunc
	closer io.Closer // nil for stdout/stderr sinks
}

var (
	mu      sync.Mutex
	current *loggingState
)

// Init creates the VeinHandler, starts the AsyncWriter goroutine, and sets slog.Default.
// Calling Init again replaces any previously initialised logging state.
func Init(cfg config.LoggingConfig, bus *vein.Dispatcher) error {
	sinkW, closer, err := newSinkWriter(cfg)
	if err != nil {
		return fmt.Errorf("logging: %w", err)
	}

	slogH, err := newSlogHandler(cfg, sinkW)
	if err != nil {
		if closer != nil {
			_ = closer.Close()
		}
		return fmt.Errorf("logging: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	aw := newAsyncWriter(slogH, bus)
	aw.Start(ctx)

	mu.Lock()
	current = &loggingState{writer: aw, cancel: cancel, closer: closer}
	mu.Unlock()

	slog.SetDefault(slog.New(NewVeinHandler(cfg.Level, bus)))
	return nil
}

// Shutdown cancels the async writer, waits for it to drain remaining events, then closes
// the sink writer (if applicable). Respects ctx for a hard deadline on the drain wait.
func Shutdown(ctx context.Context) error {
	mu.Lock()
	s := current
	current = nil
	mu.Unlock()

	if s == nil {
		return nil
	}

	s.cancel()

	select {
	case <-s.writer.done:
	case <-ctx.Done():
		return ctx.Err()
	}

	if s.closer != nil {
		return s.closer.Close()
	}
	return nil
}
