package logging_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/igorrius/go-app-tui-template/internal/logging"
	vein "github.com/igorrius/go-vein"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVeinHandler_Enabled(t *testing.T) {
	bus := &vein.Dispatcher{}
	h := logging.NewVeinHandler(slog.LevelWarn, bus)

	assert.False(t, h.Enabled(context.Background(), slog.LevelDebug))
	assert.False(t, h.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, h.Enabled(context.Background(), slog.LevelWarn))
	assert.True(t, h.Enabled(context.Background(), slog.LevelError))
}

func TestVeinHandler_Handle_PublishesEvent(t *testing.T) {
	bus := &vein.Dispatcher{}

	received := make(chan logging.LogEvent, 1)
	sub := vein.SubscribeTo[logging.LogEvent](bus)
	sub.On(func(e logging.LogEvent) { received <- e })

	h := logging.NewVeinHandler(slog.LevelDebug, bus)

	ts := time.Now().Truncate(time.Millisecond)
	r := slog.NewRecord(ts, slog.LevelInfo, "hello", 0)
	r.AddAttrs(slog.String("key", "val"))

	require.NoError(t, h.Handle(context.Background(), r))

	select {
	case e := <-received:
		assert.Equal(t, slog.LevelInfo, e.Level)
		assert.Equal(t, "hello", e.Message)
		assert.Equal(t, ts, e.Time)
		require.Len(t, e.Attrs, 1)
		assert.Equal(t, "key", e.Attrs[0].Key)
		assert.Equal(t, "val", e.Attrs[0].Value.String())
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for log event")
	}
}

func TestVeinHandler_WithAttrs(t *testing.T) {
	bus := &vein.Dispatcher{}
	received := make(chan logging.LogEvent, 1)
	sub := vein.SubscribeTo[logging.LogEvent](bus)
	sub.On(func(e logging.LogEvent) { received <- e })

	h := logging.NewVeinHandler(slog.LevelDebug, bus)
	h2 := h.WithAttrs([]slog.Attr{slog.String("persistent", "yes")})

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	require.NoError(t, h2.Handle(context.Background(), r))

	select {
	case e := <-received:
		require.Len(t, e.Attrs, 1)
		assert.Equal(t, "persistent", e.Attrs[0].Key)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for log event")
	}
}

func TestVeinHandler_WithGroup(t *testing.T) {
	bus := &vein.Dispatcher{}
	received := make(chan logging.LogEvent, 1)
	sub := vein.SubscribeTo[logging.LogEvent](bus)
	sub.On(func(e logging.LogEvent) { received <- e })

	h := logging.NewVeinHandler(slog.LevelDebug, bus)
	h2 := h.WithGroup("req").WithAttrs([]slog.Attr{slog.String("id", "123")})

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	require.NoError(t, h2.Handle(context.Background(), r))

	select {
	case e := <-received:
		require.Len(t, e.Attrs, 1)
		assert.Equal(t, "req", e.Attrs[0].Key)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for log event")
	}
}
