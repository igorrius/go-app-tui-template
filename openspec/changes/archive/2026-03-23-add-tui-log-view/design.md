## Context

The TUI currently routes between screens via the root `App` model in `internal/tui/app.go`. Screens implement `tea.Model` and are registered in `internal/tui/screen/screen.go`. The `KeyBar` component in `internal/tui/component/keybar.go` lists function-key shortcuts.

Structured logs are published to the `go-vein` Dispatcher as `LogEvent` messages by `VeinHandler`. The async writer (`AsyncWriter`) in `internal/logging/writer.go` subscribes to the bus and forwards events to the configured sink (file / stdout). The `App` is wired through `samber/do` in `internal/cfg/tui.go`, which currently injects only `*config.AppConfig`.

## Goals / Non-Goals

**Goals:**
- Add a `LogView` screen that displays live log entries as plain text lines inside the TUI.
- Activate `LogView` on F9; show it in the key bar.
- Subscribe `LogView` to `LogEvent` via the vein Dispatcher so it mirrors all structured logs.
- Implement a custom text formatter (`logging.TextFormatter`) that renders each log record as: `<Time> <LEVEL> <MODULE> [key=value ...]` where `Module` is extracted from the `"module"` attribute (uppercased) or defaults to `"APP"`.
- On activation (first subscription), emit `slog.Info("Subscription to log was performed", "module", "LOG_VIEW")`.
- No scroll support — the view shows the most recent entries that fit the current height.

**Non-Goals:**
- Interactive filtering, search, or keyboard navigation within the log view.
- Colour / lipgloss styling of individual log levels (plain text output only).
- Persistence or export of the in-view log buffer.
- Changes to the existing file/stdout sink pipeline.

## Decisions

### 1. LogView receives LogEvent via Bubbletea message dispatch (not direct channel)

`go-vein` subscribers return a `tea.Cmd` (channel-based). Wrapping the vein subscription as a `tea.Cmd` means `LogEvent` messages flow through the normal Bubble Tea `Update` cycle, keeping the model mutation single-threaded and race-free.

**Alternative considered**: a background goroutine writing directly into a shared slice — rejected because it requires a mutex and breaks the Elm architecture invariant.

### 2. LogView subscription fires on Init()

`LogView.Init()` returns the vein subscription `tea.Cmd`. This is the idiomatic Bubble Tea hook for starting side-effects. The info log `"Subscription to log was performed"` is emitted immediately before returning the cmd.

### 3. LogView is injected with *vein.Dispatcher via ProvideTUIApp

`ProvideTUIApp` in `internal/cfg/tui.go` already uses `do.Invoke`. It will additionally resolve `*vein.Dispatcher` and pass it to `tui.NewApp`. `NewApp` passes the dispatcher down to `LogView` at construction time.

**Alternative considered**: resolve the dispatcher lazily inside `screen.New` — rejected because it would require threading the injector through the screen registry, adding undesirable coupling.

### 4. Custom TextFormatter as a standalone type in internal/logging

`logging.TextFormatter` is a pure function-object with a single exported method `Format(LogEvent) string`. It has no dependency on `slog.Handler` internals and can be unit-tested independently. The `LogView` holds a `TextFormatter` value and calls it inside `Update` when a `LogEvent` arrives.

### 5. No-scroll display: fixed-size ring of the last N lines

The view tracks a circular slice of the last `height` rendered lines. On every resize or new event the slice is trimmed to fit. This gives O(1) rendering with no scroll state.

## Risks / Trade-offs

- **High log volume**: If the application emits logs at very high frequency, every `LogEvent` triggers a Bubble Tea re-render. Mitigation: the vein handler already throttles via the async pipeline; the view stores only the last N lines so memory stays bounded.
- **Module attribute missing from existing log calls**: Existing log sites do not pass `module=…`; they will all render as `APP`. This is acceptable short-term — module tagging is an opt-in convention for future log sites.
- **Dispatcher injected through NewApp signature change**: Adding `*vein.Dispatcher` to `tui.NewApp` is a minor API break inside the monorepo but has only one call site (`cfg/tui.go`), so the impact is minimal.

## Migration Plan

1. Add `logging.TextFormatter` with unit tests.
2. Add `screen.LogViewID` and `LogView` model; update `screen.New`.
3. Update `tui.NewApp` to accept `*vein.Dispatcher`; pass it to `LogView`.
4. Update `cfg/tui.go` to resolve and inject the dispatcher.
5. Add F9 entry to `KeyBar` items.
6. Add F9 case to `App.Update` router.
7. Run `golangci-lint && go test ./...`.

Rollback: all changes are additive except the `NewApp` signature; reverting the signature and the cfg provider restores the previous state fully.

## Open Questions

- Should the log view be the default screen on startup, or always start on Dashboard? → Defaulting to Dashboard (F1) is consistent with current behaviour; F9 activates the log view on demand.
