# Logging Vein Handler

## Purpose

This spec defines the purpose of the Logging Vein Handler capability.

## Requirements

### Requirement: VeinHandler implements slog.Handler
The system SHALL provide a `VeinHandler` struct that implements `slog.Handler` and publishes a `LogEvent` to the application message bus on each `Handle` call.

#### Scenario: Handle publishes event to bus
- **WHEN** `VeinHandler.Handle` is called with a `slog.Record`
- **THEN** the handler SHALL publish exactly one `LogEvent` containing the log level, message, timestamp, and all attributes to the configured `go-vein` bus topic

#### Scenario: Enabled respects configured level
- **WHEN** `VeinHandler.Enabled` is called with a level below the configured minimum
- **THEN** the method SHALL return `false` and no event SHALL be published

#### Scenario: Enabled passes records at or above minimum level
- **WHEN** `VeinHandler.Enabled` is called with a level at or above the configured minimum
- **THEN** the method SHALL return `true`

#### Scenario: WithAttrs returns new handler with attributes
- **WHEN** `VeinHandler.WithAttrs` is called with a set of attributes
- **THEN** it SHALL return a new `VeinHandler` that includes those attributes in every subsequent `LogEvent`

#### Scenario: WithGroup returns new handler with group prefix
- **WHEN** `VeinHandler.WithGroup` is called with a group name
- **THEN** it SHALL return a new `VeinHandler` that namespaces subsequent attributes under that group name

### Requirement: LogEvent carries all slog record data
The system SHALL define a `LogEvent` struct used as the message bus payload, containing: `Time time.Time`, `Level slog.Level`, `Message string`, `Attrs []slog.Attr`.

#### Scenario: LogEvent round-trips slog record fields
- **WHEN** a `slog.Record` with time, level, message, and two attrs is handled
- **THEN** the published `LogEvent` SHALL contain every field from the record with no loss
