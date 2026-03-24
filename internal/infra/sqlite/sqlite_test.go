package sqlite_test

import (
	"path/filepath"
	"testing"

	"github.com/igorrius/go-app-tui-template/internal/infra/sqlite"
	"github.com/stretchr/testify/require"
)

func TestOpen_PRAGMAsAreSet(t *testing.T) {
	dir := t.TempDir()
	dsn := filepath.Join(dir, "test.db")

	db, err := sqlite.Open(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	tests := []struct {
		pragma   string
		query    string
		expected string
	}{
		{
			pragma:   "journal_mode",
			query:    "PRAGMA journal_mode",
			expected: "wal",
		},
		{
			pragma:   "busy_timeout",
			query:    "PRAGMA busy_timeout",
			expected: "5000",
		},
		{
			pragma:   "synchronous",
			query:    "PRAGMA synchronous",
			expected: "1", // NORMAL = 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.pragma, func(t *testing.T) {
			var val string
			row := db.QueryRow(tt.query)
			require.NoError(t, row.Scan(&val))
			require.Equal(t, tt.expected, val)
		})
	}
}
