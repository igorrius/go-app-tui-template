## Purpose

Provides a production-ready SQLite connection factory using `modernc.org/sqlite` (pure-Go, CGO-free). Opens connections with WAL mode and concurrent-access PRAGMAs, configures the `database/sql` connection pool, and registers the `*sql.DB` in the DI container.

## Requirements

### Requirement: SQLite connection opens with WAL mode and concurrent-access PRAGMAs
The system SHALL open a SQLite database connection with `journal_mode=WAL`, `busy_timeout=5000`, and `synchronous=NORMAL` applied immediately after the connection is established.

#### Scenario: Connection is opened for the first time
- **WHEN** the application initialises the SQLite provider
- **THEN** the database file SHALL be created if it does not exist, and all three PRAGMAs SHALL be set before any query is executed

#### Scenario: Concurrent readers do not block each other
- **WHEN** multiple goroutines execute read queries simultaneously
- **THEN** all reads SHALL complete without returning `SQLITE_BUSY` errors

#### Scenario: Writer contention is retried automatically
- **WHEN** a write operation encounters a locked database
- **THEN** the driver SHALL retry for up to 5000 ms before returning an error

### Requirement: SQLite package is isolated under `internal/infra/sqlite`
The system SHALL expose the SQLite factory only through the `internal/infra/sqlite` package; no other package SHALL import a SQLite driver directly.

#### Scenario: Driver import is centralised
- **WHEN** the application is built
- **THEN** the SQLite driver registration (`_ "modernc.org/sqlite"`) SHALL appear only in `internal/infra/sqlite/sqlite.go`

### Requirement: Database connection is registered in the DI container
The system SHALL provide a `*sql.DB` instance via the `samber/do/v2` DI container so that any service can declare it as a dependency.

#### Scenario: DI provider resolves a *sql.DB
- **WHEN** a service requests `do.MustInvoke[*sql.DB](injector)`
- **THEN** the DI container SHALL return the configured connection without error

#### Scenario: Database path is configurable
- **WHEN** the application configuration specifies a `database.path` value
- **THEN** the SQLite factory SHALL open the file at that path

### Requirement: Connection pool is tuned for concurrent access
The system SHALL configure `database/sql` pool parameters appropriate for SQLite's single-writer model.

#### Scenario: Pool is initialised with sensible defaults
- **WHEN** the SQLite provider initialises
- **THEN** `SetMaxOpenConns`, `SetMaxIdleConns`, and `SetConnMaxLifetime` SHALL be set to values documented in the package

### Requirement: Database file is stored under a `db/` subfolder
The system SHALL default the database file path to `db/go-app-tui-template.db` (relative to the working directory). The `db/` directory SHALL be excluded from version control and Docker build context.

#### Scenario: Default database path uses `db/` subfolder
- **WHEN** no `database.path` is specified in the configuration
- **THEN** the application SHALL use `db/go-app-tui-template.db` as the database file path, creating the `db/` directory if it does not exist

#### Scenario: `db/` directory is excluded from version control
- **WHEN** the repository is cloned or `git status` is run
- **THEN** the `db/` directory and its contents SHALL be listed in `.gitignore` and not tracked by git

#### Scenario: `db/` directory is excluded from Docker build context
- **WHEN** a Docker image is built
- **THEN** the `db/` directory SHALL be listed in `.dockerignore` and not copied into the image
