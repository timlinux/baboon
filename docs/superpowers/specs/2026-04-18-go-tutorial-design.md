# Go Tutorial Design Specification

**Date:** 2026-04-18
**Status:** Approved
**Author:** Claude Code

## Overview

A comprehensive Go programming tutorial integrated into Baboon's Hugo documentation site. The tutorial uses Baboon's codebase as the primary learning vehicle, walking through packages in order of conceptual complexity while teaching Go concepts as they naturally appear in production code.

## Target Audience

Developers with programming experience in other languages (Python, JavaScript, etc.) who are new to Go. The tutorial assumes familiarity with general programming concepts but not Go-specific syntax or idioms.

## Learning Objectives

By completing this tutorial, learners will:
- Understand Go syntax, types, and idioms
- Be able to read and modify the Baboon codebase confidently
- Know how to structure Go projects with clean package boundaries
- Understand Go's approach to interfaces, error handling, and concurrency
- Be familiar with common Go patterns for HTTP servers, database access, and authentication

---

## Tutorial Structure

### Location

`hugo/content/go-tutorial/` within the existing Hugo documentation site.

### Chapter Organization

The tutorial follows a **Progressive Package Tour** approach, walking through Baboon's packages in order of conceptual complexity:

| Chapter | Package(s) | Go Concepts |
|---------|-----------|-------------|
| 01 | Getting Started | Installation, tooling, project structure |
| 02 | `settings/` | Basic types, structs, JSON tags, file I/O |
| 03 | `words/` | Slices, maps, randomization |
| 04 | `stats/` | Complex structs, methods, persistence |
| 05 | `backend/api.go` | Interfaces, type definitions |
| 06 | `backend/engine.go` | Implementing interfaces, business logic |
| 07 | `frontend/` | Bubble Tea TUI, Model-Update-View |
| 08 | `backend/server.go` | HTTP servers, routing, concurrency |
| 09 | `frontend/client.go` | HTTP clients, interface over network |
| 10 | `database/` | SQL, repository pattern, migrations |
| 11 | `auth/` | OAuth2, JWT, middleware |
| 12 | `main.go` | CLI flags, signals, dependency wiring |
| Appendix | - | Tooling, patterns summary, further reading |

### Chapter Format

Each chapter follows a consistent structure:

1. **Introduction** - What this package does in Baboon, why it exists
2. **Concepts** - Go concepts demonstrated, with code extracted from the package
3. **Deep Dive** - Detailed walkthrough of key functions and design decisions
4. **Exercises** - Hands-on modifications with collapsible solutions

### File Structure

```
hugo/content/go-tutorial/
├── _index.md
├── 01-getting-started/
│   ├── _index.md
│   └── hello-baboon.md
├── 02-settings/
│   ├── _index.md
│   ├── basic-types.md
│   ├── structs.md
│   ├── json-tags.md
│   └── file-io.md
├── 03-words/
│   ├── _index.md
│   ├── slices.md
│   ├── maps.md
│   └── randomization.md
├── 04-stats/
│   ├── _index.md
│   ├── complex-structs.md
│   ├── methods.md
│   └── persistence.md
├── 05-interfaces/
│   ├── _index.md
│   ├── defining.md
│   ├── implementing.md
│   └── polymorphism.md
├── 06-game-engine/
│   ├── _index.md
│   ├── state-management.md
│   ├── algorithms.md
│   └── business-logic.md
├── 07-terminal-ui/
│   ├── _index.md
│   ├── bubble-tea.md
│   ├── styling.md
│   └── animations.md
├── 08-http-server/
│   ├── _index.md
│   ├── routing.md
│   ├── middleware.md
│   └── concurrency.md
├── 09-http-client/
│   ├── _index.md
│   ├── requests.md
│   ├── json-api.md
│   └── interface-over-http.md
├── 10-database/
│   ├── _index.md
│   ├── sql-basics.md
│   ├── repository-pattern.md
│   └── migrations.md
├── 11-authentication/
│   ├── _index.md
│   ├── oauth2.md
│   ├── jwt.md
│   └── middleware.md
├── 12-main-and-cli/
│   ├── _index.md
│   ├── flags.md
│   ├── signals.md
│   └── wiring.md
└── appendix/
    ├── go-tooling.md
    ├── common-patterns.md
    └── further-reading.md
```

---

## Validation Tooling

### Purpose

Prevent tutorial rot by automatically validating that code references, extracted snippets, and standalone examples remain accurate as the Baboon codebase evolves.

### Validation Mechanisms

#### 1. Code Reference Markers

Source files contain special comments marking code for tutorial extraction:

```go
// TUTORIAL:settings/structs:start
type Settings struct {
    AdvanceKey  AdvanceKey `json:"advance_key"`
    PerfectMode bool       `json:"perfect_mode"`
}
// TUTORIAL:settings/structs:end
```

Tutorial markdown references these markers:

```markdown
{{</* code-ref id="settings/structs" */>}}
```

The Hugo shortcode extracts current code at build time, ensuring snippets are always fresh.

#### 2. Line Reference Validation

Inline references like `settings.go:47` are validated to ensure:
- File exists at specified path
- Line number is within file bounds
- Optionally: content at line matches expected pattern

#### 3. Standalone Example Compilation

Code blocks tagged with `compile` are extracted and compiled:

~~~markdown
```go compile
package main

func main() {
    println("Hello")
}
```
~~~

The validator extracts to temp files and runs `go build` to verify compilation.

### Tool Structure

```
tools/
└── tutorialcheck/
    ├── main.go           # CLI entry point
    ├── extract.go        # Code extraction from markers
    ├── validate.go       # Reference validation
    ├── compile.go        # Standalone example compilation
    └── report.go         # Error reporting
```

### CLI Interface

```bash
# Validate all tutorial content
go run ./tools/tutorialcheck --all

# Validate specific files (for pre-commit)
go run ./tools/tutorialcheck hugo/content/go-tutorial/02-settings/structs.md

# Extract all code snippets (for debugging)
go run ./tools/tutorialcheck --extract-only

# Verbose output
go run ./tools/tutorialcheck --all --verbose
```

---

## Pre-commit & CI Integration

### Pre-commit Configuration

File: `.pre-commit-config.yaml`

```yaml
repos:
  - repo: local
    hooks:
      - id: tutorial-check
        name: Validate Go Tutorial
        entry: go run ./tools/tutorialcheck
        language: system
        files: '\.(md|go)$'
        pass_filenames: true

      - id: go-fmt
        name: Go Format
        entry: gofmt -w
        language: system
        files: '\.go$'

      - id: go-vet
        name: Go Vet
        entry: go vet ./...
        language: system
        pass_filenames: false

      - id: go-build
        name: Go Build
        entry: go build ./...
        language: system
        pass_filenames: false
```

### GitHub Actions Workflow

File: `.github/workflows/tutorial.yml`

```yaml
name: Tutorial Validation

on:
  push:
    paths:
      - 'hugo/content/go-tutorial/**'
      - '**.go'
      - 'tools/tutorialcheck/**'
  pull_request:
    paths:
      - 'hugo/content/go-tutorial/**'
      - '**.go'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Validate tutorial references
        run: go run ./tools/tutorialcheck --all
      - name: Build Hugo site
        run: cd hugo && hugo --minify
```

### Makefile Targets

```makefile
.PHONY: tutorial-check tutorial-extract tutorial-serve tutorial-new-chapter

tutorial-check:           ## Validate all tutorial references
	go run ./tools/tutorialcheck --all

tutorial-extract:         ## Re-extract all code snippets (debug)
	go run ./tools/tutorialcheck --extract-only

tutorial-serve:           ## Hugo server with live reload
	cd hugo && hugo server -D

tutorial-new-chapter:     ## Scaffold a new chapter (interactive)
	@read -p "Chapter number (e.g., 13): " num; \
	read -p "Chapter name (e.g., testing): " name; \
	mkdir -p hugo/content/go-tutorial/$$num-$$name; \
	echo "---\ntitle: \"$$name\"\nweight: $$num\n---\n" > hugo/content/go-tutorial/$$num-$$name/_index.md; \
	echo "Created hugo/content/go-tutorial/$$num-$$name/"
```

---

## Hugo Shortcodes

Hugo has two shortcode syntaxes:
- `{{</* shortcode */>}}` - Content is not processed as markdown
- `{{%/* shortcode */%}}` - Content is processed as markdown

Our shortcodes use `{{</* */>}}` for code (preserves formatting) and `{{%/* */%}}` for notice boxes (allows markdown inside).

All source file paths are relative to the repository root (e.g., `settings/settings.go`).

### code-ref

Extracts code between `TUTORIAL:id:start` and `TUTORIAL:id:end` markers.

**Usage:**
```markdown
{{</* code-ref id="settings/structs" */>}}
{{</* code-ref id="backend/gameapi" lang="go" title="The GameAPI Interface" */>}}
```

**Implementation:** `hugo/layouts/shortcodes/code-ref.html`

Reads source files at build time, finds markers, extracts code, renders as syntax-highlighted block.

### file-ref

Creates a link to source file with optional line numbers.

**Usage:**
```markdown
{{</* file-ref path="settings/settings.go" lines="47-55" */>}}
```

**Output:** Link to GitHub source with line highlighting.

### exercise

Collapsible exercise block with solution.

**Usage:**
```markdown
{{</* exercise title="Add a new setting" */>}}
Modify the Settings struct to add a `SoundEnabled` field...

{{</* solution */>}}
```go
type Settings struct {
    SoundEnabled bool `json:"sound_enabled"`
}
```
{{</* /solution */>}}
{{</* /exercise */>}}
```

---

## Content Standards

### Front Matter

Each tutorial page includes:

```yaml
---
title: "Structs and Types"
weight: 2
chapter: "02-settings"
concepts:
  - structs
  - type-definitions
  - json-tags
sources:
  - settings/settings.go
prerequisites:
  - 01-getting-started
---
```

### Callout Boxes

```markdown
{{% notice tip %}}
Go tip or best practice.
{{% /notice %}}

{{% notice warning %}}
Common mistake or gotcha.
{{% /notice %}}

{{% notice info %}}
Additional context or background.
{{% /notice %}}
```

### Code Conventions

- Use `code-ref` shortcode for Baboon source code
- Use fenced code blocks with `go compile` for standalone examples
- Use inline code (`backticks`) for function names, variables, commands

---

## Developer Workflow

### First-Time Setup

```bash
# Install pre-commit (if not already installed)
pip install pre-commit

# Install hooks
pre-commit install

# Verify setup
make tutorial-check
```

### Adding Tutorial Content

1. Identify code to reference in Baboon source
2. Add `TUTORIAL:chapter/section:start` and `:end` markers
3. Write tutorial markdown using `{{</* code-ref */>}}` shortcode
4. Add exercises with solutions
5. Run `make tutorial-check` to validate
6. Commit (pre-commit hook runs automatically)

### Modifying Baboon Source

When modifying source code that has tutorial markers:
1. Pre-commit hook validates tutorial still builds
2. If validation fails, update tutorial content to match
3. Both changes committed together

---

## Success Criteria

1. **Tutorial builds** - Hugo generates site without errors
2. **All references valid** - `tutorialcheck --all` passes
3. **Examples compile** - All `go compile` blocks build successfully
4. **Pre-commit enforced** - Hooks installed and running in CI
5. **Comprehensive coverage** - All listed chapters have content
6. **Exercises work** - Solutions compile and are correct

---

## Implementation Order

1. Create `tools/tutorialcheck/` validation tool
2. Create Hugo shortcodes (`code-ref`, `file-ref`, `exercise`)
3. Set up pre-commit configuration
4. Add GitHub Actions workflow
5. Add Makefile targets
6. Write Chapter 01: Getting Started
7. Add tutorial markers to `settings/settings.go`
8. Write Chapter 02: Settings
9. Continue with remaining chapters
10. Write appendix content

---

## Open Questions

None - all design decisions have been approved.

---

## Revision History

| Date | Change |
|------|--------|
| 2026-04-18 | Initial design approved |
