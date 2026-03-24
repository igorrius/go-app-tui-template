## 1. Preparation

- [x] 1.1 Update `go.mod` to use `charm.land/bubbletea/v2@v2.0.2`
- [x] 1.2 Run `go mod tidy` to update `go.sum`

## 2. Verification

- [x] 2.1 Run `go build ./...` to check for compilation errors
- [x] 2.2 Fix any breaking changes introduced by `v2.0.2` (e.g., package paths, renamed methods)
- [x] 2.3 Run the TUI locally to ensure dashboard and terminal interaction works as expected
