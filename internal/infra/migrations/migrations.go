package migrations

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

// Runner wraps a *sql.DB and applies goose migrations from the embedded SQL directory.
type Runner struct {
	db *sql.DB
}

// New creates a Runner for the given database connection.
func New(db *sql.DB) *Runner {
	goose.SetBaseFS(embedMigrations)
	_ = goose.SetDialect("sqlite")
	return &Runner{db: db}
}

// Up applies all pending migrations.
func (r *Runner) Up(ctx context.Context) error {
	return goose.UpContext(ctx, r.db, "sql")
}

// Down rolls back the most recently applied migration.
func (r *Runner) Down(ctx context.Context) error {
	return goose.DownContext(ctx, r.db, "sql")
}

// Status prints the current migration status to stdout.
func (r *Runner) Status(ctx context.Context) error {
	return goose.StatusContext(ctx, r.db, "sql")
}
