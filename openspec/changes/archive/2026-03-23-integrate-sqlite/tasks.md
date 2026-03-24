## 1. Dependencies

- [x] 1.1 Add `modernc.org/sqlite` to `go.mod` as a direct dependency
- [x] 1.2 Add `pressly/goose/v3` to `go.mod` as a direct dependency
- [x] 1.3 Run `go mod tidy` and verify the build compiles cleanly

## 2. SQLite Storage Package

- [x] 2.1 Create `internal/infra/sqlite/sqlite.go` with an `Open(dsn string) (*sql.DB, error)` factory function that registers the `modernc.org/sqlite` driver and applies `PRAGMA journal_mode=WAL`, `PRAGMA busy_timeout=5000`, and `PRAGMA synchronous=NORMAL`
- [x] 2.2 Configure `database/sql` pool parameters (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`) in the factory function
- [x] 2.3 Ensure the SQLite driver import (`_ "modernc.org/sqlite"`) appears only in `internal/infra/sqlite/sqlite.go`
- [x] 2.4 Write a unit test for `Open` that verifies the three PRAGMAs are set on the returned connection

## 3. Configuration

- [x] 3.1 Add a `Database` sub-struct to the application config struct in `internal/config/config.go` with a `Path string` field (YAML key: `database.path`)
- [x] 3.2 Add `database.path` to `app-config.dist.yaml` with a sensible default (`db/go-app-tui-template.db`)
- [x] 3.3 Add `db/` to `.gitignore`
- [x] 3.4 Create `.dockerignore` (if it does not exist) and add `db/` to it

## 4. DI Provider for Database

- [x] 4.1 Create `internal/cfg/sqlite.go` with a `do` provider function that reads `config.Database.Path` and calls `sqlite.Open`, returning `*sql.DB`
- [x] 4.2 Register the `*sql.DB` provider in the DI injector setup in `internal/cfg/config.go` (or equivalent bootstrap file)

## 5. Migrations Package

- [x] 5.1 Create `internal/infra/migrations/sql/` directory and add the first migration file `00001_init.sql` with a `-- +goose Up` / `-- +goose Down` scaffolding comment (empty schema for now)
- [x] 5.2 Create `internal/infra/migrations/migrations.go` with:
  - An `//go:embed sql/*.sql` directive to embed migration files
  - A `Runner` struct holding a `*sql.DB` reference
  - `New(db *sql.DB) *Runner` constructor
  - `Up(ctx context.Context) error` method calling `goose.UpContext`
  - `Down(ctx context.Context) error` method calling `goose.DownContext`
  - `Status(ctx context.Context) error` method calling `goose.StatusContext`
- [x] 5.3 Write unit tests for `Runner.Up`, `Runner.Down`, and `Runner.Status` using a temporary SQLite file database

## 6. DI Provider for Migration Runner

- [x] 6.1 Create `internal/cfg/migrations.go` with a `do` provider function that resolves `*sql.DB` and returns `*migrations.Runner`
- [x] 6.2 Register the `*migrations.Runner` provider in the DI injector setup

## 7. `migrate` CLI Subcommand

- [x] 7.1 Create `cmd/flag/migrate.go` (or inline) defining the `migrate up`, `migrate down`, and `migrate status` sub-subcommand descriptors using `urfave/cli/v3`
- [x] 7.2 Implement the `migrate up` action: resolve `*migrations.Runner` from the DI container and call `Runner.Up`; print success message or write error to stderr and exit non-zero
- [x] 7.3 Implement the `migrate down` action: resolve `*migrations.Runner` and call `Runner.Down`; handle errors the same way
- [x] 7.4 Implement the `migrate status` action: resolve `*migrations.Runner` and call `Runner.Status`; output is printed by goose internally
- [x] 7.5 Register the top-level `migrate` command in `cmd/go-app-tui-template/main.go` alongside existing commands

## 8. Validation & Integration

- [x] 8.1 Build the binary (`go build ./...`) and verify no compilation errors
- [x] 8.2 Run `go-app-tui-template migrate status` against a fresh database and confirm it prints the migration list
- [x] 8.3 Run `go-app-tui-template migrate up` and confirm it applies `00001_init.sql` without error
- [x] 8.4 Run `go-app-tui-template migrate down` and confirm it rolls back the migration
- [x] 8.5 Run `go test ./internal/infra/...` and confirm all tests pass
- [x] 8.6 Verify `go-app-tui-template --help` lists `migrate` as a top-level command
