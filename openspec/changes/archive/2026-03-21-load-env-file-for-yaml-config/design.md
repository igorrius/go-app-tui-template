## Context

The current configuration system loads a YAML file and stores the data in a `map[string]any`. To support environment variable interpolation, we need to:
1. Load `.env` files if they exist before parsing the configuration.
2. Replace `${VAR}` placeholders in the configuration data with their actual values from the environment.

## Goals / Non-Goals

**Goals:**
- Load `.env` file automatically on start.
- Support `${VAR}` and `${VAR:-default}` interpolation in the configuration.
- Maintain existing dot-notation access to configuration values.

**Non-Goals:**
- Interpolating complex objects (only string-based interpolation in scalar values).
- Supporting complex bash-like substitutions (e.g., `${VAR:offset:length}`).

## Decisions

- **Dependency: `github.com/joho/godotenv`**
  - Use `godotenv` to load the `.env` file. It's the standard for Go projects and simplifies the loading logic.
  - *Alternative:* Manual parsing of `.env` files. Rationale: More complex to implement correctly and maintain.
- **Interpolation implementation: `os.Expand` or similar**
  - Use a custom interpolation function or a library like `github.com/drone/envsubst` to handle `${VAR:-default}`.
  - *Alternative:* Use `os.ExpandEnv`. Rationale: `os.ExpandEnv` does not support default values (`${VAR:-default}`).
- **Where to interpolate:**
  - Interpolate values after the YAML is parsed but before it's used in the `AppConfig`. Perform a recursive traversal of the `map[string]any` to find and replace string values.
  - *Alternative:* Interpolate the raw YAML string before parsing. Rationale: Risky if `${VAR}` contains special YAML characters that would break the parser.

## Risks / Trade-offs

- **Risk:** Interpolation could break configuration if not handled carefully.
  - *Mitigation:* Only interpolate string values. Ensure that if interpolation fails or a variable is missing, it behaves predictably (e.g., returns empty string or the placeholder itself).
- **Risk:** Performance impact of recursive traversal.
  - *Mitigation:* Configuration loading happens only once at startup, so the impact is negligible.
