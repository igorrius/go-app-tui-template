## MODIFIED Requirements

### Requirement: Log view screen displays live log entries
The system SHALL provide a `LogView` screen that receives `LogEvent` messages from a provided channel and renders each received log entry as a formatted text line using `logging.ColorTextFormatter`.

#### Scenario: LogView renders received log events
- **WHEN** a `LogEvent` is received via the provided channel
- **THEN** the `LogView` screen SHALL append a formatted text line to its display buffer

#### Scenario: LogView lifecycle and dependency injection
- **WHEN** the `LogView` screen is initialized
- **THEN** it SHALL use a log subscription channel provided during creation, avoiding the creation of new subscriptions on every view change.

#### Scenario: LogView limits displayed lines to terminal height
- **WHEN** the number of buffered log lines exceeds the available view height
- **THEN** the `LogView` SHALL display only the most recent lines that fit, without scroll controls

#### Scenario: LogView resizes with terminal
- **WHEN** the terminal is resized
- **THEN** the `LogView` SHALL adjust its visible line count to match the new height
