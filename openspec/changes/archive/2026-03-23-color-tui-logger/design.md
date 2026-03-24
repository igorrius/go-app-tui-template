## Context

The current TUI log view uses a plain `TextFormatter` that renders all log fields in the same color. This makes it difficult for users to quickly scan the logs for errors, module sources, or timestamps. The system already uses `lipgloss` or standard ANSI escape codes for TUI styling, so we should leverage these capabilities to enhance the log view.

## Goals / Non-Goals

**Goals:**
- Implement `ColorTextFormatter` in `internal/logging/formatter.go`.
- Use ANSI color codes for specific log fields:
    - Dimmed gray for timestamps.
    - Level-specific colors (Red: Error, Yellow: Warn, Green: Info, Blue: Debug).
    - Bold white for module names.
    - White for attribute keys.
- Replace `TextFormatter` with `ColorTextFormatter` in the `LogView` screen.
- Ensure the base log message is uncolored.

**Non-Goals:**
- Add complex styling like backgrounds or borders.
- Implement a generic styling engine; the colors are fixed for the TUI use case.
- Change the log entry data structure.

## Decisions

- **Decision: Hardcoded ANSI Color Codes vs. Library Styles.**
    - **Rationale:** While `lipgloss` is used in the project for higher-level TUI layouts, the log formatter should ideally produce raw strings with ANSI codes that can be efficiently appended to a buffer. However, `lipgloss` provides convenient and readable style definitions.
    - **Outcome:** We will use `lipgloss` for defining styles (or standard ANSI sequences if `lipgloss` isn't preferred for simple log formatting) but ensure the output is a standard string. If `lipgloss` is used, it should be the same version as the rest of the project. Given the prompt rules, we'll manually apply ANSI codes or use `lipgloss` if it's already a dependency for the TUI.
    - **Alternatives:** Manual ANSI escape code concatenation (more performant, but less maintainable).

- **Decision: Update `logview.go` to use the new formatter.**
    - **Rationale:** The `LogView` screen is specifically designed for TUI display, where color is beneficial. Other sinks (like file logging) should continue using the plain `TextFormatter`.
    - **Outcome:** Update `internal/tui/screen/logview.go` to use `ColorTextFormatter` explicitly.

## Risks / Trade-offs

- **Risk: Terminal compatibility.**
    - **Mitigation:** Use standard ANSI-16 or ANSI-256 colors that are universally supported by modern TUI terminals.
- **Risk: Performance of formatting.**
    - **Mitigation:** The `ColorTextFormatter` will be used only for the TUI display buffer, which is sized to the terminal height. The overhead of adding ANSI codes is minimal.
