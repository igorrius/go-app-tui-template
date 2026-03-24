## ADDED Requirements

### Requirement: Layered package structure
The system SHALL organize code into `internal/domain/`, `internal/infra/`, and `internal/service/` packages representing distinct architectural layers.

#### Scenario: Domain layer isolation
- **WHEN** domain types are defined
- **THEN** they SHALL reside in `internal/domain/` and SHALL NOT import from `infra` or `service` packages

#### Scenario: Service layer depends on domain
- **WHEN** service logic requires domain types
- **THEN** `internal/service/` SHALL import from `internal/domain/` but NOT from `internal/infra/`

#### Scenario: Infra layer implements domain interfaces
- **WHEN** infrastructure adapters are needed
- **THEN** `internal/infra/` SHALL implement interfaces defined in `internal/domain/`

### Requirement: Entry point in cmd directory
The application binary entry point SHALL reside at `cmd/go-app-tui-template/main.go`.

#### Scenario: Build target
- **WHEN** `go build ./cmd/go-app-tui-template` is run
- **THEN** a working binary SHALL be produced

### Requirement: Scaffold packages with documentation
Each architectural layer package SHALL include a `doc.go` file describing its purpose.

#### Scenario: Domain doc exists
- **WHEN** a developer opens `internal/domain/`
- **THEN** `doc.go` SHALL describe the domain layer's purpose and constraints
