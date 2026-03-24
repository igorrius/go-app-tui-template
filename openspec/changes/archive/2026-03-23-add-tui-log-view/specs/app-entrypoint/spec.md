## MODIFIED Requirements

### Requirement: CLI entry point bootstraps the application
The system SHALL provide a CLI entry point using `urfave/cli/v3` that initializes the DI container, loads configuration, **initializes the logging subsystem**, resolves the `*vein.Dispatcher` from the DI container and injects it into the TUI `App` so that `LogView` can subscribe to `LogEvent` messages, and launches the TUI as the default action. On exit (including interrupt/SIGTERM) the application SHALL call `logging.Shutdown` to drain any buffered log events before the process terminates.

#### Scenario: Default launch starts TUI
- **WHEN** the user runs the binary without subcommands
- **THEN** the application SHALL initialize the DI container, load config, **initialize logging**, and start the bubbletea TUI in full-terminal mode

#### Scenario: Version flag
- **WHEN** the user runs the binary with `--version`
- **THEN** the application SHALL print the version and exit

#### Scenario: Shutdown drains log events
- **WHEN** the application receives a shutdown signal or the TUI exits normally
- **THEN** `logging.Shutdown` SHALL be called and complete before the process exits
