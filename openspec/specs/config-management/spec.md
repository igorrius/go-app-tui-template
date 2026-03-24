## ADDED Requirements

### Requirement: YAML config file loading
The system SHALL load application configuration from a YAML file, defaulting to `app-config.yaml` with fallback to `app-config.dist.yaml`. Additionally, the system SHALL load a `.env` file if it exists and perform environment variable interpolation on the YAML content. The loaded configuration SHALL include a `logging` section that, when absent, resolves to default logging values.

#### Scenario: Load local config with interpolation
- **WHEN** `app-config.yaml` exists and contains `${VAR}` syntax and a `.env` file exists with `VAR=foo`
- **THEN** the system SHALL parse, interpolate, and use the result as the application configuration

#### Scenario: Fallback to dist config with interpolation
- **WHEN** `app-config.yaml` does not exist, `app-config.dist.yaml` exists and contains `${VAR}` syntax, and a `.env` file exists with `VAR=foo`
- **THEN** the system SHALL parse the dist config, interpolate, and use the result as the application configuration

#### Scenario: No config file found
- **WHEN** neither `app-config.yaml` nor `app-config.dist.yaml` exists
- **THEN** the system SHALL return an error during initialization

#### Scenario: Logging defaults applied when section absent
- **WHEN** the YAML config does not contain a `logging` section
- **THEN** the loaded `LoggingConfig` SHALL have `level: info`, `format: text`, `sink: stdout`

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
