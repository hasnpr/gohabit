# ADR-002: Logging Architecture

**Status:** ✅ Implemented  
**Date:** 2025-01-XX  
**Reviewers:** Claude AI, CodeRabbit

## Context
Need structured logging solution for microservice with proper error handling and testability.

## Decision
Use embedded `*zap.Logger` struct approach instead of global logger variables.

## Implementation
```go
type Logger struct {
    *zap.Logger
}

func New(level zapcore.Level) (*Logger, error)
func NewString(level string) (*Logger, error)  
func NewDefault() *Logger
func (l *Logger) Close() error
```

## Alternatives Considered
- **Global logger with init()**: Rejected due to testing difficulties and inflexibility
- **Direct zap usage**: Rejected due to lack of custom error handling

## Benefits
- **Clean API**: All zap methods directly available
- **Testable**: Easy to mock individual instances
- **Flexible**: Multiple loggers with different levels
- **Robust**: Handles terminal sync errors properly

## Trade-offs
- Slight abstraction overhead vs direct zap usage
- Additional struct vs global simplicity

## Review Links
- [CodeRabbit PR#1](https://github.com/hasnpr/gohabit/pull/1)