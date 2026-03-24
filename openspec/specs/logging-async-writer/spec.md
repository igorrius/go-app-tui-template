## ADDED Requirements

### Requirement: AsyncWriter subscribes to log events and writes to sink
The system SHALL provide an `AsyncWriter` that subscribes to the `go-vein` bus topic for `LogEvent` messages and for each event writes a log line to the configured sink using a standard `slog.TextHandler` or `slog.JSONHandler`.

#### Scenario: Text format writes human-readable line
- **WHEN** the `AsyncWriter` is configured with `format: "text"` and receives a `LogEvent`
- **THEN** it SHALL write a plain-text log line (as produced by `slog.NewTextHandler`) to the configured sink

#### Scenario: JSON format writes JSON line
- **WHEN** the `AsyncWriter` is configured with `format: "json"` and receives a `LogEvent`
- **THEN** it SHALL write a JSON log line (as produced by `slog.NewJSONHandler`) to the configured sink

#### Scenario: Errors from the sink writer are non-fatal
- **WHEN** the sink writer returns an error on a write call
- **THEN** the `AsyncWriter` SHALL continue processing subsequent events and SHALL NOT panic or stop the process

### Requirement: AsyncWriter runs as a managed goroutine
The `AsyncWriter` SHALL be started with `Start(ctx context.Context)` and stop cleanly when the context is cancelled after draining any remaining events from the bus subscription.

#### Scenario: Writer stops after context cancellation
- **WHEN** the context passed to `Start` is cancelled
- **THEN** the goroutine SHALL process remaining queued events and then exit

#### Scenario: Writer is concurrent-safe
- **WHEN** multiple goroutines emit log records simultaneously
- **THEN** all records SHALL be written without data races or interleaved partial lines
