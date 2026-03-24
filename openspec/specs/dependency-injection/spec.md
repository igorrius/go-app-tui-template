## ADDED Requirements

### Requirement: DI container with Provider pattern
The system SHALL use `samber/do` as the dependency injection container, with all providers defined as `func ProvideXxx(i do.Injector) (T, error)` in the `internal/cfg` package.

#### Scenario: Register and resolve config provider
- **WHEN** the DI container is initialized
- **THEN** `ProvideConfig` SHALL be registered and resolvable, returning a loaded `*AppConfig`

#### Scenario: Register and resolve TUI app provider
- **WHEN** the DI container is initialized
- **THEN** `ProvideTUIApp` SHALL be registered and resolvable, returning a configured TUI application model

#### Scenario: Dependency resolution chain
- **WHEN** `ProvideTUIApp` is resolved
- **THEN** it SHALL automatically resolve its dependency on `*AppConfig` through the injector

### Requirement: All application dependencies registered in cfg package
The `internal/cfg` package SHALL be the single location where all DI providers are registered, providing a centralized view of the dependency graph.

#### Scenario: New dependency registration
- **WHEN** a new application component is created
- **THEN** its provider function SHALL be added to `internal/cfg/providers.go`
