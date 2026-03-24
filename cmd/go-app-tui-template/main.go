package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/igorrius/go-app-tui-template/cmd/flag"
	"github.com/igorrius/go-app-tui-template/internal/cfg"
	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	"github.com/igorrius/go-app-tui-template/internal/tui"
	vein "github.com/igorrius/go-vein"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

func main() {
	root := &cli.Command{
		Name:    "go-app-tui-template",
		Usage:   "A terminal UI application",
		Version: "0.1.0",
		Flags: []cli.Flag{
			flag.ConfigFilePath,
		},
		Commands: []*cli.Command{
			migrateCommand,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			injector := do.New()
			do.Provide(injector, cfg.ProvideConfig(cmd))
			do.Provide(injector, cfg.ProvideDispatcher())
			do.Provide(injector, cfg.ProvideTUIApp())

			appCfg, err := do.Invoke[*config.AppConfig](injector)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			loggingCfg, err := config.GetLogging(appCfg)
			if err != nil {
				return fmt.Errorf("loading logging config: %w", err)
			}

			bus, err := do.Invoke[*vein.Dispatcher](injector)
			if err != nil {
				return fmt.Errorf("initializing bus: %w", err)
			}

			if err := logging.Init(loggingCfg, bus); err != nil {
				return fmt.Errorf("initializing logging: %w", err)
			}
			defer func() {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = logging.Shutdown(shutdownCtx)
			}()

			app, err := do.Invoke[*tui.App](injector)
			if err != nil {
				return fmt.Errorf("initializing application: %w", err)
			}

			p := tea.NewProgram(app)
			if _, err := p.Run(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	if err := root.Run(context.Background(), os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
