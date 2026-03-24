package logging

import (
	"log/slog"
	"time"
)

// LogEvent is the message bus payload for a structured log record.
type LogEvent struct {
	Time    time.Time
	Level   slog.Level
	Message string
	Attrs   []slog.Attr
}
