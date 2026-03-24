## ADDED Requirements

### Requirement: Migration SQL files are embedded in the binary
The system SHALL embed all migration SQL files from `internal/infra/migrations/sql/` into the binary using `//go:embed`, so that no external migration files are needed at runtime.

#### Scenario: Binary carries migration files
- **WHEN** the binary is built
- **THEN** all `*.sql` files under `internal/infra/migrations/sql/` SHALL be accessible via the embedded filesystem at runtime

#### Scenario: New migration file is added during development
- **WHEN** a developer adds a new versioned SQL file (e.g., `00002_add_users.sql`) to `internal/infra/migrations/sql/`
- **THEN** a subsequent build SHALL include that file in the embedded filesystem automatically

### Requirement: Migration runner applies pending migrations upward
The system SHALL provide a `Runner` type in `internal/infra/migrations` that applies all pending migrations in version order when `Up` is called.

#### Scenario: No pending migrations
- **WHEN** `Runner.Up` is called and the database schema is already at the latest version
- **THEN** the runner SHALL return successfully without modifying the database

#### Scenario: One or more pending migrations
- **WHEN** `Runner.Up` is called and there are unapplied migrations
- **THEN** each migration SHALL be applied in ascending version order, and the goose version table SHALL be updated accordingly

#### Scenario: Migration SQL error causes rollback
- **WHEN** a migration SQL statement returns an error during `Runner.Up`
- **THEN** the failed migration's transaction SHALL be rolled back and the runner SHALL return the error without applying subsequent migrations

### Requirement: Migration runner rolls back the latest migration
The system SHALL provide a `Runner.Down` method that rolls back only the most recently applied migration.

#### Scenario: Roll back the latest migration
- **WHEN** `Runner.Down` is called
- **THEN** the most recently applied migration SHALL be reversed and the goose version table SHALL be updated

#### Scenario: No migrations to roll back
- **WHEN** `Runner.Down` is called but no migrations have been applied
- **THEN** the runner SHALL return successfully without error

### Requirement: Migration runner reports applied and pending migration status
The system SHALL provide a `Runner.Status` method that returns the list of migrations with their applied/pending state.

#### Scenario: Status lists all migrations
- **WHEN** `Runner.Status` is called
- **THEN** the returned list SHALL include every migration file with an indication of whether it has been applied

### Requirement: Migration runner is registered in the DI container
The system SHALL provide a `migrations.Runner` instance via the `samber/do/v2` DI container.

#### Scenario: DI resolves the migration runner
- **WHEN** a component requests `do.MustInvoke[*migrations.Runner](injector)`
- **THEN** the container SHALL return a runner connected to the application's `*sql.DB`

### Requirement: `migrate` CLI subcommand exposes `up`, `down`, and `status` operations
The system SHALL register a `migrate` top-level subcommand in the `go-app-tui-template` binary with three sub-subcommands: `up`, `down`, and `status`.

#### Scenario: `migrate up` applies pending migrations
- **WHEN** the user runs `go-app-tui-template migrate up`
- **THEN** all pending migrations SHALL be applied and the command SHALL exit 0 on success

#### Scenario: `migrate down` rolls back the latest migration
- **WHEN** the user runs `go-app-tui-template migrate down`
- **THEN** the most recently applied migration SHALL be rolled back and the command SHALL exit 0 on success

#### Scenario: `migrate status` prints migration state
- **WHEN** the user runs `go-app-tui-template migrate status`
- **THEN** the command SHALL print a table showing each migration version, name, and applied/pending state, then exit 0

#### Scenario: `migrate up` fails on SQL error
- **WHEN** the user runs `go-app-tui-template migrate up` and a migration returns an error
- **THEN** the command SHALL print the error to stderr and exit with a non-zero code
