# Go App TUI Template

A foundational template repository for building production-ready Go Terminal User Interface (TUI) applications. This template is designed with a layered architecture, dependency injection, and essential infrastructure to help you build maintainable, concurrent, and robust TUI apps quickly.

## Features

- **TUI Framework**: Built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) (`v2`) and [Lipgloss](https://github.com/charmbracelet/lipgloss) for a modern, responsive, full-terminal experience. Includes a multi-screen router and an HTOP-style function-key bar.
- **Dependency Injection**: Uses [`samber/do`](https://github.com/samber/do) for a clean, type-safe DI container, isolating components into `cfg` providers.
- **Layered Architecture**: Pre-scaffolded `domain`, `infra`, and `service` packages to enforce clean boundaries and maintainable business logic.
- **Configuration Management**: 
  - Loads from `app-config.yaml` (local override) or `app-config.dist.yaml` (default).
  - Supports environment variable interpolation (e.g., `${LOG_LEVEL:-info}`).
  - Automatically loads `.env` files using `godotenv`.
  - Provides generic, dot-path access (e.g., `config.MustGet[string]("app.name")`).
- **Structured Logging**: 
  - Global `slog.Logger` powered by a decoupled message bus ([`go-vein`](https://github.com/igorrius/go-vein)).
  - Supports asynchronous writes to `stdout`/`stderr` or time-rotated log files.
  - Built-in `LogView` TUI screen (press **F9**) with colorized, real-time log streaming.
- **Embedded Database & Migrations**: 
  - Uses [`modernc.org/sqlite`](https://gitlab.com/cznic/sqlite) for a pure-Go, CGO-free SQLite driver.
  - Configured out-of-the-box for high-concurrency with WAL mode.
  - Includes a built-in migration runner using [`pressly/goose`](https://github.com/pressly/goose) with embedded SQL files.
- **CLI Management**: Powered by [`urfave/cli`](https://github.com/urfave/cli) (`v3`) for elegant subcommands, flags, and environment bindings.

## Project Structure

```text
.
├── cmd/
│   └── app/                   # Application entrypoint and CLI subcommands
├── internal/
│   ├── cfg/                   # Dependency Injection providers
│   ├── config/                # YAML and Env configuration loading
│   ├── domain/                # Core business models and interfaces
│   ├── infra/                 # Database, migrations, and external adapters
│   ├── logging/               # Structured logging handlers and formatters
│   ├── service/               # Application logic orchestrating domain and infra
│   └── tui/                   # Bubble Tea components, screens, and routing
├── openspec/                  # OpenSpec artifacts and specs
├── app-config.dist.yaml       # Default configuration schema
└── go.mod                     # Go module dependencies
```

## Getting Started

### Prerequisites
- Go 1.26 or later

### Installation

1. Clone this repository (or use it as a template on GitHub):
   ```bash
   git clone https://github.com/igorrius/go-app-tui-template.git my-app
   cd my-app
   ```

2. Tidy dependencies and build:
   ```bash
   go mod tidy
   go build ./cmd/app
   ```

### Running the App

Run the built binary to start the TUI:
```bash
./app
```

**TUI Navigation:**
- **F1**: Dashboard View
- **F9**: Real-time Log View
- **F10** or **Ctrl+C**: Quit

### Configuration

Configuration is managed via YAML. By default, the app loads `app-config.dist.yaml`. 
To override settings locally, create an `app-config.yaml` (which is git-ignored) or use a `.env` file.

Example `.env`:
```env
LOG_LEVEL=debug
APP_CONFIG=./app-config.yaml
```

### Database Migrations

The binary includes a built-in migration runner for SQLite. Migrations are written as plain SQL files in `internal/infra/migrations/sql/` and embedded directly into the binary.

```bash
# Apply pending migrations
./app migrate up

# Rollback the last migration
./app migrate down

# Check migration status
./app migrate status
```

## Development Workflow

This template was designed with an OpenSpec-driven development workflow in mind, capturing design decisions, capabilities, and tasks in the `openspec/` directory. Check out the `.github/skills` directory if you use an AI assistant for managing specs.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.