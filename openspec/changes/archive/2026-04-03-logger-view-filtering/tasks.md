## 1. FilterBar Component Implementation

- [x] 1.1 Create `internal/tui/component/filterbar.go` to handle log level checkboxes and search input display.
- [x] 1.2 Implement `FilterBar.Update` and `FilterBar.View` methods for rendering the status line and managing input focus.
- [x] 1.3 Add helper methods to `FilterBar` for toggling specific levels and retrieving the current search text.

## 2. LogView Screen Refactoring

- [x] 2.1 Update `LogView` struct in `internal/tui/screen/logview.go` to store `masterEvents []logging.LogEvent` and current filtering state.
- [x] 2.2 Re-implement `LogView.Update` to handle keyboard events for filtering:
  - Numerical keys `1-4` for toggling log levels.
  - `/` key to activate/focus search input in `FilterBar`.
  - `Enter` and `Esc` for search control (commit, clear, close).
- [x] 2.3 Implement the filtering loop in `LogView` to re-generate the visible lines based on `masterEvents`, `enabledLevels`, and `searchQuery`.
- [x] 2.4 Update `LogView.View` to return a composite view containing log lines and the `FilterBar` status line.

## 3. UI/UX Polish

- [x] 3.1 Adjust `LogView.trimBuffer` and height calculations to subtract 1 line for the `FilterBar`.
- [x] 3.2 Add case-insensitive search matching for both log messages and structured attributes.
- [x] 3.3 Ensure visual feedback for active log levels (e.g., color `[x]` vs `[ ]`).
- [x] 3.4 Verify that filters are applied immediately upon keystroke for a smooth interactive experience.
