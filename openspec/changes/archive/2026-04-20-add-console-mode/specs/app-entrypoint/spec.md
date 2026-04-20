## MODIFIED Requirements

### Requirement: CLI entry point bootstraps the application
The system SHALL provide a CLI entry point using `urfave/cli/v3` that initializes the DI container, loads configuration, **initializes the logging subsystem**, resolves the `*vein.Dispatcher` from the DI container and injects it into the TUI `App` so that `LogView` can subscribe to `LogEvent` messages, and launches the TUI as the default action **unless `--console` is set, in which case it SHALL run in console mode**. On exit (including interrupt/SIGTERM) the application SHALL call `logging.Shutdown` to drain any buffered log events before the process terminates.

#### Scenario: Default launch starts TUI
- **WHEN** the user runs the binary without subcommands and without `--console`
- **THEN** the application SHALL initialize the DI container, load config, **initialize logging**, and start the bubbletea TUI in full-terminal mode

#### Scenario: Console flag starts console mode
- **WHEN** the user runs the binary with `--console`
- **THEN** the application SHALL initialize config, initialize logging to stdout, skip TUI initialization, and block until interrupted

#### Scenario: Version flag
- **WHEN** the user runs the binary with `--version`
- **THEN** the application SHALL print the version and exit

#### Scenario: Shutdown drains log events
- **WHEN** the application receives a shutdown signal or the TUI exits normally
- **THEN** `logging.Shutdown` SHALL be called and complete before the process exits

## MODIFIED Requirements

### Requirement: CLI flags and env variables are centralized
The system SHALL define all CLI flags and environment variable bindings in the `cmd/flag` package, including the `--console` flag with its `APP_CONSOLE` environment variable binding.

#### Scenario: Config file path override
- **WHEN** the user provides a `--config` flag with a file path
- **THEN** the application SHALL use the specified file instead of the default `app-config.yaml`

#### Scenario: Env variable binding
- **WHEN** an environment variable corresponding to a flag is set
- **THEN** the flag value SHALL be populated from the environment variable if not explicitly provided

#### Scenario: Console flag defined centrally
- **WHEN** the `--console` flag or `APP_CONSOLE` env var is used
- **THEN** the flag SHALL be defined in `cmd/flag` and registered on the root command
