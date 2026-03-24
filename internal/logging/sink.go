package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/igorrius/go-app-tui-template/internal/config"
)

// newSinkWriter returns an io.Writer (and optional io.Closer) for the configured sink.
// Returns a nil closer for stdout/stderr; the caller must check before closing.
func newSinkWriter(cfg config.LoggingConfig) (io.Writer, io.Closer, error) {
	switch cfg.Sink {
	case "stdout":
		return os.Stdout, nil, nil
	case "stderr":
		return os.Stderr, nil, nil
	case "file":
		if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {
			return nil, nil, fmt.Errorf("creating log directory %q: %w", cfg.Dir, err)
		}
		pattern := filepath.Join(cfg.Dir, "app-%Y-%m-%d_%H-%M.log")
		link := filepath.Join(cfg.Dir, "app.log")
		rl, err := rotatelogs.New(
			pattern,
			rotatelogs.WithLinkName(link),
			rotatelogs.WithRotationTime(cfg.RotationInterval),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("creating rotating log writer: %w", err)
		}
		return rl, rl, nil
	default:
		return nil, nil, fmt.Errorf("unknown log sink %q: expected stdout, stderr, or file", cfg.Sink)
	}
}

// newSlogHandler returns a slog.Handler writing to w in the format specified by cfg.
func newSlogHandler(cfg config.LoggingConfig, w io.Writer) (slog.Handler, error) {
	opts := &slog.HandlerOptions{Level: cfg.Level}
	switch cfg.Format {
	case "text":
		return slog.NewTextHandler(w, opts), nil
	case "json":
		return slog.NewJSONHandler(w, opts), nil
	default:
		return nil, fmt.Errorf("unknown log format %q: expected text or json", cfg.Format)
	}
}
