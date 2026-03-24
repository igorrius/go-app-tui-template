## MODIFIED Requirements

### Requirement: YAML config file loading
The system SHALL load application configuration from a YAML file, defaulting to `app-config.yaml` with fallback to `app-config.dist.yaml`. Additionally, the system SHALL load a `.env` file if it exists and perform environment variable interpolation on the YAML content.

#### Scenario: Load local config with interpolation
- **WHEN** `app-config.yaml` exists and contains `${VAR}` syntax and a `.env` file exists with `VAR=foo`
- **THEN** the system SHALL parse, interpolate, and use the result as the application configuration

#### Scenario: Fallback to dist config with interpolation
- **WHEN** `app-config.yaml` does not exist, `app-config.dist.yaml` exists and contains `${VAR}` syntax, and a `.env` file exists with `VAR=foo`
- **THEN** the system SHALL parse the dist config, interpolate, and use the result as the application configuration

#### Scenario: No config file found
- **WHEN** neither `app-config.yaml` nor `app-config.dist.yaml` exists
- **THEN** the system SHALL return an error during initialization
