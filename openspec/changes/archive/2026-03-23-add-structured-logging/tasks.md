## 1. Dependencies & Project Setup

- [x] 1.1 Add `github.com/igorrius/go-vein` to `go.mod` / `go.sum` via `go get`
- [x] 1.2 Add a time-based log rotation library (e.g., `github.com/lestrrat-go/file-rotatelogs`) to `go.mod` / `go.sum`
- [x] 1.3 Add `logs/` to `.gitignore`
- [x] 1.4 Add `logs/` to `.dockerignore` (create the file if it does not exist)

## 2. Configuration

- [x] 2.1 Define `LoggingConfig` struct in `internal/cfg/config.go` with fields: `Level`, `Format`, `Sink`, `Dir`, `RotationInterval`
- [x] 2.2 Embed `LoggingConfig` into the root `Config` struct under the `logging` YAML key
- [x] 2.3 Register defaults in the config loader (`level: info`, `format: text`, `sink: stdout`, `dir: logs`, `rotation_interval: 1h`)
- [x] 2.4 Add the `logging` section to `app-config.dist.yaml` documenting all keys and defaults

## 3. LogEvent & VeinHandler

- [x] 3.1 Create `internal/logging/event.go` — define `LogEvent` struct (`Time`, `Level`, `Message`, `Attrs`)
- [x] 3.2 Create `internal/logging/handler.go` — implement `VeinHandler` struct satisfying `slog.Handler`
- [x] 3.3 Implement `VeinHandler.Enabled` — return `true` only when record level >= configured minimum level
- [x] 3.4 Implement `VeinHandler.Handle` — build `LogEvent` from `slog.Record` and publish to `go-vein` bus
- [x] 3.5 Implement `VeinHandler.WithAttrs` — return a new `VeinHandler` with accumulated attributes
- [x] 3.6 Implement `VeinHandler.WithGroup` — return a new `VeinHandler` with the group prefix applied

## 4. AsyncWriter

- [x] 4.1 Create `internal/logging/writer.go` — define `AsyncWriter` struct holding the sink `io.Writer` and the configured `slog.Handler`
- [x] 4.2 Implement `AsyncWriter.Start(ctx context.Context)` — subscribe to the `go-vein` bus topic, start processing goroutine
- [x] 4.3 Implement the event processing loop: receive `LogEvent`, reconstruct `slog.Record`, call `handler.Handle`
- [x] 4.4 Implement graceful shutdown: on context cancellation drain remaining events, then return
- [x] 4.5 Handle sink write errors non-fatally (log to stderr as fallback, continue)

## 5. Sink Factory

- [x] 5.1 Create `internal/logging/sink.go` — implement `newSinkWriter(cfg LoggingConfig) (io.Writer, io.Closer, error)`
- [x] 5.2 Handle `sink: stdout` — return `os.Stdout`
- [x] 5.3 Handle `sink: stderr` — return `os.Stderr`
- [x] 5.4 Handle `sink: file` — create log directory if absent, open rotating file writer with configured `dir` and `rotation_interval`, return the writer and its closer
- [x] 5.5 Return error for unrecognised `sink` value
- [x] 5.6 Implement `newSlogHandler(cfg LoggingConfig, w io.Writer) (slog.Handler, error)` — return `slog.NewTextHandler` or `slog.NewJSONHandler` based on `format`; return error for unknown format

## 6. Init & Shutdown

- [x] 6.1 Create `internal/logging/logging.go` — implement `Init(cfg LoggingConfig, bus *vein.Bus) error`
  - create sink writer, slog handler, `VeinHandler`, `AsyncWriter`
  - call `slog.SetDefault`
  - start `AsyncWriter`
- [x] 6.2 Implement `Shutdown(ctx context.Context) error` — cancel async writer context, wait for stop, close sink writer

## 7. Application Wiring

- [x] 7.1 In `cmd/go-app-tui-template/main.go` call `logging.Init(cfg.Logging, bus)` after config and bus are set up, propagate any error
- [x] 7.2 Ensure `logging.Shutdown` is deferred (or called in cleanup path) before the process exits
- [x] 7.3 Update the DI wiring in `internal/cfg/` if the bus is not yet available at startup (wire bus init before logging init)

## 8. Tests

- [x] 8.1 Write unit tests for `VeinHandler.Handle` — verify `LogEvent` fields match input `slog.Record`
- [x] 8.2 Write unit tests for `VeinHandler.Enabled` — verify level filtering
- [x] 8.3 Write unit tests for `AsyncWriter` — verify event is written to sink in correct format (text/json)
- [x] 8.4 Write unit tests for `newSinkWriter` — verify error on unknown sink and directory creation for `file` sink
- [x] 8.5 Write integration test for `Init`/`Shutdown` cycle — emit a log record and confirm it appears in the sink
