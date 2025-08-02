# Justfile for gohabit project
# Run commands with: just <command>

# Default command - show available commands
default:
    @just --list

# Build the application with version information
build:
    @echo "Building gohabit..."
    @./scripts/build.sh

# Run the application
run:
    @go run .

# Run the application with banner
run-banner:
    @go run . start --show-banner

# Run tests
test:
    @go test ./...

# Run tests with coverage
test-coverage:
    @go test -cover ./...

# Format code
fmt:
    @go fmt ./...

# Vet code for issues
vet:
    @go vet ./...

# Lint and format (combination of fmt and vet)
lint: fmt vet

# Clean build artifacts
clean:
    @rm -rf bin/
    @echo "Cleaned build artifacts"

# Install dependencies
deps:
    @go mod download
    @go mod tidy

# Development build and run
dev: build
    @./bin/gohabit start --show-banner

# Show build information without building
info:
    @echo "Project: gohabit"
    @echo "Go version: $(go version)"
    @echo "Git branch: $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')"
    @echo "Git commit: $(git rev-parse HEAD 2>/dev/null || echo 'unknown')"

# Run all checks (test, lint, vet)
check: test lint vet
    @echo "All checks passed ✓"