## ADDED Requirements

### Requirement: CLI entry point registers the `migrate` subcommand
The system SHALL register a `migrate` top-level command in the `urfave/cli/v3` application so it is discoverable via `go-app-tui-template --help`.

#### Scenario: `migrate` appears in help output
- **WHEN** the user runs `go-app-tui-template --help`
- **THEN** the help output SHALL list `migrate` as an available command with a short description

#### Scenario: Running `go-app-tui-template migrate` without a subcommand shows migrate help
- **WHEN** the user runs `go-app-tui-template migrate` with no further arguments
- **THEN** the application SHALL print the help for the `migrate` command listing `up`, `down`, and `status` subcommands
