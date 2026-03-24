## Context

The go-app-tui-template project is a greenfield Go TUI application. There is no existing code beyond a `go.mod` file. The project needs a well-structured scaffold that establishes architectural patterns for all future development.

Technology stack:
- **Go 1.26.1** (module: `github.com/igorrius/go-app-tui-template`)
- **bubbletea** (`github.com/charmbracelet/bubbletea`) — TUI framework
- **samber/do** (`github.com/samber/do/v2`) — dependency injection
- **urfave/cli** (`github.com/urfave/cli/v3`) — CLI command/flag management
- **gopkg.in/yaml.v3** — YAML config parsing

## Goals / Non-Goals

**Goals:**
- Establish a domain / infra / service layered architecture
- Create a working TUI with full-terminal layout and HTOP-style function-key bar
- Implement generic typed config access (`Config.Get[T]("path")` / `Config.MustGet[T]("path")`)
- Set up DI container with `samber/do` using `Provider[T]` pattern in a `cfg` package
- Wire CLI entry point with `urfave/cli` for command/flag/env management in a `cmd/flag` package
- Provide `app-config.dist.yaml` (tracked) and `app-config.yaml` (gitignored) for business logic configuration

**Non-Goals:**
- Implementing any business domain logic beyond the scaffold
- Adding network/API layers
- Building multiple screens beyond Dashboard placeholder
- Persistent data storage or database integration

## Decisions

### 1. Package Layout — Domain / Infra / Service Layers

```
cmd/
  go-app-tui-template/
    main.go            # Entry point
internal/
  cfg/                 # DI container providers (samber/do)
    providers.go
  cmd/
    flag/              # urfave/cli commands, flags, env bindings
      root.go
  config/              # YAML config loader with generic access
    config.go
  domain/              # Domain types and interfaces (empty scaffold)
  infra/               # Infrastructure adapters (empty scaffold)
  service/             # Application services (empty scaffold)
  tui/
    app.go             # Root bubbletea model (main window)
    keybar.go          # Bottom function-key bar component
    screen/
      dashboard.go     # Dashboard screen ("In construction")
```

**Rationale:** Follows standard Go project layout (`cmd/` + `internal/`). The domain/infra/service split provides clear boundaries without over-abstracting at scaffold stage. TUI components are isolated in `internal/tui/`.

### 2. Dependency Injection — samber/do with cfg Package

All components are registered as providers in `internal/cfg/providers.go`:

```go
func ProvideConfig(cli.Context) do.Provider[*config.AppConfig] { ... }
func ProvideTUIApp() do.Provider[do.Provider] { ... }
```

The DI container is initialized at startup and passed through the application. Each provider function follows the `func ProvideXxx() do.Provider[T]` pattern.

**Rationale:** `samber/do` provides a lightweight DI container. Centralizing providers in `cfg` makes the dependency graph visible and easy to extend.

### 3. Configuration — Generic Typed Access

Config is loaded from YAML files into a nested `map[string]any` structure. Access is provided through generic functions:

```go
func Get[T any](cfg *AppConfig, path string) (T, error)
func MustGet[T any](cfg *AppConfig, path string) T
```

Path uses dot notation: `"server.port"`, `"app.name"`. `MustGet` panics on missing/invalid keys.

Config loading priority: `app-config.yaml` (local, gitignored) falls back to `app-config.dist.yaml` (tracked).

**Rationale:** Generics eliminate type-switch boilerplate. Dot-path access is natural for nested YAML. Two-file pattern (dist + local override) is a well-established convention.

### 4. CLI Management — urfave/cli in cmd/flag Package

`internal/cmd/flag/root.go` defines the root CLI command with `urfave/cli/v3`. Subcommands and flags are registered here. The default action launches the TUI.

**Rationale:** Separating CLI definition from `main.go` keeps the entry point minimal. `urfave/cli/v3` provides a mature, idiomatic CLI framework with built-in env variable binding.

### 5. TUI Architecture — bubbletea with Screen Router

The root model (`tui/app.go`) manages:
- Current active screen (enum-based routing)
- Terminal size (tracked via `tea.WindowSizeMsg`)
- Bottom key bar rendering

The key bar (`tui/keybar.go`) renders function key labels: `F1 Dashboard`, `F10 Quit`. Global key handling is done at the root model level.

Each screen implements a common interface:

```go
type Screen interface {
    Update(msg tea.Msg) (Screen, tea.Cmd)
    View() string
}
```

**Rationale:** Screen interface allows adding new screens without modifying the root model. The root model handles global concerns (resize, key bar, function keys) while delegating content to the active screen.

### 6. Color Palette — bubbletea/lipgloss Default

Use `lipgloss` default styling (from bubbletea ecosystem). No custom theme at scaffold stage.

**Rationale:** Bubbletea's default palette works across terminal color schemes. Custom theming can be layered on later.

## Risks / Trade-offs

- **[Generic config access is runtime-typed]** → Config path errors surface at runtime, not compile time. Mitigation: `MustGet` fails fast during startup for required values.
- **[Empty domain/infra/service packages]** → May confuse contributors initially. Mitigation: Each package gets a minimal doc.go with purpose description.
- **[Single config file format (YAML only)]** → Locks to YAML. Mitigation: Config package abstracts the format behind its API; switching parsers later is internal.
