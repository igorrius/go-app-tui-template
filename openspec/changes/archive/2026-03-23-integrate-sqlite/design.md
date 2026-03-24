## Context

The project currently has no persistent storage layer. The `internal/infra/` package is scaffolded but empty. The application already uses `samber/do/v2` for dependency injection and `urfave/cli/v3` for CLI commands. Future usage anticipates high-concurrency read/write access to the SQLite database (many goroutines reading, occasional write bursts), so the driver choice and connection tuning must be deliberate from the start.

## Goals / Non-Goals

**Goals:**
- Provide a production-ready SQLite connection factory with WAL mode enabled and sensible PRAGMA defaults for concurrent workloads
- Implement a self-contained migration runner using embedded SQL files so the binary carries its own schema management
- Expose a `migrate` CLI subcommand (`up`, `down`, `status`) for operational use
- Register database and migrator providers in the DI container so all future services can depend on `*sql.DB` cleanly

**Non-Goals:**
- ORM or query-builder layer (raw `database/sql` is sufficient)
- Distributed or networked database support
- Automatic migration on application startup (migrations are explicit operator action)
- Multiple database backends or abstraction behind an interface (SQLite only for now)

## Decisions

### Decision 1: Use `modernc.org/sqlite` as the SQLite driver

**Choice**: `modernc.org/sqlite` (pure-Go, CGO-free transpiled driver)

**Rationale**: The project has no CGO dependencies today. Introducing `mattn/go-sqlite3` (the most common alternative) would require a C toolchain in CI and cross-compilation environments. `modernc.org/sqlite` is a pure-Go transpilation of the SQLite C source, supports the full SQLite feature set including WAL mode, and keeps the build simple. Performance is within ~10–20% of the native CGO driver for typical workloads, which is acceptable.

**Alternatives considered**:
- `mattn/go-sqlite3`: Best raw performance, CGO required — rejected to keep build portable
- `zombiezen.com/go/sqlite`: Low-level CGO bindings with excellent concurrency control — rejected for same CGO reason; revisit if performance profiling shows `modernc` as a bottleneck
- `glebarez/go-sqlite`: Another pure-Go driver based on `modernc` — redundant; use `modernc` directly

### Decision 2: Use `pressly/goose/v3` for migrations

**Choice**: `pressly/goose/v3` with embedded SQL files via `//go:embed`

**Rationale**: Goose supports embedded filesystems natively (`goose.SetBaseFS`), sequential versioned SQL files, and exposes a clean programmatic API suitable for both CLI and testing. It outputs minimal dependencies compared to `golang-migrate/migrate/v4` (which has a large driver matrix). Goose's `VersionedMigration` model and `goose.Up(db, dir)` API are straightforward to wrap.

**Alternatives considered**:
- `golang-migrate/migrate/v4`: More feature-rich but heavier dependency tree; overkill for embedded SQL files
- Hand-rolled migrations: Brittle, no rollback support — rejected

### Decision 3: PRAGMA tuning for concurrent access

**WAL mode**: `PRAGMA journal_mode=WAL` — enables concurrent readers while a writer holds its lock, eliminating the "database is locked" error under read-heavy load.

**Busy timeout**: `PRAGMA busy_timeout=5000` — instructs SQLite to retry for up to 5 seconds before returning `SQLITE_BUSY`, smoothing write contention spikes.

**Synchronous**: `PRAGMA synchronous=NORMAL` — safest mode compatible with WAL that avoids fsync on every write; provides acceptable durability (only risks data loss on OS crash, not application crash).

**Connection pool**: Use `db.SetMaxOpenConns(1)` for the **write** connection to serialise writes at the `database/sql` level, preventing `SQLITE_BUSY` from concurrent writers. A separate read pool with `db.SetMaxOpenConns(N)` (N = runtime.NumCPU) can be exposed for read-heavy queries. For the initial implementation, a single pool with WAL is sufficient; read/write pool separation is noted as a future optimisation.

### Decision 4: Package layout

```
internal/
  infra/
    sqlite/
      sqlite.go          // Open() → *sql.DB factory + PRAGMA setup
    migrations/
      migrations.go      // Runner wrapping goose
      sql/               // Embedded *.sql migration files
        00001_init.sql
        ...
cmd/
  go-app-tui-template/
    main.go              // registers migrate subcommand
internal/
  cfg/
    sqlite.go            // DI provider: *sql.DB
    migrations.go        // DI provider: migrations.Runner
```

The `migrate` subcommand is registered alongside existing commands in `cmd/go-app-tui-template/main.go` using `urfave/cli/v3`'s `Commands` field.

## Risks / Trade-offs

- **WAL file left on disk**: WAL mode creates a `-wal` and `-shm` file next to the database file. Operators must be aware that all three files constitute the database. → Mitigation: document this in the `migrate` subcommand help text and README.
- **Single-writer bottleneck**: Even with WAL, only one writer can hold the write lock at a time. Under extreme write concurrency, `busy_timeout` may still expire. → Mitigation: `SetMaxOpenConns(1)` on the write path serialises at the Go layer before hitting SQLite, preventing cascading busy errors.
- **`modernc` performance ceiling**: Pure-Go interpreting the SQLite bytecode is measurably slower than the CGO driver under sustained heavy writes. → Mitigation: acceptable for current scale; documented as a known limitation with a clear upgrade path to `zombiezen.com/go/sqlite` if profiling demands it.
- **No rollback by default**: `goose down` rolls back one migration; bulk rollbacks are destructive on production data. → Mitigation: CLI `down` subcommand requires explicit invocation; no automatic rollback on startup.

## Migration Plan

1. Add `modernc.org/sqlite` and `pressly/goose/v3` to `go.mod` / `go.sum`
2. Implement `internal/infra/sqlite` package and DI provider; default DSN is `db/go-app-tui-template.db` (directory created at runtime if absent)
3. Implement `internal/infra/migrations` package with embedded SQL directory
4. Register DI providers in `internal/cfg/`
5. Add `migrate` CLI subcommand in `cmd/go-app-tui-template/main.go`
6. Add `db/` to `.gitignore`; create `.dockerignore` (if absent) and add `db/` to it
7. Write unit tests for the migration runner (using an in-memory `file::memory:?cache=shared` SQLite URL)

**Rollback**: The feature is purely additive. Removing the `migrate` subcommand and unregistering DI providers is sufficient to revert; the database file is only created on first `migrate up`.

## Open Questions

- Should `migrate up` be automatically run on application startup behind a feature flag, or always remain manual? (Current decision: always manual.)
- Is a `file::memory:?cache=shared` DSN acceptable for integration tests, or should tests use a temp-file database? (Prefer temp-file to test WAL mode behaviour accurately.)
