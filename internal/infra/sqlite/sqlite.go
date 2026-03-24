package sqlite

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	_ "modernc.org/sqlite"
)

// Open opens a SQLite database at the given DSN, applying WAL mode and
// concurrent-access PRAGMAs, then configures the connection pool.
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite database: %w", err)
	}

	if err := applyPragmas(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(runtime.NumCPU())
	db.SetMaxIdleConns(runtime.NumCPU())
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func applyPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA synchronous=NORMAL",
	}
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("applying %q: %w", pragma, err)
		}
	}
	return nil
}
