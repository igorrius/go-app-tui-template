## 1. Implement Color Log Formatter

- [x] 1.1 Add `ColorTextFormatter` to `internal/logging/formatter.go`.
- [x] 1.2 Implement `Format` method for `ColorTextFormatter` with the following color rules:
    - Timestamps in dimmed gray (e.g., `#555555`).
    - ERROR: Red, WARN: Yellow, INFO: Green, DEBUG: Blue.
    - Module name (capitalized) in bold white.
    - Attribute keys in white.
    - Base message remains uncolored.

## 2. Update TUI Log View

- [x] 2.1 Refactor `internal/tui/screen/logview.go` to use `ColorTextFormatter` instead of `TextFormatter`.
- [x] 2.2 Verify that the `LogView` still renders entries correctly and fits the terminal height.

## 3. Testing and Validation

- [x] 3.1 Update `internal/logging/formatter_test.go` to include test cases for `ColorTextFormatter`.
- [x] 3.2 Add test cases that check for presence of ANSI escape codes in the output.
- [x] 3.3 Run existing tests to ensure no regressions in `TextFormatter`.
