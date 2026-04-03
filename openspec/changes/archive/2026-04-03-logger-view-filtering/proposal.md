## Why

The current TUI log view displays all logs chronologically but lacks the ability to filter by log level or search for specific content. As the system grows and more modules generate logs, finding specific information or troubleshooting errors becomes difficult due to the volume of logs.

## What Changes

- Add a filtering bar at the bottom of the log screen, above the functional panel (F-keys).
- Implementation of log level checkboxes (DEBUG, INFO, WARN, ERROR) that can be toggled using keys 1-4.
- Implementation of a search filter triggered by the `/` key.
- Search should filter log rows where either the body or any attribute (name or value) matches the filter text.
- Interactive controls for the search input: `Enter` to finish, `Esc` to clear or close.
- Visual display of the current filter string: `Filter: <string>`.

## Capabilities

### New Capabilities
- `logger-view-filtering`: Interactive log level toggling and full-text search in the TUI log view.

### Modified Capabilities
- `tui-log-view`: Add interactive filtering and search requirements.

## Impact

- `internal/tui/`: Log screen component and related models will be updated.
- Keyboard event handling in the TUI will be extended to handle filtering keys.
