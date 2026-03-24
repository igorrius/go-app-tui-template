## ADDED Requirements

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

### Requirement: CLI flags and env variables are centralized
The system SHALL define all CLI flags and environment variable bindings in the `internal/cmd/flag` package.

#### Scenario: Config file path override
- **WHEN** the user provides a `--config` flag with a file path
- **THEN** the application SHALL use the specified file instead of the default `app-config.yaml`

#### Scenario: Env variable binding
- **WHEN** an environment variable corresponding to a flag is set
- **THEN** the flag value SHALL be populated from the environment variable if not explicitly provided

### Requirement: CLI entry point registers the `migrate` subcommand
The system SHALL register a `migrate` top-level command in the `urfave/cli/v3` application so it is discoverable via `go-app-tui-template --help`.

#### Scenario: `migrate` appears in help output
- **WHEN** the user runs `go-app-tui-template --help`
- **THEN** the help output SHALL list `migrate` as an available command with a short description

#### Scenario: Running `go-app-tui-template migrate` without a subcommand shows migrate help
- **WHEN** the user runs `go-app-tui-template migrate` with no further arguments
- **THEN** the application SHALL print the help for the `migrate` command listing `up`, `down`, and `status` subcommands
