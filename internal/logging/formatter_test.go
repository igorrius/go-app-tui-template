package logging_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestColorTextFormatter_Basic(t *testing.T) {
	// Force color profile for testing to ensure ANSI codes are generated
	lipgloss.SetColorProfile(termenv.TrueColor)
	ts := time.Date(2026, 3, 23, 14, 5, 0, 0, time.UTC)
	e := makeEvent(ts, slog.LevelInfo, "msg",
		slog.String("module", "db"),
		slog.String("query", "select"),
	)

	f := logging.ColorTextFormatter{}
	got := f.Format(e)

	// Check for presence of ANSI escape codes (starts with \x1b[)
	assert.Contains(t, got, "\x1b[")
	// Check for the content itself (it might be interleaved with codes)
	assert.Contains(t, got, "2026-03-23 14:05:00")
	assert.Contains(t, got, "INFO")
	assert.Contains(t, got, "DB")
	assert.Contains(t, got, "msg")
	assert.Contains(t, got, "query")
	assert.Contains(t, got, "select")
}

func TestColorTextFormatter_Levels(t *testing.T) {
	// Force color profile for testing to ensure ANSI codes are generated
	lipgloss.SetColorProfile(termenv.TrueColor)
	ts := time.Date(2026, 3, 23, 14, 5, 0, 0, time.UTC)
	f := logging.ColorTextFormatter{}

	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}
	for _, l := range levels {
		e := makeEvent(ts, l, "msg")
		got := f.Format(e)
		assert.Contains(t, got, l.String())
		assert.Contains(t, got, "\x1b[")
	}
}

func makeEvent(t time.Time, level slog.Level, msg string, attrs ...slog.Attr) logging.LogEvent {
	return logging.LogEvent{
		Time:    t,
		Level:   level,
		Message: msg,
		Attrs:   attrs,
	}
}

func TestTextFormatter_FullRecordWithModule(t *testing.T) {
	ts := time.Date(2026, 3, 23, 14, 5, 0, 0, time.UTC)
	e := makeEvent(ts, slog.LevelInfo, "msg",
		slog.String("module", "db"),
		slog.String("query", "select"),
	)

	f := logging.TextFormatter{}
	got := f.Format(e)

	assert.Equal(t, "[2026-03-23 14:05:00] INFO DB msg query=select", got)
}

func TestTextFormatter_MissingModuleDefaultsToAPP(t *testing.T) {
	ts := time.Date(2026, 3, 23, 14, 5, 0, 0, time.UTC)
	e := makeEvent(ts, slog.LevelWarn, "msg")

	f := logging.TextFormatter{}
	got := f.Format(e)

	assert.Equal(t, "[2026-03-23 14:05:00] WARN APP msg", got)
}

func TestTextFormatter_ModuleIsUppercased(t *testing.T) {
	ts := time.Date(2026, 3, 23, 14, 5, 0, 0, time.UTC)

	f := logging.TextFormatter{}

	e := makeEvent(ts, slog.LevelInfo, "msg", slog.String("module", "log_view"))
	assert.Equal(t, "[2026-03-23 14:05:00] INFO LOG_VIEW msg", f.Format(e))

	e = makeEvent(ts, slog.LevelInfo, "msg", slog.String("module", "LogView"))
	assert.Equal(t, "[2026-03-23 14:05:00] INFO LOGVIEW msg", f.Format(e))
}

func TestTextFormatter_ModuleNotDuplicatedInAttrs(t *testing.T) {
	ts := time.Date(2026, 3, 23, 14, 5, 0, 0, time.UTC)
	e := makeEvent(ts, slog.LevelInfo, "msg",
		slog.String("module", "db"),
		slog.String("query", "select"),
	)

	f := logging.TextFormatter{}
	got := f.Format(e)

	assert.Contains(t, got, "DB")
	assert.Contains(t, got, "msg")
	assert.Contains(t, got, "query=select")
	assert.NotContains(t, got, "module=db")
}

func TestTextFormatter_TimeFormat(t *testing.T) {
	ts := time.Date(2026, 3, 23, 14, 5, 0, 0, time.UTC)
	e := makeEvent(ts, slog.LevelInfo, "msg")

	f := logging.TextFormatter{}
	got := f.Format(e)

	assert.Equal(t, "[2026-03-23 14:05:00]", got[:21])
}
