## Why

The application needs persistent local storage to support data-driven features. SQLite is the right fit for an embedded, single-node tool, but future usage requires handling high-concurrency reads and writes safely — so the integration must be designed with WAL mode, proper connection pooling, and a clean migration strategy from day one.

## What Changes

- Add a dedicated `internal/infra/sqlite` package providing a configured SQLite connection pool with WAL mode enabled for concurrent access
- Add a `internal/infra/migrations` package for managing schema migrations programmatically
- Add a `migrate` CLI subcommand (under the existing `go-app-tui-template` binary) exposing `up`, `down`, and `status` operations
- Register the database and migration runner in the DI container (`samber/do/v2`)

## Capabilities

### New Capabilities

- `sqlite-storage`: SQLite connection package with WAL mode, tuned connection pool (`PRAGMA journal_mode=WAL`, `PRAGMA busy_timeout`, `PRAGMA synchronous=NORMAL`) and a factory registered in the DI container
- `database-migrations`: Embedded SQL migration files (`internal/infra/migrations/sql/`), a migration runner using `pressly/goose` (or `golang-migrate/migrate`), and a `migrate` CLI subcommand with `up`, `down`, and `status` sub-commands

### Modified Capabilities

- `app-entrypoint`: The CLI entrypoint gains a new top-level `migrate` subcommand

## Impact

- **New dependency**: `modernc.org/sqlite` (pure-Go, CGO-free SQLite driver) for portability with WAL mode support; `pressly/goose/v3` for embedded SQL migration management
- **DI container**: New providers for `*sql.DB` and the migration runner added to `internal/cfg/`
- **CLI**: `cmd/go-app-tui-template/main.go` gains the `migrate` subcommand registered via `urfave/cli/v3`
- **No breaking changes** to existing TUI, config, or startup flow
