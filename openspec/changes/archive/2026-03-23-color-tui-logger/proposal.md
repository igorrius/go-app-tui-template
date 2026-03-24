## Why

The current `TextFormatter` in the TUI log view lacks visual hierarchy, making it difficult to distinguish timestamps, log levels, and modules at a glance. Adding color-coded formatting will improve readability and help users quickly identify important information in the logs.

## What Changes

- Replace `TextFormatter` with `ColorTextFormatter` in the TUI application.
- Implement color-coded rendering for log entries:
    - Dimmed timestamps to reduce visual noise.
    - Level-specific colors (e.g., Red for ERROR, Yellow for WARN, Green for INFO).
    - Bold and white module names for clear separation.
    - White-colored attribute keys for better visibility.
- Ensure the base log message remains uncolored for maximum clarity.

## Capabilities

### New Capabilities
- `tui-color-logger`: Implements a specialized color-aware log formatter for the Bubble Tea based TUI.

### Modified Capabilities
- `tui-log-view`: Update the log view to use the new color-aware formatter.

## Impact

- `internal/logging/formatter.go`: Add `ColorTextFormatter` and its implementation.
- `internal/tui/screen/logview.go`: Switch to using `ColorTextFormatter`.
- `internal/logging/formatter_test.go`: Add tests for the new formatter.
