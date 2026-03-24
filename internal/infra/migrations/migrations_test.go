package migrations_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/igorrius/go-app-tui-template/internal/infra/migrations"
	"github.com/igorrius/go-app-tui-template/internal/infra/sqlite"
	"github.com/stretchr/testify/require"
)

func openTestDB(t *testing.T) *migrations.Runner {
	t.Helper()
	dir := t.TempDir()
	db, err := sqlite.Open(filepath.Join(dir, "test.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return migrations.New(db)
}

func TestRunner_Up(t *testing.T) {
	r := openTestDB(t)
	ctx := context.Background()
	require.NoError(t, r.Up(ctx))
	require.NoError(t, r.Up(ctx))
}

func TestRunner_Down(t *testing.T) {
	r := openTestDB(t)
	ctx := context.Background()
	require.NoError(t, r.Up(ctx))
	require.NoError(t, r.Down(ctx))
}

func TestRunner_Status(t *testing.T) {
	r := openTestDB(t)
	ctx := context.Background()
	require.NoError(t, r.Status(ctx))
	require.NoError(t, r.Up(ctx))
	require.NoError(t, r.Status(ctx))
}
