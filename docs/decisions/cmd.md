# ADR-001: Project Structure

**Status:** ✅ Implemented
**Date:** 2025-01-XX
**Reviewers:** Claude AI, CodeRabbit

## Context
Need to organize Go project structure for microservice with proper separation of concerns.

## Decision
- Place `cmd/` package at project root following Go conventions
- Keep constants domain-specific, avoid global constants package

## Implementation
```
gohabit/
├── cmd/           # CLI commands at root
├── internal/      # Private application code
├── pkg/           # Public library code
└── main.go        # Application entry point
```

## Alternatives Considered
- **internal/cmd/**: Rejected - mixing entry points with internal code
- **Global constants package**: Rejected - creates tight coupling

## Benefits
- **Standard Go layout**: Follows community conventions
- **Scalable**: Easy to add multiple binaries
- **Clear separation**: Commands vs library vs internal code
- **Maintainable**: Domain-specific constants reduce coupling

## Trade-offs
- More directories vs flat structure
- Need to understand Go project conventions

## Constants Strategy
```
// Domain-specific constants
internal/auth/constants.go
internal/handlers/constants.go

// Application-level constants
internal/app/constants.go

// Shared constants (minimal)
internal/shared/errors.go
```

## Review Links
- [CodeRabbit PR#2](https://github.com/hasnpr/gohabit/pull/2)
