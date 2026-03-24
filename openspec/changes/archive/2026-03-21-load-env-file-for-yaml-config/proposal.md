## Why

The current configuration management system only supports static YAML file loading. There is no support for environment variables or `.env` files, which are essential for managing sensitive information (like API keys) and environment-specific changes without modifying the YAML files directly.

## What Changes

- **Environment Variable Interpolation**: Support `${VAR}` or `${VAR:-default}` syntax in YAML configuration files to inject environment variables.
- **Automatic .env Loading**: The application will automatically look for and load a `.env` file (if present) upon startup to populate environment variables.
- **Config Management Update**: Update the configuration loading logic to perform interpolation after parsing the YAML structure.

## Capabilities

### New Capabilities
- `env-interpolation`: The system SHALL support interpolating environment variables into configuration values using `${VAR}` syntax.

### Modified Capabilities
- `config-management`: The configuration loading process now includes loading `.env` files and performing environment variable expansion on the loaded YAML content.

## Impact

- `internal/cfg/config.go`: Loading logic needs to be updated.
- `internal/config/config.go`: May need updates if it handles the low-level parsing.
- New dependency on a `.env` loading library (e.g., `github.com/joho/godotenv`) or manual implementation.
