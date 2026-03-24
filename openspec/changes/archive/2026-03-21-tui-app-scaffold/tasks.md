## 1. Project Structure & Dependencies

- [x] 1.1 Create directory structure: `cmd/go-app-tui-template/`, `internal/domain/`, `internal/infra/`, `internal/service/`, `internal/cfg/`, `internal/config/`, `internal/tui/screen/`, `internal/cmd/flag/`
- [x] 1.2 Add `doc.go` files for `domain`, `infra`, and `service` packages describing each layer's purpose
- [x] 1.3 Add Go dependencies: `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`, `github.com/samber/do/v2`, `github.com/urfave/cli/v3`, `gopkg.in/yaml.v3`

## 2. Configuration Management

- [x] 2.1 Create `app-config.dist.yaml` with sample structure (app name, version, placeholder sections)
- [x] 2.2 Add `app-config.yaml` to `.gitignore`
- [x] 2.3 Implement `internal/config/config.go`: YAML file loader with fallback logic (`app-config.yaml` → `app-config.dist.yaml`), `AppConfig` struct holding parsed `map[string]any`
- [x] 2.4 Implement generic `Get[T](cfg *AppConfig, path string) (T, error)` and `MustGet[T](cfg *AppConfig, path string) T` functions with dot-path traversal

## 3. Dependency Injection

- [x] 3.1 Implement `internal/cfg/providers.go`: `ProvideConfig` provider using `samber/do` that loads and returns `*config.AppConfig`
- [x] 3.2 Add `ProvideTUIApp` provider that resolves config dependency and returns the root TUI model

## 4. CLI Entry Point

- [x] 4.1 Implement `internal/cmd/flag/root.go`: root `urfave/cli/v3` command with `--config` flag and env variable support
- [x] 4.2 Default action wires DI container, resolves TUI app, and starts bubbletea program in full-terminal (alt screen) mode
- [x] 4.3 Implement `cmd/go-app-tui-template/main.go`: minimal entry point that calls the root CLI command

## 5. TUI Core

- [x] 5.1 Implement `internal/tui/screen/dashboard.go`: Dashboard screen implementing `Screen` interface, renders "In construction" centered
- [x] 5.2 Implement `internal/tui/keybar.go`: bottom function-key bar component rendering `F1 Dashboard` and `F10 Quit` labels using lipgloss
- [x] 5.3 Implement `internal/tui/app.go`: root bubbletea model — manages active screen, terminal size tracking, key bar rendering, global F-key handling (F1 → Dashboard, F10 → quit)
- [x] 5.4 Define `Screen` interface in `internal/tui/screen.go` with `Update(tea.Msg) (Screen, tea.Cmd)` and `View() string` methods

## 6. Verification

- [x] 6.1 Run `go build ./cmd/go-app-tui-template` and verify the binary compiles
- [x] 6.2 Run the binary and verify: full-terminal layout, bottom key bar visible, Dashboard shows "In construction", F10 exits
