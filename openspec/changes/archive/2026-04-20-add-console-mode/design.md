## Context

The application currently always starts the bubbletea TUI. Log events are routed through the `vein` event bus: `slog.Default` (a `VeinHandler`) publishes events to the bus, and an `AsyncWriter` consumes them and writes to the configured sink (file, stdout, or stderr). When running in CI/CD pipelines, Docker containers, or other non-TTY environments this design is unusable — the TUI requires a real terminal and its output is unreadable in log aggregators.

The `logging.Init` call sets `slog.Default` to a `VeinHandler` that feeds the vein bus. In console mode there is no TUI subscriber, so we need a direct stdout path for `slog`.

## Goals / Non-Goals

**Goals:**
- Add a `--console` boolean CLI flag to the root command
- When `--console` is set, skip TUI initialization entirely
- In console mode, route all `slog` output directly to stdout using the configured log format (text/json from config)
- In console mode, block until OS interrupt/SIGTERM so the process lifecycle is well-defined
- Preserve all existing TUI behavior when `--console` is not set

**Non-Goals:**
- Changing the log format or level selection logic
- Supporting a configurable console sink (always stdout in console mode)
- Sharing the vein bus event pipeline in console mode (not needed without TUI)
- Running application business logic in console mode (this is a scaffold; the app has no background services yet)

## Decisions

### Decision: `--console` flag lives in `cmd/flag/` alongside other flags
**Rationale**: All CLI flag definitions are centralized in `cmd/flag/`. Adding `ConsoleModeFlag` there keeps the pattern consistent and the flag discoverable from `flag.go`.

**Alternative considered**: Inline the flag definition in `main.go`. Rejected — breaks the existing convention.

### Decision: Console mode bypasses the VeinHandler and sets `slog.Default` to a direct stdout handler
**Rationale**: The `VeinHandler` exists to feed the TUI log view via the vein bus. In console mode there is no TUI subscriber, so pumping events through the bus adds complexity with no benefit. Setting `slog.Default` to `slog.NewTextHandler(os.Stdout, ...)` (or JSON, per config) is simpler and correct.

The existing `logging.Init` will not be called in console mode; instead a lightweight `logging.InitConsole(cfg)` function initialises only a direct stdout handler.

**Alternative considered**: Call `logging.Init` with a stdout sink and leave `slog.Default` as `VeinHandler`. This still creates and starts the `AsyncWriter` goroutine and vein bus unnecessarily. Rejected for simplicity.

### Decision: Console mode blocks on `signal.NotifyContext` (SIGINT/SIGTERM)
**Rationale**: After logging is initialised the process must stay alive until signalled. Using `signal.NotifyContext` is idiomatic Go and produces clean shutdown with the existing `logging.Shutdown` defer.

**Alternative considered**: `select{}` infinite block. Rejected — it ignores shutdown signals and never runs deferred cleanup.

### Decision: `main.go` branches between TUI and console runner in the `Action` handler
**Rationale**: The branch point is the CLI action. A simple `if cmd.Bool(flag.ConsoleModeFlag)` check keeps the two paths clearly separated without abstracting prematurely.

## Risks / Trade-offs

- [Risk] `logging.InitConsole` duplicates some logic from `logging.Init` → Mitigation: keep it small — only a handler + `slog.SetDefault`; extract shared helpers if they grow.
- [Risk] Future background services added to the app need console mode to keep them running → Mitigation: the signal-block loop will naturally extend to service start/stop hooks when they are introduced.

## Migration Plan

No migration needed. The change is purely additive — the `--console` flag defaults to `false`, leaving existing TUI behavior unchanged.
