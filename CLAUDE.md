# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library project called `gohabit` that provides logging functionality. The project is in early stages with minimal components.

## Module Information

- **Module**: `github.com/hasnpr/gohabit`
- **Go Version**: 1.24.3
- **Main Dependencies**: 
  - `go.uber.org/zap v1.27.0` (structured logging)
  - `go.uber.org/multierr v1.10.0` (indirect)

## Development Commands

### Build and Test
```bash
# Build the module
go build ./...

# Run tests (when they exist)
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Tidy module dependencies
go mod tidy
```

### Linting
No specific linter configuration found. Use standard Go tools:
```bash
go fmt ./...
go vet ./...
```

## Code Architecture

### Package Structure
- `pkg/logger/` - Logging utilities using Uber Zap
  - Provides structured JSON logging to stdout/stderr
  - Configurable log levels via `zap.AtomicLevel`
  - Custom time format: "2006-01-02 15:04:05.000"
  - Production-ready configuration with stacktrace and caller info disabled

### Key Components

#### Logger Package (`pkg/logger/log.go`)
- **Global Variables**: `defaultLogger` and `DefaultLevel` for application-wide logging
- **Configuration**: JSON encoding, production settings, stdout/stderr outputs
- **Cleanup**: `CloseLogger()` function handles proper logger cleanup with error handling for TTY/file descriptor issues

## Development Notes

- This is a library project with no main application entry point
- No test files or build scripts currently exist
- Uses standard Go module structure
- Logger is initialized in `init()` function with production defaults
- Error handling includes specific checks for `syscall.ENOTTY` and `syscall.EBADF` in logger cleanup