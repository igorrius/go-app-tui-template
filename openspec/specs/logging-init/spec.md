# Logging Init

## Purpose

This spec defines the purpose of the Logging Init capability.

## Requirements

### Requirement: Logger initialisation at startup
The system SHALL provide a `logging.Init(cfg LoggingConfig, bus *vein.Dispatcher) error` function that creates a `VeinHandler`, starts the `AsyncWriter` goroutine, and calls `slog.SetDefault` with the resulting logger.

#### Scenario: Successful init with stdout sink
- **WHEN** `logging.Init` is called with `sink: "stdout"` and a valid bus
- **THEN** `slog.Default()` SHALL return a logger backed by the `VeinHandler` and any `slog.Info(...)` call SHALL eventually write a line to stdout

#### Scenario: Successful init with file sink
- **WHEN** `logging.Init` is called with `sink: "file"` and a valid bus
- **THEN** `slog.Default()` SHALL return a logger backed by the `VeinHandler` and log entries SHALL appear in the configured log directory

#### Scenario: Init returns error on invalid config
- **WHEN** `logging.Init` is called with an unrecognised `sink` value
- **THEN** the function SHALL return a descriptive error and leave `slog.Default()` unchanged

### Requirement: Logger shutdown drains in-flight events
The system SHALL provide a `logging.Shutdown(ctx context.Context) error` function that signals the `AsyncWriter` to stop, waits for it to drain queued events, and closes any open file descriptors.

#### Scenario: Clean shutdown within deadline
- **WHEN** `logging.Shutdown` is called with a context that has sufficient deadline
- **THEN** all events already published to the bus SHALL be written to the sink before the function returns

#### Scenario: Shutdown respects context cancellation
- **WHEN** `logging.Shutdown` is called with an already-cancelled context
- **THEN** the function SHALL return `context.Canceled` promptly without blocking
