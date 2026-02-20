# Contributing to Baboon

Thank you for your interest in contributing to Baboon! This guide will help you get started.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Node.js 18+ (for web frontend)
- Git
- Optional: Nix (for reproducible builds)

### Setting Up the Development Environment

#### Option 1: Nix (Recommended)

```bash
# Clone the repository
git clone https://github.com/timlinux/baboon.git
cd baboon

# Enter the development shell
nix develop

# All dependencies are now available!
```

#### Option 2: Manual Setup

```bash
# Clone the repository
git clone https://github.com/timlinux/baboon.git
cd baboon

# Install Go dependencies
go mod download

# Install web dependencies
cd web
npm install
cd ..
```

## Project Structure

```
baboon/
├── main.go              # Application entry point
├── backend/             # Game engine and REST API
│   ├── api.go          # GameAPI interface
│   ├── engine.go       # Game logic
│   └── server.go       # REST server
├── frontend/            # Terminal UI
│   ├── model.go        # Bubble Tea model
│   ├── views.go        # Rendering
│   ├── styles.go       # Lipgloss styles
│   ├── animations.go   # Spring animations
│   └── client.go       # REST client
├── font/                # Block letter font
│   └── font.go
├── words/               # Word dictionary
│   └── words.go
├── stats/               # Statistics
│   ├── stats.go        # Types and persistence
│   └── keyboard.go     # Layout mappings
├── web/                 # React web frontend
│   └── src/
├── scripts/             # Management scripts
└── docs/                # Documentation (MkDocs)
```

## Development Workflow

### Building

```bash
# Build the binary
go build -o baboon .

# Or use make
make build
```

### Running

```bash
# Run directly
go run .

# With punctuation mode
go run . -p

# Server mode
go run . -server
```

### Testing

```bash
# Run all tests
go test ./...

# With verbose output
go test -v ./...

# Run specific package tests
go test ./stats/...
```

### Code Formatting

```bash
# Format Go code
go fmt ./...

# Or use make
make fmt
```

### Linting

```bash
# Run go vet
go vet ./...
```

## Making Changes

### Branch Naming

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation
- `refactor/description` - Code refactoring

### Code Style

- Follow standard Go conventions
- Use meaningful variable names
- Add comments for complex logic
- Keep functions focused and small

### British English

Use British spellings throughout:

- colour (not color)
- behaviour (not behavior)
- centre (not center)

### Commit Messages

Write clear, concise commit messages:

```
Add per-finger accuracy tracking

- Track accuracy for each of 8 fingers
- Display in results screen with colour coding
- Store in stats.json for persistence
```

## Areas for Contribution

### Good First Issues

- Documentation improvements
- Adding new words to the dictionary
- UI tweaks and improvements
- Test coverage improvements

### Feature Ideas

- Multiple word lists (programming, languages, etc.)
- Practice modes (timed, accuracy, etc.)
- Leaderboards
- Custom themes
- Sound effects
- Accessibility improvements

### Bug Fixes

Check the [Issues](https://github.com/timlinux/baboon/issues) page for known bugs.

## Pull Request Process

### Before Submitting

1. **Test your changes**:
   ```bash
   go test ./...
   go vet ./...
   ```

2. **Format your code**:
   ```bash
   go fmt ./...
   ```

3. **Update documentation** if needed

4. **Update SPECIFICATION.md** if behaviour changes

### Submitting

1. Push your branch to your fork
2. Open a Pull Request against `main`
3. Fill out the PR template
4. Wait for review

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation
- [ ] Refactoring

## Testing
How did you test your changes?

## Checklist
- [ ] Tests pass
- [ ] Code is formatted
- [ ] Documentation updated
- [ ] SPECIFICATION.md updated (if applicable)
```

## Architecture Guidelines

### Backend/Frontend Separation

All game logic lives in `backend/`:

- Frontend communicates only through `GameAPI` interface
- No direct access to game state
- Clean API boundaries

### Timing

All timing-critical operations happen on the frontend:

- WPM calculations
- Seek time measurements
- Round duration tracking

This eliminates network latency effects.

### Statistics

Statistics are managed by the `stats` package:

- Persistence to JSON
- Validation on load
- Cumulative tracking across sessions

## Web Frontend

### Technology

- React 18 with hooks
- Chakra UI for components
- Framer Motion for animations

### Development

```bash
# Start dev server
cd web
npm start

# Build for production
npm run build
```

### Adding Components

1. Create component in `web/src/components/`
2. Use Chakra UI components where possible
3. Follow existing patterns for animations

## Documentation

### Building Docs

```bash
# Install MkDocs
pip install mkdocs-material

# Serve locally
mkdocs serve

# Build static site
mkdocs build
```

### Adding Pages

1. Create `.md` file in `docs/`
2. Add to `nav` in `mkdocs.yml`
3. Follow existing formatting

## Getting Help

- **Questions**: Open a Discussion on GitHub
- **Bugs**: Open an Issue
- **Ideas**: Open a Discussion or Issue

## Code of Conduct

Be respectful and constructive. We're all here to make Baboon better!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for helping make Baboon awesome! 🎉
