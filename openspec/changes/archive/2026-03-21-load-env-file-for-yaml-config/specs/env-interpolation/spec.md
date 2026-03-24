## ADDED Requirements

### Requirement: Interpolate environment variables in YAML config
The system SHALL support interpolating environment variables into YAML configuration values using the `${VAR}` or `${VAR:-default}` syntax.

#### Scenario: Interpolate existing environment variable
- **WHEN** an entry in the config file has the value `"${MY_VAR}"` and `MY_VAR` is set to `"foo"`
- **THEN** the system SHALL return `"foo"` when accessing that config path

#### Scenario: Interpolate with default value
- **WHEN** an entry in the config file has the value `"${MY_VAR:-default_val}"` and `MY_VAR` is not set
- **THEN** the system SHALL return `"default_val"` when accessing that config path

#### Scenario: Missing environment variable without default
- **WHEN** an entry in the config file has the value `"${MY_VAR}"` and `MY_VAR` is not set
- **THEN** the system SHALL treat it as an empty string or return an error depending on the interpolation implementation (ideally empty string to avoid breaking optional configs)

### Requirement: Load .env file on startup
The system SHALL attempt to load a `.env` file from the current working directory upon application startup before parsing the configuration files.

#### Scenario: .env file exists
- **WHEN** a `.env` file exists in the current working directory
- **THEN** its contents SHALL be loaded into the process's environment variables

#### Scenario: .env file does not exist
- **WHEN** no `.env` file exists
- **THEN** the system SHALL continue normally without error
