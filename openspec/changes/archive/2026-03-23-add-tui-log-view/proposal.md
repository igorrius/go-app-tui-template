## Why

The application currently has no way to observe structured log output directly within the TUI. Adding a dedicated log view accessible via F9 allows developers and operators to monitor real-time log activity without leaving the terminal UI, improving observability and debuggability.

## What Changes

- Add a new `LogView` screen (`screen.LogViewID`) that displays live log entries inside the TUI.
- Register F9 key binding in the root `App` to toggle/navigate to the log view.
- Add F9 entry to the `KeyBar` component.
- Implement a custom log text formatter that outputs fields in order: `Time`, `Level`, `Module` (from `"module"` attribute, defaulting to `"APP"` in uppercase), then remaining attributes.
- Subscribe the `LogView` to `LogEvent` messages via the existing message bus so it receives all structured log records.
- The log view renders entries without scroll support (latest entries visible, no scrolling).
- On activation, the `LogView` emits a first info-level log: `"Subscription to log was performed"` with `module: LOG_VIEW`.

## Capabilities

### New Capabilities

- `tui-log-view`: A TUI screen that subscribes to `LogEvent` and renders live log entries with a custom formatter (Time → Level → Module → remaining attrs). Activated via F9. No scrolling.
- `logging-text-formatter`: A custom `slog`-compatible text formatter that formats log records as `Time Level Module [attrs...]`.

### Modified Capabilities

- `tui-layout`: F9 key binding added to the key bar and the root app router.
- `app-entrypoint`: The log view must be wired into the dependency-injection / dispatcher so it receives `LogEvent` from the message bus.

## Impact

- `internal/tui/screen/`: new `logview.go` file; `screen.go` registry updated with `LogViewID`.
- `internal/tui/app.go`: F9 route added.
- `internal/tui/component/keybar.go`: F9 item added.
- `internal/logging/`: new `formatter.go` file with the custom text formatter.
- `internal/cfg/tui.go` (or dispatcher): wire `LogView` as a subscriber to `LogEvent`.
- No new external dependencies required.
