package logging_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	vein "github.com/igorrius/go-vein"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitShutdown_StdoutSink runs a full Init → log → Shutdown cycle and
// verifies that the log record appears in the configured sink.
func TestInitShutdown_StdoutSink(t *testing.T) {
	bus := &vein.Dispatcher{}

	// Capture output by redirecting through a buffer-backed text handler.
	var buf bytes.Buffer
	cfg := config.LoggingConfig{
		Level:  slog.LevelDebug,
		Format: "text",
		Sink:   "stdout",
	}

	// Build the slog handler backed by the buffer and an AsyncWriter manually
	// so we can capture output without touching os.Stdout.
	slogH, err := logging.NewSlogHandlerForTest(cfg, &buf)
	require.NoError(t, err)

	aw := logging.NewAsyncWriterForTest(slogH, bus)
	ctx, cancel := context.WithCancel(context.Background())
	aw.Start(ctx)

	veinH := logging.NewVeinHandler(cfg.Level, bus)
	slog.SetDefault(slog.New(veinH))

	slog.Info("integration-test-message", "key", "value")

	time.Sleep(50 * time.Millisecond)

	cancel()
	<-aw.Done()

	output := buf.String()
	assert.Contains(t, output, "integration-test-message")
	assert.Contains(t, output, "key=value")
}

// TestInit_InvalidSink verifies that Init returns an error for unknown sink values
// and does not alter slog.Default.
func TestInit_InvalidSink(t *testing.T) {
	original := slog.Default()
	bus := &vein.Dispatcher{}

	err := logging.Init(config.LoggingConfig{
		Level:  slog.LevelInfo,
		Format: "text",
		Sink:   "kafka",
	}, bus)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "kafka")
	assert.Same(t, original, slog.Default(), "slog.Default must not change on error")
}

// TestShutdown_AlreadyCancelledContext verifies that Shutdown returns
// ctx.Err() promptly when the context is already cancelled.
func TestShutdown_AlreadyCancelledContext(t *testing.T) {
	bus := &vein.Dispatcher{}
	require.NoError(t, logging.Init(config.LoggingConfig{
		Level:  slog.LevelInfo,
		Format: "text",
		Sink:   "stdout",
	}, bus))

	cancelled, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	err := logging.Shutdown(cancelled)
	assert.ErrorIs(t, err, context.Canceled)
}
