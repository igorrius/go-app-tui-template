## Context

The current `go.mod` file uses `github.com/charmbracelet/bubbletea v1.3.10`. The user clarified that the latest version is `v2.0.2`.

## Goals / Non-Goals

**Goals:**
- Update `github.com/charmbracelet/bubbletea` to version `v2.0.2`.
- Ensure all TUI components continue to function properly.
- Update `go.mod` and `go.sum` correctly.

**Non-Goals:**
- Redefining the TUI logic or architecture.
- Implementing any new features for the app beyond what's required for the update.

## Decisions

- **Version Selection**: I will use `v1.3.10` if that's truly the latest, or the newest version found via `go get`.
- **Mod Tidy**: I will run `go mod tidy` after the change to clean up `go.sum`.

## Risks / Trade-offs

- [Risk] Breaking change in `bubbletea` → [Mitigation] Review library's release notes and perform a build/test run.
