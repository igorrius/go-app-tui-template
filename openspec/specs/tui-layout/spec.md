# TUI Layout

## Purpose

This spec defines the purpose of the TUI Layout capability.

## Requirements

### Requirement: Full-terminal window layout
The system SHALL render the TUI occupying the entire available terminal area, dynamically adjusting to terminal resize events.

#### Scenario: Initial render fills terminal
- **WHEN** the TUI starts
- **THEN** the root view SHALL occupy the full width and height of the terminal

#### Scenario: Terminal resize adjusts layout
- **WHEN** the terminal is resized
- **THEN** the TUI layout SHALL re-render to fill the new dimensions

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

### Requirement: Screen routing
The system SHALL support multiple screens routed by an enum-based identifier, with a common `Screen` interface for update/view lifecycle.

#### Scenario: Default screen is Dashboard
- **WHEN** the TUI starts
- **THEN** the Dashboard screen SHALL be the active screen

#### Scenario: Screen switch updates view
- **WHEN** a function key for a different screen is pressed
- **THEN** the root model SHALL replace the active screen and re-render

### Requirement: Dashboard screen placeholder
The Dashboard screen SHALL display the text "In construction" centered in the available area.

#### Scenario: Dashboard renders placeholder
- **WHEN** the Dashboard screen is active
- **THEN** the view SHALL display "In construction" centered horizontally and vertically
