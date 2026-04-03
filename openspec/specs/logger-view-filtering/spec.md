# Capability: Interactive Logger Filtering

## Purpose
Enables real-time interactive log filtering and searching in the TUI, improving system observability.

## Requirements

### Requirement: Interactive Log Level Toggling
The system SHALL provide interactive toggling of log levels DEBUG, INFO, WARN, ERROR via numerical keys 1 through 4 while the `LogView` screen is active.

#### Scenario: Key 1-4 toggles log levels
- **GIVEN** `LogView` is the active screen
- **WHEN** the user presses key `1`
- **THEN** the DEBUG level filter SHALL be toggled (enabled/disabled)
- **AND** the `LogView` SHALL immediately update its visible entries

### Requirement: Interactive Full-Text Search
The system SHALL provide an interactive search input triggered by the `/` key while the `LogView` screen is active.

#### Scenario: `/` key activates search input
- **GIVEN** `LogView` is the active screen
- **WHEN** the user presses key `/`
- **THEN** the `FilterBar` SHALL focus the search input field
- **AND** the input text SHALL be blank initially if no filter was previously set

#### Scenario: Real-time search filtering
- **GIVEN** the search input is focused
- **WHEN** the user types characters
- **THEN** the `LogView` SHALL update its display in real-time, filtering for the current search text

#### Scenario: Search navigation keys
- **GIVEN** the search input is focused
- **WHEN** the user presses key `Enter`
- **THEN** the `FilterBar` SHALL unfocus the search input field but keep the filter active
- **WHEN** the user presses key `Esc`
- **THEN** if a filter was present, the filter SHALL be cleared and the input field SHALL be emptied
- **AND** if the input field was already empty and focused, it SHALL be unfocused (closed)

### Requirement: FilterBar provides visual feedback
The system SHALL display the state of all filters in a single line `FilterBar`.

#### Scenario: FilterBar displays level checkboxes
- **WHEN** rendered
- **THEN** the `FilterBar` SHALL show checkboxes for each log level (e.g., `1[x] INFO`) indicating which levels are currently being displayed

#### Scenario: FilterBar displays current search string
- **WHEN** a search filter is active (even if unfocused)
- **THEN** the `FilterBar` SHALL show the label `Filter: ` followed by the current search text
