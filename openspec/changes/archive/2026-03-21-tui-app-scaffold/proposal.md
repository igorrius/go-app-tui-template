## Why

The go-app-tui-template project needs its foundational TUI application scaffold. Without a structured base, building terminal UI features becomes ad-hoc and inconsistent. Establishing the core architecture now — with proper dependency injection, configuration management, and a layered domain/infra/service design — sets the standard for all future development.

## What Changes

- Create the main application entry point with `urfave/cli` for command/flag management
- Set up `bubbletea`-based TUI with a full-terminal window layout and bottom function-key bar (HTOP-style: F1 Dashboard, F10 Quit)
- Implement dependency injection using `samber/do` with a `cfg` package exposing `Provider[T]` functions
- Add YAML-based application configuration via `app-config.yaml` / `app-config.dist.yaml` with generic typed access (`Config.Get[T]("path")`, `Config.MustGet[T]("path")`)
- Organize code using domain / infra / service layers
- Default view: Dashboard screen displaying "In construction"

## Capabilities

### New Capabilities
- `app-entrypoint`: CLI entry point with urfave/cli, command/flag registration, and application bootstrap
- `tui-layout`: Full-terminal bubbletea window with bottom function-key bar and screen routing (F1 Dashboard, F10 Quit)
- `config-management`: YAML config loader with generic typed access (Get[T]/MustGet[T]) via a dedicated package
- `dependency-injection`: samber/do-based DI container with Provider[T] pattern in a cfg package
- `project-structure`: Domain/infra/service layered package organization

### Modified Capabilities

## Impact

- New Go dependencies: `github.com/charmbracelet/bubbletea`, `github.com/samber/do`, `github.com/urfave/cli/v3`, a YAML parser
- New directories: `cmd/`, `internal/domain/`, `internal/infra/`, `internal/service/`, `internal/tui/`, `internal/cfg/`, `internal/config/`
- New config files: `app-config.yaml` (gitignored), `app-config.dist.yaml`
- `go.mod` updated with new dependencies
