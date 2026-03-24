## 1. Setup

- [x] 1.1 Add `github.com/joho/godotenv` to `go.mod`
- [x] 1.2 Add `github.com/drone/envsubst` (or similar) for interpolation to `go.mod`

## 2. Environment Loading

- [x] 2.1 Implement `.env` file loading in `internal/config/config.go`
- [x] 2.2 Update `Load` function to call `.env` loading before parsing YAML

## 3. Interpolation Logic

- [x] 3.1 Implement recursive interpolation for `map[string]any` in `internal/config/config.go`
- [x] 3.2 Add support for `${VAR}` and `${VAR:-default}` syntax
- [x] 3.3 Ensure interpolation is applied to the YAML content after parsing but before creating `AppConfig`

## 4. Verification

- [x] 4.1 Verify `.env` file values are correctly interpolated into YAML configuration
- [x] 4.2 Verify default values are used when environment variables are missing
- [x] 4.3 Verify that fallback to `app-config.dist.yaml` still works with interpolation
