## 1. Logging Text Formatter

- [x] 1.1 Create `internal/logging/formatter.go` with `TextFormatter` struct and `Format(LogEvent) string` method
- [x] 1.2 Implement field order: Time (`time.DateTime` layout) → Level (uppercased) → Module (from `"module"` attr, default `"APP"`, always uppercase) → remaining attrs as `key=value`
- [x] 1.3 Ensure `"module"` attribute is excluded from the trailing key=value output
- [x] 1.4 Write unit tests in `internal/logging/formatter_test.go` covering: full record with module, missing module defaults to APP, module uppercasing, module not duplicated in attrs, time format

## 2. LogView Screen

- [x] 2.1 Create `internal/tui/screen/logview.go` with `LogView` struct holding `*vein.Dispatcher`, `logging.TextFormatter`, line buffer `[]string`, `width`, `height`
- [x] 2.2 Implement `LogView.Init()`: emit `slog.Info("Subscription to log was performed", "module", "LOG_VIEW")`, return `vein.Subscribe[LogEvent]` tea.Cmd
- [x] 2.3 Implement `LogView.Update()`: handle `tea.WindowSizeMsg` (update dimensions, trim buffer) and `logging.LogEvent` (format and append to buffer, trim to height)
- [x] 2.4 Implement `LogView.View()`: render the last `height` lines as plain text, no scroll controls
- [x] 2.5 Add `LogViewID ID = "logview"` constant to `internal/tui/screen/screen.go`
- [x] 2.6 Register `LogViewID` in the `screen.New` factory, passing a `*vein.Dispatcher` argument to `NewLogView`

## 3. TUI App Router and KeyBar

- [x] 3.1 Add `{key: "F9", label: "Logs"}` entry to `keyBarItems` in `internal/tui/component/keybar.go`
- [x] 3.2 Add `case "f9":` branch in `App.Update` key handling in `internal/tui/app.go` that calls `a.routeTo(screen.LogViewID)`
- [x] 3.3 Update `tui.NewApp` signature to accept `*vein.Dispatcher` and pass it through to `screen.New(screen.LogViewID, dispatcher)` (or store it on `App` for use during routing)

## 4. Dependency Injection Wiring

- [x] 4.1 Update `cfg/tui.go` `ProvideTUIApp` to resolve `*vein.Dispatcher` via `do.Invoke` and pass it to `tui.NewApp`
- [x] 4.2 Update `screen.New` factory signature to accept an optional `*vein.Dispatcher` parameter (or move dispatcher storage to the `App` struct and supply it when constructing `LogView` in `routeTo`)

## 5. Validation

- [x] 5.1 Run `go build ./...` and confirm no compilation errors
- [x] 5.2 Run `go test ./...` and confirm all tests pass including new formatter tests
- [x] 5.3 Run `golangci-lint run` and confirm no lint violations
- [x] 5.4 Manual smoke test: start the app, press F9, verify log view appears and shows the subscription log entry
