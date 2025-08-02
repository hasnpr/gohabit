# Architectural Decision Records (ADRs)

This directory contains architectural decision records for the GoHabit project. Each decision document captures the context, options considered, and rationale behind significant technical choices made during development.

## Purpose

These decision records serve to:
- **Document rationale** - Capture why specific architectural choices were made
- **Share knowledge** - Help team members understand design decisions
- **Track evolution** - Maintain a history of how the architecture evolved
- **Enable reviews** - Provide context for code reviews and future changes
- **Onboard developers** - Help new team members understand the codebase quickly

## Review Process

All architectural decisions documented here have been reviewed by:
- **Claude AI** - For technical best practices and Go idioms
- **CodeRabbit** - For code quality and maintainability analysis

## Decision Documents

### [ADR-001: Project Structure](./decisions/cmd.md)
**Status:** ✅ Implemented  
**Decision:** Place cmd package at project root, use domain-specific constants

**Key Points:**
- Use `cmd/` at project root following Go conventions
- Avoid global constants package, keep them domain-specific
- Standard Go project layout for scalability

**Benefits:** Community standard, clear separation, maintainable

---

### [ADR-002: Logging Architecture](./decisions/logger.md)
**Status:** ✅ Implemented  
**Decision:** Use embedded `*zap.Logger` struct approach

**Key Points:**
- Embed `*zap.Logger` directly in custom Logger struct
- Provide convenience constructors (`NewDefault`, `NewString`)
- Proper cleanup handling for terminal sync errors

**Benefits:** Clean API, testable, flexible, robust error handling

---

### [ADR-003: Build Information System](./decisions/build_info.md)
**Status:** ✅ Implemented  
**Decision:** Structured build info in `internal/app` with automated injection

**Key Points:**
- Centralized build information in `internal/app/info.go`
- Structured `Info` type with JSON tags
- Automated build script with ldflags injection

**Benefits:** Centralized, API-ready, automated workflow

## Documentation Standards

Each decision document follows this structure:

1. **ADR Title** - Clear, numbered identifier (ADR-XXX)
2. **Status/Date** - Implementation status and date
3. **Context** - Problem being solved
4. **Decision** - Chosen solution with reasoning
5. **Implementation** - Code examples and approach
6. **Alternatives Considered** - Other options evaluated
7. **Benefits/Trade-offs** - Advantages and disadvantages
8. **Review Links** - CodeRabbit PR references

## Adding New Decisions

When documenting new architectural decisions:

1. Create numbered file: `docs/decisions/ADR-XXX-topic.md`
2. Use next available ADR number (004, 005, etc.)
3. Follow the standard template structure
4. Include concise implementation examples
5. Link to related PRs or issues
6. Update this index document

## Related Documentation

- [CLAUDE.md](../CLAUDE.md) - Development guidance for Claude Code
- [README.md](../README.md) - Project overview and setup
- [justfile](../justfile) - Development commands and workflows

---

*This documentation approach helps maintain architectural consistency and provides valuable context for future development decisions.*