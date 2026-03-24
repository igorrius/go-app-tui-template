## MODIFIED Requirements

### Requirement: Bottom function-key bar
The system SHALL display a bottom status bar showing available function key actions, styled similar to HTOP.

#### Scenario: Key bar shows available actions
- **WHEN** the TUI is active
- **THEN** the bottom row SHALL display `F1 Dashboard`, `F9 Logs`, and `F10 Quit` labels

#### Scenario: F10 quits the application
- **WHEN** the user presses F10
- **THEN** the application SHALL exit cleanly

#### Scenario: F1 switches to Dashboard
- **WHEN** the user presses F1
- **THEN** the active screen SHALL switch to the Dashboard view

#### Scenario: F9 switches to LogView
- **WHEN** the user presses F9
- **THEN** the active screen SHALL switch to the LogView screen
