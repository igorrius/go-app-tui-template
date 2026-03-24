## ADDED Requirements

### Requirement: Text formatter renders log records in canonical field order
The system SHALL provide a `logging.TextFormatter` type with a `Format(LogEvent) string` method that renders each log record in the following fixed order: `Time`, `Level`, `Module`, then all remaining attributes as `key=value` pairs.

#### Scenario: Full record with module attribute
- **WHEN** `Format` is called with a `LogEvent` whose `Attrs` contain a `"module"` key with value `"db"`
- **THEN** the returned string SHALL begin with the time, followed by the uppercased level, then `DB`, then any remaining attributes

#### Scenario: Record without module attribute defaults to APP
- **WHEN** `Format` is called with a `LogEvent` whose `Attrs` do not contain a `"module"` key
- **THEN** the returned string SHALL include `APP` as the module field

### Requirement: Module field is always rendered in uppercase
The system SHALL render the module field in uppercase regardless of the casing stored in the `"module"` attribute value.

#### Scenario: Lowercase module value is uppercased
- **WHEN** `Format` is called with `module=log_view`
- **THEN** the rendered module field SHALL be `LOG_VIEW`

#### Scenario: Mixed-case module value is uppercased
- **WHEN** `Format` is called with `module=LogView`
- **THEN** the rendered module field SHALL be `LOGVIEW`

### Requirement: Module attribute is excluded from remaining attributes output
The system SHALL NOT repeat the `"module"` attribute in the trailing key=value pairs section of the formatted line.

#### Scenario: Module not duplicated in trailing attrs
- **WHEN** `Format` is called with attrs `[module=db, query=select]`
- **THEN** the output SHALL contain `DB` as the module field and `query=select` as a trailing attr, but SHALL NOT contain `module=db` again

### Requirement: Time field uses a fixed human-readable format
The system SHALL format the `Time` field of each log record using `time.DateTime` layout (`2006-01-02 15:04:05`), wrapped in square brackets.

#### Scenario: Time field format
- **WHEN** `Format` is called with a `LogEvent` whose `Time` is `2026-03-23 14:05:00 UTC`
- **THEN** the output SHALL start with `[2026-03-23 14:05:00]`
