# Building Baboon

Instructions for building Baboon from source.

## Prerequisites

### Required

- **Go 1.21+**: [Download](https://go.dev/dl/)
- **Git**: For cloning the repository

### Optional

- **Nix**: For reproducible builds
- **Node.js 18+**: For web frontend
- **Make**: For convenience targets

## Quick Build

### Clone and Build

```bash
git clone https://github.com/timlinux/baboon.git
cd baboon
go build -o baboon .
```

### Run

```bash
./baboon
```

## Using Nix

Nix provides the most reproducible build environment.

### Enter Development Shell

```bash
nix develop
```

This provides:

- Go compiler
- All Go dependencies
- Node.js and npm
- Development tools

### Build with Nix

```bash
nix build
./result/bin/baboon
```

### Run Directly

```bash
nix run github:timlinux/baboon
```

## Using Make

The Makefile provides convenient targets:

### Build Targets

```bash
# Build binary
make build

# Run directly
make run

# Run with punctuation mode
make run-p
```

### Test Targets

```bash
# Run all tests
make test

# Run with verbose output
make test-v
```

### Format Targets

```bash
# Format Go code
make fmt

# Run linter
make vet
```

### Web Frontend

```bash
# Install npm dependencies
make web-install

# Start development server
make web-dev

# Build for production
make web-build

# Start backend + web together
make web-start
```

## Cross-Compilation

Build for multiple platforms:

### Linux

```bash
# AMD64
GOOS=linux GOARCH=amd64 go build -o baboon-linux-amd64 .

# ARM64
GOOS=linux GOARCH=arm64 go build -o baboon-linux-arm64 .
```

### macOS

```bash
# Intel
GOOS=darwin GOARCH=amd64 go build -o baboon-darwin-amd64 .

# Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o baboon-darwin-arm64 .
```

### Windows

```bash
GOOS=windows GOARCH=amd64 go build -o baboon-windows-amd64.exe .
```

## Package Building

### DEB Package (Debian/Ubuntu)

```bash
# Create package structure
mkdir -p dist/baboon_1.0.0_amd64/DEBIAN
mkdir -p dist/baboon_1.0.0_amd64/usr/bin

# Copy binary
cp baboon dist/baboon_1.0.0_amd64/usr/bin/

# Create control file
cat > dist/baboon_1.0.0_amd64/DEBIAN/control << EOF
Package: baboon
Version: 1.0.0
Architecture: amd64
Maintainer: Tim Sutton <tim@example.com>
Description: Typing practice application
EOF

# Build package
dpkg-deb --build dist/baboon_1.0.0_amd64
```

### RPM Package (Fedora/RHEL)

Create `baboon.spec`:

```spec
Name:           baboon
Version:        1.0.0
Release:        1%{?dist}
Summary:        Typing practice application

License:        MIT
URL:            https://github.com/timlinux/baboon

%description
A cross-platform typing practice application.

%install
mkdir -p %{buildroot}/usr/bin
install -m 755 baboon %{buildroot}/usr/bin/

%files
/usr/bin/baboon
```

Build with `rpmbuild`.

## Build Optimisations

### Smaller Binary

```bash
# Strip debug information
go build -ldflags="-s -w" -o baboon .
```

### Static Binary

```bash
CGO_ENABLED=0 go build -o baboon .
```

### Version Information

```bash
go build -ldflags="-X main.Version=1.0.0" -o baboon .
```

## Web Frontend Build

### Development

```bash
cd web
npm install
npm start
```

Opens http://localhost:3000 with hot reload.

### Production

```bash
cd web
npm run build
```

Creates optimised build in `web/build/`.

### Serving Production Build

```bash
# Using npx serve
npx serve -s web/build -l 3000

# Or copy to a web server
cp -r web/build/* /var/www/baboon/
```

## Testing the Build

### Verify Binary

```bash
./baboon --help
```

Expected output:

```
Baboon - Typing Practice Application

Usage:
  baboon [flags]

Flags:
  -p              Enable punctuation mode
  -port int       Server port (default 8787)
  -server         Run in server-only mode
  -client         Run in client-only mode
```

### Run Tests

```bash
go test ./...
```

### Check for Race Conditions

```bash
go test -race ./...
```

## Continuous Integration

### GitHub Actions

The repository includes CI workflows:

- **test.yml**: Runs tests on push/PR
- **build.yml**: Cross-platform build verification
- **release.yml**: Automated releases on tag push

### Local CI Simulation

```bash
# Run what CI runs
go test ./...
go vet ./...
go build .
```

## Troubleshooting

### Missing Dependencies

```bash
# Update Go modules
go mod download

# Verify modules
go mod verify
```

### Build Errors

```bash
# Clean build cache
go clean -cache

# Rebuild
go build .
```

### Web Build Issues

```bash
# Clear npm cache
cd web
rm -rf node_modules
npm cache clean --force
npm install
```

## Release Process

1. Update version in code
2. Update SPECIFICATION.md
3. Update CHANGELOG
4. Create git tag: `git tag v1.0.0`
5. Push tag: `git push origin v1.0.0`
6. GitHub Actions builds and publishes release

## Next Steps

- [Architecture](architecture.md) - System design
- [API Reference](api.md) - REST API docs
- [Contributing](contributing.md) - How to contribute
