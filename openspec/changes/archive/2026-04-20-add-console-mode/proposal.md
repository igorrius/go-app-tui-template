## Why

When running the application in CI/CD pipelines, Docker containers, or any non-interactive environment, the full TUI is undesirable — it requires a TTY and produces unreadable output. A `--console` flag lets the app run as a plain process with structured log output to stdout, enabling integration with standard log aggregation tools.

## What Changes

- Add a `--console` CLI flag (boolean) to the root command
- When `--console` is set, skip TUI initialization and run the application in a headless console mode
- In console mode, all log events are written directly to stdout instead of the TUI log view
- The logging subsystem continues to operate normally; only the sink destination changes based on the mode
- The DI container wiring selects the appropriate runner (TUI vs console) based on the flag

## Capabilities

### New Capabilities

- `console-mode`: Defines the console execution mode — flag handling, headless application lifecycle, and log routing to stdout when TUI is disabled

### Modified Capabilities

- `app-entrypoint`: The CLI entry point must accept the `--console` flag and branch between TUI startup and console-mode startup

## Impact

- `cmd/flag/config.go`: Add `--console` flag definition
- `cmd/go-app-tui-template/main.go`: Branch on `--console` flag to start console runner instead of TUI
- `internal/cfg/`: DI wiring may need a console-mode provider or conditional binding for the log sink
- `internal/tui/app.go`: No changes; TUI is simply not started in console mode
- No breaking changes to existing TUI behavior
