# ADR-003: Build Information System

**Status:** ✅ Implemented  
**Date:** 2025-01-XX  
**Reviewers:** Claude AI, CodeRabbit

## Context
Need build information (version, git commit, build date) for debugging and deployment tracking.

## Decision
Keep build info in `internal/app` package with structured approach and automated injection.

## Implementation
```go
// internal/app/info.go
type Info struct {
    Name            string `json:"name"`
    Version         string `json:"version"`
    GitCommit       string `json:"git_commit"`
    GitRef          string `json:"git_ref"`
    GitTag          string `json:"git_tag"`
    BuildDate       string `json:"build_date"`
    CompilerVersion string `json:"compiler_version"`
}

func GetInfo() Info
```

## Build System
```bash
# scripts/build.sh with ldflags injection
go build -ldflags "\
    -X 'github.com/hasnpr/gohabit/internal/app.Version=${VERSION}' \
    -X 'github.com/hasnpr/gohabit/internal/app.GitCommit=${GIT_COMMIT}' \
    ..." -o bin/gohabit .
```

## Alternatives Considered
- **pkg/build_info**: Rejected - not reusable library code
- **internal/build_info**: Rejected - unnecessary separation
- **Global variables**: Rejected - harder to test and organize

## Benefits
- **Centralized**: Single source of build information
- **Structured**: JSON-ready for API endpoints
- **Automated**: Build script handles injection
- **Accessible**: Easy banner and version command integration

## Trade-offs
- Build script dependency vs manual version management
- ldflags complexity vs simple constants

## Usage
```go
info := app.GetInfo()
fmt.Printf("Version: %s\n", info.Version)
```