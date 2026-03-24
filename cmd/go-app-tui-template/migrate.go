package main

import (
	"context"
	"fmt"
	"os"

	"github.com/igorrius/go-app-tui-template/internal/cfg"
	"github.com/igorrius/go-app-tui-template/internal/infra/migrations"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

var migrateCommand = &cli.Command{
	Name:  "migrate",
	Usage: "Manage database schema migrations",
	Commands: []*cli.Command{
		{
			Name:  "up",
			Usage: "Apply all pending migrations",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				runner, err := newMigrationRunner(cmd)
				if err != nil {
					return err
				}
				if err := runner.Up(ctx); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "migrate up: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("migrations applied successfully")
				return nil
			},
		},
		{
			Name:  "down",
			Usage: "Roll back the most recently applied migration",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				runner, err := newMigrationRunner(cmd)
				if err != nil {
					return err
				}
				if err := runner.Down(ctx); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "migrate down: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("migration rolled back successfully")
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "Print the current migration status",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				runner, err := newMigrationRunner(cmd)
				if err != nil {
					return err
				}
				if err := runner.Status(ctx); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "migrate status: %v\n", err)
					os.Exit(1)
				}
				return nil
			},
		},
	},
}

func newMigrationRunner(cmd *cli.Command) (*migrations.Runner, error) {
	injector := do.New()
	do.Provide(injector, cfg.ProvideConfig(cmd.Root()))
	do.Provide(injector, cfg.ProvideSQLiteDB())
	do.Provide(injector, cfg.ProvideMigrationRunner())

	runner, err := do.Invoke[*migrations.Runner](injector)
	if err != nil {
		return nil, fmt.Errorf("initializing migration runner: %w", err)
	}
	return runner, nil
}
