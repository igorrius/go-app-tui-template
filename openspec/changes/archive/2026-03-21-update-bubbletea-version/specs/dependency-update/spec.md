## ADDED Requirements

### Requirement: Dependency alignment
The project MUST use `github.com/charmbracelet/bubbletea` version `v2.0.2`.

#### Scenario: Verify go.mod version
- **WHEN** checked in `go.mod`
- **THEN** it should be `v2.0.2`

### Requirement: TUI compilation
The TUI code MUST compile successfully with the updated version of `bubbletea`.

#### Scenario: Build project
- **WHEN** running `go build ./...`
- **THEN** the project should build without errors
