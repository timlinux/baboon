.PHONY: all build clean test run server client install help

# Default target
all: build

# Build the binary
build:
	go build -o baboon .

# Build with nix (reproducible)
nix-build:
	nix build

# Clean build artifacts
clean:
	rm -f baboon
	rm -rf result

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

# Help
help:
	@echo "Baboon - Terminal typing practice"
	@echo ""
	@echo "Build targets:"
	@echo "  make build       - Build the binary"
	@echo "  make nix-build   - Build with nix (reproducible)"
	@echo "  make clean       - Remove build artifacts"
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
	@echo "Development:"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Lint code"
	@echo "  make vendor      - Vendor dependencies"
	@echo "  make deps        - Update dependencies"
