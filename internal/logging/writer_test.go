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

func TestAsyncWriter_TextFormat(t *testing.T) {
	bus := &vein.Dispatcher{}
	var buf bytes.Buffer

	cfg := config.LoggingConfig{Level: slog.LevelDebug, Format: "text", Sink: "stdout"}
	slogH, err := logging.NewSlogHandlerForTest(cfg, &buf)
	require.NoError(t, err)

	aw := logging.NewAsyncWriterForTest(slogH, bus)
	ctx, cancel := context.WithCancel(context.Background())
	aw.Start(ctx)

	vein.PublishTo[logging.LogEvent](bus, logging.LogEvent{
		Time:    time.Now(),
		Level:   slog.LevelInfo,
		Message: "text-test",
		Attrs:   nil,
	})

	// Allow event to be processed
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-aw.Done()

	assert.Contains(t, buf.String(), "text-test")
}

func TestAsyncWriter_JSONFormat(t *testing.T) {
	bus := &vein.Dispatcher{}
	var buf bytes.Buffer

	cfg := config.LoggingConfig{Level: slog.LevelDebug, Format: "json", Sink: "stdout"}
	slogH, err := logging.NewSlogHandlerForTest(cfg, &buf)
	require.NoError(t, err)

	aw := logging.NewAsyncWriterForTest(slogH, bus)
	ctx, cancel := context.WithCancel(context.Background())
	aw.Start(ctx)

	vein.PublishTo[logging.LogEvent](bus, logging.LogEvent{
		Time:    time.Now(),
		Level:   slog.LevelWarn,
		Message: "json-test",
		Attrs:   nil,
	})

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-aw.Done()

	output := buf.String()
	assert.Contains(t, output, `"msg":"json-test"`)
}

func TestAsyncWriter_DrainOnCancel(t *testing.T) {
	bus := &vein.Dispatcher{}
	var buf bytes.Buffer

	cfg := config.LoggingConfig{Level: slog.LevelDebug, Format: "text", Sink: "stdout"}
	slogH, err := logging.NewSlogHandlerForTest(cfg, &buf)
	require.NoError(t, err)

	aw := logging.NewAsyncWriterForTest(slogH, bus)
	ctx, cancel := context.WithCancel(context.Background())
	aw.Start(ctx)

	// Publish before cancel — should be drained
	for i := 0; i < 5; i++ {
		vein.PublishTo[logging.LogEvent](bus, logging.LogEvent{
			Time:    time.Now(),
			Level:   slog.LevelInfo,
			Message: "drain-event",
		})
	}
	cancel()
	<-aw.Done()

	assert.Contains(t, buf.String(), "drain-event")
}
