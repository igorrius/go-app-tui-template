## Why

Updating `github.com/charmbracelet/bubbletea` to the latest version (v2.0.2) ensures we have the latest features, bug fixes, and security updates from the library.

## What Changes

- Update `github.com/charmbracelet/bubbletea` in `go.mod` to version `v2.0.2`.
- Update `go.sum` via `go mod tidy`.
- Verify that TUI components (Bubble Tea programs) still function correctly with the latest version. **BREAKING**: v2.0.2 may contain breaking changes compared to v1.x.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- None. (This is a dependency update that does not change the core requirements of our capabilities).

## Impact

- `go.mod` and `go.sum` will be modified.
- All files importing `github.com/charmbracelet/bubbletea` (e.g., [internal/tui/app.go](internal/tui/app.go)) may be affected by any breaking changes in the library, although typically minor version updates are backwards compatible.
