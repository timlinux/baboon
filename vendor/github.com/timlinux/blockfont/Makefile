.PHONY: all build clean test run demo fmt lint docs-dev docs-build docs-clean docs-open demo-record demo-play release version help

# Default target
all: build

# Build the library (just verify it compiles)
build:
	go build ./...

# Build with nix (reproducible)
nix-build:
	nix build

# Clean build artifacts
clean:
	rm -rf result
	rm -f demo-bin simple animated editor

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run the demo application
demo: build
	go run ./examples/demo

# Run simple example
run-simple:
	go run ./examples/simple

# Run animated example
run-animated:
	go run ./examples/animated

# Run editor example
run-editor:
	go run ./examples/editor

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Lint code (requires golangci-lint)
lint:
	golangci-lint run ./...

# Vendor dependencies
vendor:
	go mod vendor

# Update dependencies
deps:
	go mod tidy

# Documentation (MkDocs)
docs-dev:
	cd docs && mkdocs serve

docs-build:
	cd docs && mkdocs build

docs: docs-build
	@echo "Documentation built in docs/site/"

docs-clean:
	rm -rf docs/site

docs-open:
	xdg-open http://localhost:8000 2>/dev/null || open http://localhost:8000 2>/dev/null || echo "Open http://localhost:8000 in your browser"

# Demo recording (asciinema)
demo-record:
	nix run .#demo-record

demo-play:
	nix run .#demo-play

# Release management
release:
	nix run .#release

version:
	@grep 'version = ' flake.nix | head -1 | sed 's/.*version = "\([^"]*\)".*/Current version: \1/'
	@echo -n "Latest git tag:  " && (git describe --tags --abbrev=0 2>/dev/null || echo "none")

# Install pre-commit hooks
pre-commit-install:
	pre-commit install

# Run pre-commit on all files
pre-commit:
	pre-commit run --all-files

# Help
help:
	@echo "blockfont - Unicode block letter rendering library"
	@echo ""
	@echo "Build targets:"
	@echo "  make build         - Build/verify the library"
	@echo "  make nix-build     - Build with nix (reproducible)"
	@echo "  make clean         - Remove build artifacts"
	@echo ""
	@echo "Run targets:"
	@echo "  make demo          - Run the interactive demo"
	@echo "  make run-simple    - Run simple example"
	@echo "  make run-animated  - Run animated example"
	@echo "  make run-editor    - Run vim editor example"
	@echo ""
	@echo "Development:"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Lint code"
	@echo "  make vendor        - Vendor dependencies"
	@echo "  make deps          - Update dependencies (go mod tidy)"
	@echo ""
	@echo "Documentation (MkDocs):"
	@echo "  make docs-dev      - Start MkDocs dev server"
	@echo "  make docs-build    - Build documentation"
	@echo "  make docs          - Build documentation (alias)"
	@echo "  make docs-clean    - Remove built documentation"
	@echo "  make docs-open     - Open docs in browser"
	@echo ""
	@echo "Demo Recording (asciinema):"
	@echo "  make demo-record   - Record a terminal demo"
	@echo "  make demo-play     - Play the recorded demo"
	@echo ""
	@echo "Release Management:"
	@echo "  make version       - Show current version and latest tag"
	@echo "  make release       - Interactive version bump and release"
	@echo ""
	@echo "Pre-commit:"
	@echo "  make pre-commit-install - Install pre-commit hooks"
	@echo "  make pre-commit         - Run pre-commit on all files"
