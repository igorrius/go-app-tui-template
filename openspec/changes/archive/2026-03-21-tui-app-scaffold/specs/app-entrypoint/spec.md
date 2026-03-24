## ADDED Requirements

### Requirement: CLI entry point bootstraps the application
The system SHALL provide a CLI entry point using `urfave/cli/v3` that initializes the DI container, loads configuration, and launches the TUI as the default action.

#### Scenario: Default launch starts TUI
- **WHEN** the user runs the binary without subcommands
- **THEN** the application SHALL initialize the DI container, load config, and start the bubbletea TUI in full-terminal mode

#### Scenario: Version flag
- **WHEN** the user runs the binary with `--version`
- **THEN** the application SHALL print the version and exit

### Requirement: CLI flags and env variables are centralized
The system SHALL define all CLI flags and environment variable bindings in the `internal/cmd/flag` package.

#### Scenario: Config file path override
- **WHEN** the user provides a `--config` flag with a file path
- **THEN** the application SHALL use the specified file instead of the default `app-config.yaml`

#### Scenario: Env variable binding
- **WHEN** an environment variable corresponding to a flag is set
- **THEN** the flag value SHALL be populated from the environment variable if not explicitly provided
