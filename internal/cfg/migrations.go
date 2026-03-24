package cfg

import (
	"database/sql"

	"github.com/igorrius/go-app-tui-template/internal/infra/migrations"
	"github.com/samber/do/v2"
)

// ProvideMigrationRunner returns a do provider that creates and returns a *migrations.Runner.
func ProvideMigrationRunner() do.Provider[*migrations.Runner] {
	return func(i do.Injector) (*migrations.Runner, error) {
		db, err := do.Invoke[*sql.DB](i)
		if err != nil {
			return nil, err
		}
		return migrations.New(db), nil
	}
}
