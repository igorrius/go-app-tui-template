## Context

The application has no logging infrastructure at all. Every subsystem that needs observability currently must use `fmt.Println` or silently swallow errors. As features grow (TUI, SQLite, migrations, future services) the absence of structured, persisted logs will make debugging significantly harder.

The application already uses a message bus via `github.com/igorrius/go-vein`. Routing log records through the bus is the natural extension of the existing event-driven architecture and allows log sinks to be decoupled from the code that produces the records.

## Goals / Non-Goals

**Goals:**
- Provide a single, globally-available `*slog.Logger` that any package can use without importing anything beyond `log/slog`.
- Decouple log producers from log sinks via the `go-vein` message bus.
- Support two sinks: (1) stdout/stderr, (2) time-rotating files under a configurable directory.
- Support two text formats: `TEXT` (human-readable) and `JSON` (machine-readable).
- Make the logging subsystem configurable through the existing YAML config.
- Exclude the `logs/` directory from git and Docker automatically.

**Non-Goals:**
- Distributed tracing or metrics collection.
- Log shipping to remote systems (Elasticsearch, Loki, etc.) — can be added later as additional bus subscribers.
- Log sampling or rate-limiting.
- Per-package log levels.

## Decisions

### D1 — Route log records through the message bus

**Chosen**: `VeinHandler` implements `slog.Handler` and publishes a typed `LogEvent` to the `go-vein` bus on every `Handle` call.

**Rationale**: Keeps producer code (any package) fully decoupled from sink configuration. Sinks become subscribers and can be added/removed without touching producer code. Aligns with the existing architecture.

**Alternative considered**: Write directly from `slog.Handler` to file/stdout. Simpler but tightly couples the handler to sink details and makes it impossible to add sinks later without modifying the handler.

---

### D2 — Async subscriber writes to the configured sink

**Chosen**: A goroutine subscribes to the `go-vein` topic for `LogEvent` messages. It holds the configured `slog.Handler` (Text or JSON) pointed at the target writer (`os.Stdout` / `os.Stderr` / rotating file). On each received event it reconstructs a `slog.Record` and calls `handler.Handle`.

**Rationale**: Fully decouples I/O latency from the calling goroutine. The bus already provides buffering. The subscriber is started during application initialisation and shut down cleanly on teardown.

**Alternative considered**: Synchronous write in `VeinHandler.Handle`. Simpler, but blocks the calling goroutine on every I/O operation and couples the handler to the sink.

---

### D3 — Time-based log file rotation

**Chosen**: Use `github.com/lestrrat-go/file-rotatelogs` (or equivalent) for time-based rotation. Files are named with a timestamp pattern (e.g., `logs/app-2006-01-02_15.log`). The rotation interval is configurable (default 1 hour).

**Rationale**: Standard Go ecosystem library, actively maintained, supports `io.Writer` interface compatible with any `slog.Handler`.

**Alternative considered**: `gopkg.in/natefinish/lumberjack.v2` — size-based rather than time-based rotation; does not fit the requirement for time-interval rotation.

---

### D4 — Single global logger via `slog.SetDefault`

**Chosen**: The `logging.Init` function creates the `VeinHandler`, wires the async subscriber, then calls `slog.SetDefault(slog.New(handler))`. All code uses `slog.Info(...)` etc. from the standard library without any extra import.

**Rationale**: Standard Go pattern. Zero coupling between producers and the logging package after init.

---

### D5 — `internal/logging` package

**Chosen**: All logging types (`VeinHandler`, `LogEvent`, `AsyncWriter`, `Init`) live in `internal/logging/`.

**Rationale**: Keeps the public API minimal. Nothing outside the module can import it, which is appropriate for application-internal infrastructure.

---

### D6 — Configurable log level

**Chosen**: The `logging.level` config key controls the minimum level (`debug`, `info`, `warn`, `error`). `VeinHandler.Enabled` respects this level to avoid publishing events on the bus that will be discarded.

**Rationale**: Avoids unnecessary bus traffic when running in production with higher log levels.

## Risks / Trade-offs

- **[Risk] Bus backpressure** → the async writer goroutine may lag behind if I/O is slow (e.g., slow disk). The `go-vein` bus will buffer up to its channel capacity and then block or drop depending on configuration. Investigate bus capacity settings during implementation and document the chosen value.
- **[Risk] Lost logs on abrupt shutdown** → if the process is killed, buffered events in the bus are lost. Mitigation: call `logging.Shutdown()` in the application teardown path to drain the bus.
- **[Risk] Circular dependency** → if any package imported by `internal/logging` also imports `internal/logging`, a cycle results. Mitigation: `internal/logging` MUST NOT import any other `internal/` package except `internal/cfg`.
- **[Trade-off] Increased complexity** → two-stage log path (handler → bus → subscriber → sink) vs direct write is harder to understand. Justified by the decoupling benefit.

## Migration Plan

1. Add `github.com/igorrius/go-vein` and rotation library to `go.mod`.
2. Create `internal/logging/` package with `VeinHandler`, `LogEvent`, `AsyncWriter`, `Init`, `Shutdown`.
3. Extend `internal/cfg/config.go` with `LoggingConfig` struct.
4. Add the `logging` section to `app-config.dist.yaml`.
5. Wire `logging.Init` / `logging.Shutdown` into `cmd/go-app-tui-template/main.go`.
6. Add `logs/` to `.gitignore` and `.dockerignore`.
7. Verify no existing `fmt.Println` diagnostic calls need replacing (out of scope for initial implementation but worth noting).

No rollback complexity — logging is purely additive; removing it is a single revert of the wiring in `main.go`.

## Open Questions

- None — requirements are fully specified in the proposal.
