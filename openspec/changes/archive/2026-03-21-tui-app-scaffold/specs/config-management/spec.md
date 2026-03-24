## ADDED Requirements

### Requirement: YAML config file loading
The system SHALL load application configuration from a YAML file, defaulting to `app-config.yaml` with fallback to `app-config.dist.yaml`.

#### Scenario: Load local config
- **WHEN** `app-config.yaml` exists
- **THEN** the system SHALL parse and use it as the application configuration

#### Scenario: Fallback to dist config
- **WHEN** `app-config.yaml` does not exist and `app-config.dist.yaml` exists
- **THEN** the system SHALL parse and use `app-config.dist.yaml` as the application configuration

#### Scenario: No config file found
- **WHEN** neither `app-config.yaml` nor `app-config.dist.yaml` exists
- **THEN** the system SHALL return an error during initialization

### Requirement: Generic typed config access via dot-path
The system SHALL provide generic functions `Get[T](path string) (T, error)` and `MustGet[T](path string) T` to access nested config values using dot-notation paths.

#### Scenario: Get string value
- **WHEN** `Config.Get[string]("app.name")` is called and the path exists with a string value
- **THEN** the function SHALL return the string value and nil error

#### Scenario: Get int value
- **WHEN** `Config.Get[int]("server.port")` is called and the path exists with an integer value
- **THEN** the function SHALL return the integer value and nil error

#### Scenario: Get with missing path
- **WHEN** `Config.Get[T]("nonexistent.path")` is called
- **THEN** the function SHALL return the zero value of T and a descriptive error

#### Scenario: Get with type mismatch
- **WHEN** `Config.Get[int]("app.name")` is called but the value is a string
- **THEN** the function SHALL return the zero value of int and a type mismatch error

#### Scenario: MustGet panics on missing path
- **WHEN** `Config.MustGet[T]("nonexistent.path")` is called
- **THEN** the function SHALL panic with a descriptive message

#### Scenario: MustGet returns value on valid path
- **WHEN** `Config.MustGet[string]("app.name")` is called and the path exists
- **THEN** the function SHALL return the value directly without error
