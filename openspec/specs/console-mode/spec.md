# Console Mode

## Purpose

Define the console mode capability that allows the application to run without the TUI, logging directly to stdout.

## Requirements

### Requirement: Console mode is activated by the --console flag
The system SHALL support a `--console` boolean CLI flag on the root command that, when set, runs the application in console mode instead of starting the TUI.

#### Scenario: Default launch is TUI mode
- **WHEN** the user runs the binary without `--console`
- **THEN** the application SHALL start in TUI mode as before

#### Scenario: Console flag activates console mode
- **WHEN** the user runs the binary with `--console`
- **THEN** the application SHALL skip TUI initialisation and run in console mode

#### Scenario: Console flag via environment variable
- **WHEN** the environment variable `APP_CONSOLE` is set to a truthy value and `--console` is not explicitly provided
- **THEN** the application SHALL activate console mode

### Requirement: Console mode initialises logging to stdout
In console mode the system SHALL initialise `slog.Default` with a handler writing directly to `os.Stdout`, using the log level and format (text/json) from the application config. The vein event bus and `AsyncWriter` SHALL NOT be started in console mode.

#### Scenario: Log output goes to stdout in console mode
- **WHEN** the application runs with `--console` and code calls `slog.Info(...)`
- **THEN** the formatted log line SHALL appear on stdout

#### Scenario: Log format respects config
- **WHEN** the application config specifies `format: json` and `--console` is set
- **THEN** log output SHALL be JSON-formatted on stdout

### Requirement: Console mode blocks until interrupted
In console mode the system SHALL block the main goroutine until SIGINT or SIGTERM is received, then perform a clean shutdown (including `logging.Shutdown`).

#### Scenario: Process stays alive until signal
- **WHEN** the application is running with `--console`
- **THEN** the process SHALL remain running until SIGINT or SIGTERM is received

#### Scenario: Shutdown drains logs
- **WHEN** SIGINT or SIGTERM is received in console mode
- **THEN** `logging.Shutdown` SHALL be called and complete before the process exits
