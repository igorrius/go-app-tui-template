# Logging Config

## Purpose

This spec defines the purpose of the Logging Config capability.

## Requirements

### Requirement: Logging configuration section in app config
The system SHALL add a `logging` top-level section to the application configuration supporting the following fields:
- `level` — minimum log level (`debug`, `info`, `warn`, `error`; default `info`)
- `format` — output format (`text` or `json`; default `text`)
- `sink` — destination (`stdout`, `stderr`, or `file`; default `stdout`)
- `dir` — log directory used when `sink: file` (default `logs`)
- `rotation_interval` — duration string for file rotation interval (default `1h`)

#### Scenario: Default configuration is valid
- **WHEN** no `logging` section is present in the YAML config
- **THEN** the system SHALL use default values (`level: info`, `format: text`, `sink: stdout`)

#### Scenario: Override via YAML
- **WHEN** the YAML config contains `logging.sink: file` and `logging.dir: /var/log/app`
- **THEN** the loaded `LoggingConfig` SHALL reflect those values

#### Scenario: Unknown sink value is rejected
- **WHEN** `logging.sink` is set to an unrecognised value
- **THEN** `logging.Init` SHALL return an error describing the invalid value

#### Scenario: Unknown format value is rejected
- **WHEN** `logging.format` is set to an unrecognised value
- **THEN** `logging.Init` SHALL return an error describing the invalid value
