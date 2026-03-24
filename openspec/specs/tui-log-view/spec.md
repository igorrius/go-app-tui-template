## ADDED Requirements

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

### Requirement: LogView emits subscription info log on activation
The system SHALL emit a structured info-level log entry when `LogView.Init()` is called, before the subscription command is returned.

#### Scenario: Subscription log on Init
- **WHEN** the `LogView` screen is initialized (Init called)
- **THEN** the system SHALL emit `slog.Info("Subscription to log was performed", "module", "LOG_VIEW")` exactly once

### Requirement: LogView activated via F9 key
The system SHALL activate the `LogView` screen when the user presses F9.

#### Scenario: F9 navigates to LogView
- **WHEN** the user presses F9 while any other screen is active
- **THEN** the root `App` SHALL switch the active screen to `LogView`

#### Scenario: F9 pressed while LogView is already active
- **WHEN** the user presses F9 while `LogView` is already the active screen
- **THEN** the root `App` SHALL remain on `LogView` (no-op navigation)

### Requirement: LogView has no scroll capability
The system SHALL NOT provide any scrolling mechanism in the `LogView` screen. The view is append-only and displays the tail of the log buffer.

#### Scenario: Scroll key presses are ignored
- **WHEN** the user presses any scroll or navigation key (arrow keys, page-up, page-down) while `LogView` is active
- **THEN** the displayed log content SHALL NOT change due to the key press
