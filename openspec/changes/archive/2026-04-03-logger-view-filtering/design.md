## Context

The current `LogView` displays all log events in a single stream. As the application runs, the sheer volume of logs makes it difficult to find specific information or monitor only critical errors. This design introduces interactive filtering capabilities directly into the TUI Log View.

## Goals / Non-Goals

**Goals:**
- Provide a dedicated `FilterBar` UI at the bottom of the log screen.
- Enable toggling specific log levels (DEBUG to ERROR) using keys `1` through `4`.
- Enable free-text search across log bodies and attributes using the `/` key.
- Support interactive search input with `Enter` (commit) and `Esc` (clear/close).
- Optimize the filtering process to maintain a responsive TUI.

**Non-Goals:**
- Persisting filter states between application restarts.
- Advanced query syntax or regex support in search.
- Filtering logs stored in files (this only affects the TUI display).

## Decisions

- **State Ownership**: The `LogView` screen will own the `masterEvents` buffer of `logging.LogEvent` objects, the current filter criteria (enabled levels and search query), and the `FilterBar` instance.
- **Filtering Logic**: 
  - Filtering will be performed on the `masterEvents` to produce a display buffer of formatted strings whenever the filters change or a new event arrives.
  - Text search will be case-insensitive and check both `LogEvent.Message` and all `LogEvent.Attrs` (key and value string representations).
- **Log Levels Mapping**:
  1. DEBUG
  2. INFO
  3. WARN
  4. ERROR
- **UI Interaction**:
  - Pressing `/` focuses the search input.
  - While typing, the view updates in real-time (live filtering).
  - `Enter` closes the input but keeps the filter.
  - `Esc` clears the filter if it exists; if input is active but empty, it closes the search input.
- **UI Layout**: The Log View height will be reduced by 1 line to accommodate the `FilterBar` above the global `KeyBar`.

## Risks / Trade-offs

- **Memory/CPU**: Storing and re-filtering thousands of events on every keystroke. We will cap the `masterEvents` to 1000 items initially to ensure smooth performance.
- **Library Compatibility**: Since we use `bubbletea/v2`, we will prefer a lightweight custom implementation for the text input if the standard `bubbles/textinput` (which targets v1) has compatibility issues, or we will wrap it appropriately. In this template, we follow the pattern of manual key handling for the filter bar as seen in `poliarbb`.
