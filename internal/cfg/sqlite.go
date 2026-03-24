package cfg

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/igorrius/go-app-tui-template/internal/infra/sqlite"
	"github.com/samber/do/v2"
)

// ProvideSQLiteDB returns a do provider that opens and returns the application *sql.DB.
func ProvideSQLiteDB() do.Provider[*sql.DB] {
	return func(i do.Injector) (*sql.DB, error) {
		cfg, err := do.Invoke[*config.AppConfig](i)
		if err != nil {
			return nil, err
		}

		dbCfg, err := config.GetDatabase(cfg)
		if err != nil {
			return nil, fmt.Errorf("resolving database config: %w", err)
		}

		if err := os.MkdirAll(filepath.Dir(dbCfg.Path), 0o755); err != nil {
			return nil, fmt.Errorf("creating database directory: %w", err)
		}

		return sqlite.Open(dbCfg.Path)
	}
}
