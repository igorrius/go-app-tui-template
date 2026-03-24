## ADDED Requirements

### Requirement: Color-coded log formatting for TUI
The system SHALL provide a `ColorTextFormatter` that applies ANSI color codes to `LogEvent` fields for enhanced readability in the TUI log view.

#### Scenario: Timestamp is dimmed
- **WHEN** formatting a `LogEvent`
- **THEN** the timestamp SHALL be wrapped in a faint or dark color ANSI code (e.g., `#555555`)

#### Scenario: Log level is color-coded
- **WHEN** formatting a `LogEvent` with `slog.LevelError`
- **THEN** the level (e.g., "ERROR") SHALL be colored red

#### Scenario: Module name is bold white
- **WHEN** formatting a `LogEvent`
- **THEN** the module name (e.g., "APP") SHALL be colored bold white

#### Scenario: Attribute keys are white
- **WHEN** formatting a `LogEvent` with attributes
- **THEN** each attribute key SHALL be colored white, while the value remains uncolored

#### Scenario: Base message remains uncolored
- **WHEN** formatting a `LogEvent`
- **THEN** the message text (excluding level, module, etc.) SHALL not have any additional ANSI color applied to it
