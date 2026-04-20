## 1. CLI Flag

- [x] 1.1 Add `ConsoleModeFlag` (`--console`, bool) to `cmd/flag/config.go` with `APP_CONSOLE` env var binding
- [x] 1.2 Register `flag.ConsoleModeFlag` in the root command's `Flags` slice in `cmd/go-app-tui-template/main.go`

## 2. Console Logging Initialisation

- [x] 2.1 Add `InitConsole(cfg config.LoggingConfig) error` to `internal/logging/` that sets `slog.Default` to a direct stdout handler (text or json per config) without starting the vein bus or AsyncWriter
- [x] 2.2 Write a unit test for `InitConsole` verifying that a `slog.Info` call produces output on stdout and that the format matches the config

## 3. Console Mode Runner

- [x] 3.1 Add a `runConsoleMode` function in `cmd/go-app-tui-template/main.go` that: calls `logging.InitConsole`, defers `logging.Shutdown`, then blocks on `signal.NotifyContext` until SIGINT or SIGTERM
- [x] 3.2 Branch in the root `Action` handler: if `cmd.Bool(flag.ConsoleModeFlag)` is true, call `runConsoleMode`; otherwise continue with the existing TUI path

## 4. Verification

- [x] 4.1 Run the binary with `--console` and confirm it starts, logs to stdout, and exits cleanly on Ctrl-C
- [x] 4.2 Run the binary without `--console` and confirm TUI still starts normally
- [x] 4.3 Run all existing tests to confirm no regressions (`go test ./...`)
