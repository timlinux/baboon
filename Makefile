.PHONY: all build clean test run server client install help web-install web-dev web-build web-start docs docs-dev docs-build docs-clean docs-open docs-new web-bundle web-serve demo-record demo-play release version tutorial-check tutorial-check-verbose tutorial-extract tutorial-serve tutorial-new-chapter pre-commit-install

# Default target - build everything
all: build web-build docs-build
	@echo "Build complete: Go binary, web frontend, and Hugo docs"

# Extract version from flake.nix (single source of truth)
VERSION := $(shell grep 'baboonVersion = "' flake.nix | sed 's/.*baboonVersion = "\([^"]*\)".*/\1/')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)

# Build only the Go binary
build:
	go build -ldflags '$(LDFLAGS)' -o baboon .

# Build with nix (reproducible)
nix-build:
	nix build

# Clean build artifacts and caches
clean:
	rm -f baboon
	rm -rf result
	go clean -cache -testcache
	rm -rf web/node_modules web/build
	rm -rf hugo/public hugo/resources hugo/.hugo_build.lock
	@echo "All caches cleared. Next build will be fresh."

# Run tests
test:
	go test ./...

# Run in combined mode (default)
run: build
	./baboon

# Run with punctuation mode
run-p: build
	./baboon -p

# Run server only
server: build
	./baboon -server

# Run client only (connect to existing server)
client: build
	./baboon -client

# Install to GOPATH/bin
install:
	go install .

# Vendor dependencies
vendor:
	go mod vendor

# Update dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Start backend in background
start-backend: build
	./scripts/start-backend.sh

# Stop backend
stop-backend:
	./scripts/stop-backend.sh

# Check backend status
status:
	./scripts/status-backend.sh

# Launch frontend against running backend
frontend: build
	./scripts/launch-frontend.sh

# Web frontend targets
web-install:
	cd web && npm install

web-dev: web-install
	cd web && npm run dev

web-build: web-install
	cd web && npm run build

web-start: build
	@echo "Starting backend in background..."
	./baboon -server &
	@sleep 2
	@echo "Starting web frontend..."
	cd web && npm run dev

# Documentation (Hugo)
docs-dev:
	cd hugo && hugo server -D --bind 0.0.0.0

docs-build:
	cd hugo && hugo --minify

docs: docs-build
	@echo "Documentation built in hugo/public/"

docs-clean:
	rm -rf hugo/public hugo/resources web/build/docs

docs-open:
	xdg-open http://localhost:1313 2>/dev/null || open http://localhost:1313 2>/dev/null || echo "Open http://localhost:1313 in your browser"

docs-new:
	@read -p "Enter page path (e.g., posts/my-new-post): " path; \
	cd hugo && hugo new "$$path.md"

# Tutorial development
tutorial-check: ## Validate all tutorial references
	go run ./tools/tutorialcheck --all

tutorial-check-verbose: ## Validate all tutorial references (verbose)
	go run ./tools/tutorialcheck --all --verbose

tutorial-extract: ## Re-extract all code snippets (debug)
	go run ./tools/tutorialcheck --extract-only --verbose

tutorial-serve: ## Hugo server with live reload for tutorial
	cd hugo && hugo server -D --port 1314

tutorial-new-chapter: ## Scaffold a new chapter (interactive)
	@read -p "Chapter number (e.g., 13): " num; \
	read -p "Chapter name (e.g., testing): " name; \
	mkdir -p hugo/content/go-tutorial/$$num-$$name; \
	echo "---" > hugo/content/go-tutorial/$$num-$$name/_index.md; \
	echo "title: \"Chapter $$num: $$(echo $$name | sed 's/-/ /g' | sed 's/\b\(.\)/\u\1/g')\"" >> hugo/content/go-tutorial/$$num-$$name/_index.md; \
	echo "weight: $$num" >> hugo/content/go-tutorial/$$num-$$name/_index.md; \
	echo "---" >> hugo/content/go-tutorial/$$num-$$name/_index.md; \
	echo "" >> hugo/content/go-tutorial/$$num-$$name/_index.md; \
	echo "Created hugo/content/go-tutorial/$$num-$$name/"

pre-commit-install: ## Install pre-commit hooks
	pip install pre-commit
	pre-commit install
	@echo "Pre-commit hooks installed!"

# Combined production build (React + Hugo bundled together)
web-bundle: web-build docs-build
	@echo "Bundling Hugo docs into web/build/docs/..."
	rm -rf web/build/docs
	cp -r hugo/public web/build/docs
	@echo "Production bundle complete in web/build/"

# Run production server with bundled docs
web-serve: web-bundle build
	./baboon web -port 8787 -dir web/build

# Demo recording (asciinema)
demo-record:
	nix run .#demo-record

demo-play:
	nix run .#demo-play

# Release management
release:
	nix run .#release

version:
	@echo "Current version: $(VERSION)"
	@echo -n "Latest git tag:  " && (git describe --tags --abbrev=0 2>/dev/null || echo "none")

# Help
help:
	@echo "Baboon - Terminal typing practice"
	@echo ""
	@echo "Build targets:"
	@echo "  make             - Build everything (Go, web, docs)"
	@echo "  make build       - Build Go binary only"
	@echo "  make nix-build   - Build with nix (reproducible)"
	@echo "  make clean       - Remove all build artifacts and caches"
	@echo "  make install     - Install to GOPATH/bin"
	@echo ""
	@echo "Run targets:"
	@echo "  make run         - Run in combined mode"
	@echo "  make run-p       - Run with punctuation mode"
	@echo "  make server      - Run backend only"
	@echo "  make client      - Run frontend only"
	@echo ""
	@echo "Backend management:"
	@echo "  make start-backend - Start backend in background"
	@echo "  make stop-backend  - Stop backend"
	@echo "  make status        - Check backend status"
	@echo "  make frontend      - Launch frontend client"
	@echo ""
	@echo "Web frontend:"
	@echo "  make web-install - Install web dependencies"
	@echo "  make web-dev     - Start web dev server"
	@echo "  make web-build   - Build web for production"
	@echo "  make web-start   - Start backend + web frontend"
	@echo "  make web-bundle  - Build React + Hugo bundled together"
	@echo "  make web-serve   - Run production server with docs"
	@echo ""
	@echo "Development:"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Lint code"
	@echo "  make vendor      - Vendor dependencies"
	@echo "  make deps        - Update dependencies"
	@echo ""
	@echo "Documentation (Hugo):"
	@echo "  make docs-dev    - Start Hugo dev server"
	@echo "  make docs-build  - Build documentation"
	@echo "  make docs        - Build documentation (alias)"
	@echo "  make docs-clean  - Remove built documentation"
	@echo "  make docs-open   - Open docs in browser"
	@echo "  make docs-new    - Create new documentation page"
	@echo ""
	@echo "Tutorial Development:"
	@echo "  make tutorial-check   - Validate tutorial references"
	@echo "  make tutorial-check-verbose - Verbose validation"
	@echo "  make tutorial-extract - Extract code snippets (debug)"
	@echo "  make tutorial-serve   - Hugo server for tutorial"
	@echo "  make tutorial-new-chapter - Scaffold new chapter"
	@echo "  make pre-commit-install - Install pre-commit hooks"
	@echo ""
	@echo "Demo Recording (asciinema):"
	@echo "  make demo-record - Record a terminal demo"
	@echo "  make demo-play   - Play the recorded demo"
	@echo ""
	@echo "Release Management:"
	@echo "  make version     - Show current version and latest tag"
	@echo "  make release     - Interactive version bump and release"
