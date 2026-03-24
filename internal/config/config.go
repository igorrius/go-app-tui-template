package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/drone/envsubst/v2"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// AppConfig holds the parsed YAML configuration as a nested map.
type AppConfig struct {
	data map[string]any
}

// DatabaseConfig holds database-specific configuration values.
type DatabaseConfig struct {
	Path string
}

// GetDatabase returns the database configuration from AppConfig.
func GetDatabase(cfg *AppConfig) (DatabaseConfig, error) {
	path, err := Get[string](cfg, "database.path")
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("database config: %w", err)
	}
	return DatabaseConfig{Path: path}, nil
}

// LoggingConfig holds logging-related configuration values.
type LoggingConfig struct {
	Level            slog.Level
	Format           string
	Sink             string
	Dir              string
	RotationInterval time.Duration
}

// GetLogging returns the logging configuration from AppConfig.
// Missing keys resolve to defaults: level=info, format=text, sink=stdout, dir=logs, rotation_interval=1h.
func GetLogging(cfg *AppConfig) (LoggingConfig, error) {
	result := LoggingConfig{
		Level:            slog.LevelInfo,
		Format:           "text",
		Sink:             "stdout",
		Dir:              "logs",
		RotationInterval: time.Hour,
	}

	if level, err := Get[string](cfg, "logging.level"); err == nil {
		parsed, ok := parseSlogLevel(level)
		if !ok {
			return LoggingConfig{}, fmt.Errorf("logging.level: unknown level %q", level)
		}
		result.Level = parsed
	}

	if format, err := Get[string](cfg, "logging.format"); err == nil {
		result.Format = format
	}

	if sink, err := Get[string](cfg, "logging.sink"); err == nil {
		result.Sink = sink
	}

	if dir, err := Get[string](cfg, "logging.dir"); err == nil {
		result.Dir = dir
	}

	if interval, err := Get[string](cfg, "logging.rotation_interval"); err == nil {
		d, err := time.ParseDuration(interval)
		if err != nil {
			return LoggingConfig{}, fmt.Errorf("logging.rotation_interval: %w", err)
		}
		result.RotationInterval = d
	}

	return result, nil
}

func parseSlogLevel(s string) (slog.Level, bool) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}

// Load reads the configuration file at the given path, parsing it into an AppConfig.
// If path is empty, it falls back to app-config.yaml then app-config.dist.yaml.
func Load(path string) (*AppConfig, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	if path != "" {
		return loadFile(path)
	}

	if _, err := os.Stat("app-config.yaml"); err == nil {
		return loadFile("app-config.yaml")
	}

	if _, err := os.Stat("app-config.dist.yaml"); err == nil {
		return loadFile("app-config.dist.yaml")
	}

	return nil, fmt.Errorf("no config file found: tried app-config.yaml and app-config.dist.yaml")
}

func loadFile(path string) (*AppConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	var data map[string]any
	if err := yaml.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	if err := interpolate(data); err != nil {
		return nil, fmt.Errorf("interpolating config %s: %w", path, err)
	}

	return &AppConfig{data: data}, nil
}

func interpolate(data map[string]any) error {
	for k, v := range data {
		switch typedV := v.(type) {
		case string:
			interpolated, err := envsubst.EvalEnv(typedV)
			if err != nil {
				return fmt.Errorf("failed to interpolate key %q: %w", k, err)
			}
			data[k] = interpolated
		case map[string]any:
			if err := interpolate(typedV); err != nil {
				return err
			}
		case []any:
			for i, elem := range typedV {
				if strElem, ok := elem.(string); ok {
					interpolated, err := envsubst.EvalEnv(strElem)
					if err != nil {
						return fmt.Errorf("failed to interpolate index %d in key %q: %w", i, k, err)
					}
					typedV[i] = interpolated
				} else if mapElem, ok := elem.(map[string]any); ok {
					if err := interpolate(mapElem); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// Get retrieves a value at the given dot-separated path and attempts to cast it to type T.
func Get[T any](cfg *AppConfig, path string) (T, error) {
	var zero T

	val, err := traverse(cfg.data, path)
	if err != nil {
		return zero, err
	}

	typed, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("config key %q: expected %T but got %T", path, zero, val)
	}

	return typed, nil
}

// MustGet retrieves a value at the given dot-separated path, panicking on error.
func MustGet[T any](cfg *AppConfig, path string) T {
	val, err := Get[T](cfg, path)
	if err != nil {
		panic(fmt.Sprintf("config: %v", err))
	}

	return val
}

func traverse(data map[string]any, path string) (any, error) {
	parts := strings.Split(path, ".")
	var current any = data

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("config key %q: cannot traverse non-map at %q", path, part)
		}

		val, exists := m[part]
		if !exists {
			return nil, fmt.Errorf("config key %q: not found", path)
		}

		current = val
	}

	return current, nil
}
