# Logging File Rotation

## Purpose

This spec defines the purpose of the Logging File Rotation capability.

## Requirements

### Requirement: Time-based log file rotation
When using the `file` sink, the system SHALL write log entries to files inside the configured directory. Files SHALL be rotated at the configured interval, with each file named using a timestamp that reflects when it was opened (e.g., `app-2006-01-02_15-00.log`). The current log file SHALL be symlinked as `app.log` for convenience.

#### Scenario: New file created after rotation interval
- **WHEN** the rotation interval elapses since the current log file was opened
- **THEN** the system SHALL open a new log file with the current timestamp in its name and route subsequent writes to it

#### Scenario: Log directory is created if absent
- **WHEN** `logging.Init` is called with a `dir` that does not exist
- **THEN** the system SHALL create the directory (including any intermediate directories) before opening the first log file

#### Scenario: Symlink points to current log file
- **WHEN** a new log file is opened due to rotation or initial start
- **THEN** the `app.log` symlink in the log directory SHALL point to the currently active file

### Requirement: Log directory excluded from VCS and Docker
The system SHALL ensure the log directory (default `logs/`) is listed in `.gitignore` and `.dockerignore` so that log files are never accidentally committed or included in container images.

#### Scenario: logs/ ignored by git
- **WHEN** a developer runs `git status` with files present in `logs/`
- **THEN** git SHALL not report those files as untracked or modified

#### Scenario: logs/ excluded from Docker build context
- **WHEN** a Docker image is built from the project root
- **THEN** files inside `logs/` SHALL NOT be copied into the build context
