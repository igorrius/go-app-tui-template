## Why

The application currently has no structured logging infrastructure. Observability is essential for diagnosing issues in production; without it, there is no way to trace application behaviour, configuration issues, or runtime errors. Adding structured logging now establishes the observability foundation before the application grows more complex.

## What Changes

- Introduce a global `slog.Logger` initialised at startup via a dedicated logging package.
- Implement a custom `slog.Handler` (`VeinHandler`) that converts each log record into a typed event and publishes it to the application's internal message bus (`go-vein`).
- Register an asynchronous subscriber on the message bus that reads log events and writes them through a configurable standard `slog.TextHandler` or `slog.JSONHandler` to one of two configured sinks:
  1. **stdout / stderr**
  2. **Rotating log files** — written to a configurable directory (default `logs/`), rotated by time interval (default 1 hour).
- Extend the application configuration to include a `logging` section (sink type, format, log directory, rotation interval, log level).
- Ensure the `logs/` directory is excluded from version control (`.gitignore`) and Docker builds (`.dockerignore`).

## Capabilities

### New Capabilities

- `logging-init`: Global `slog` logger initialisation and wiring at application startup — creates the `VeinHandler`, subscribes the async writer, and sets `slog.SetDefault`.
- `logging-vein-handler`: Custom `slog.Handler` that serialises log records into message bus events and publishes them via `go-vein`.
- `logging-async-writer`: Asynchronous message bus subscriber that receives log events and writes them through a standard TEXT or JSON handler to the configured sink (stdout/stderr or rotating file).
- `logging-config`: Configuration schema for the logging subsystem (sink, format, directory, rotation interval, log level).
- `logging-file-rotation`: Time-based log file rotation producing timestamped files inside the configured log directory.

### Modified Capabilities

- `app-entrypoint`: Logger initialisation and teardown (flush/close) must be integrated into the application startup/shutdown lifecycle.
- `config-management`: A `logging` section must be added to the application configuration struct and YAML schema.

## Impact

- **New dependency**: `github.com/igorrius/go-vein` — message bus library.
- **New dependency**: a log-file rotation library (e.g., `gopkg.in/natefinish/lumberjack.v2` or a time-based equivalent such as `github.com/lestrrat-go/file-rotatelogs`).
- **Modified files**: `internal/cfg/config.go`, `cmd/go-app-tui-template/main.go`, `app-config.dist.yaml`.
- **New package**: `internal/logging/` — contains `VeinHandler`, async writer, and initialisation helpers.
- **Git / Docker**: `logs/` added to `.gitignore` and `.dockerignore`.
