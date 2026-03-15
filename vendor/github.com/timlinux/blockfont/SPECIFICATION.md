# blockfont Specification

## Overview

blockfont is a Go library for rendering Unicode block letters in terminal applications. It provides a comprehensive solution for displaying large, stylized text with support for vim-style editing, animations, and integration with the charmbracelet ecosystem.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         blockfont                               │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │  Widget  │──│  Buffer  │  │ Animator │  │  Style   │        │
│  │(tea.Model)│  │  (vim)  │  │ (spring) │  │ (lipgloss)│       │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │
│       │             │             │             │               │
│       └─────────────┴─────────────┴─────────────┘               │
│                           │                                     │
│  ┌──────────┐  ┌──────────┴──────────┐  ┌──────────┐           │
│  │   Font   │──│       Layout        │──│   ANSI   │           │
│  │(BlockLetters)│  (align/wrap)     │  │(utilities)│           │
│  └──────────┘  └────────────────────┘  └──────────┘           │
└─────────────────────────────────────────────────────────────────┘
```

## User Stories

### US-001: Simple Text Rendering
**As a** developer
**I want to** render text as block letters
**So that** I can display prominent text in my terminal application

**Acceptance Criteria:**
- Can call `RenderText("hello")` and get a multi-line string
- Output contains Unicode block characters (█, ◢, ◣, ◤, ◥)
- All lowercase and uppercase letters are supported
- Numbers and common punctuation are supported

### US-002: Widget Integration
**As a** developer using bubbletea
**I want to** use a pre-built Widget component
**So that** I can easily integrate block text into my TUI

**Acceptance Criteria:**
- Widget implements tea.Model interface
- Can set text with `SetText()`
- Can render with `View()`
- Supports width/alignment configuration

### US-003: Vim-Style Editing
**As a** developer building an editor
**I want to** enable vim-style text editing
**So that** users can efficiently modify block text

**Acceptance Criteria:**
- Normal mode for navigation (h, j, k, l)
- Insert mode for typing (i, a, o)
- Delete operations (x, dd, D)
- Cursor is visually indicated

### US-004: Character Highlighting
**As a** developer building a typing game
**I want to** highlight individual characters
**So that** I can show correct/incorrect input

**Acceptance Criteria:**
- Can set highlights per character
- Supports correct (green), incorrect (red), pending (gray)
- Cursor highlight inverts character
- Custom colors can be set per character

### US-005: Animation Support
**As a** developer
**I want to** animate text transitions
**So that** my application feels polished

**Acceptance Criteria:**
- Slide up/down transitions
- Fade in/out transitions
- Scale transitions
- Spring-based physics for smooth motion

## Functional Requirements

### FR-001: Font Data
- Characters defined as 6-line arrays of strings
- Uses Unicode block elements: █ (full block), ◢◣◤◥ (triangles)
- Variable width characters
- Fallback to space for undefined characters

### FR-002: Core Rendering
- `RenderWord()` returns 2D array of character lines
- `RenderText()` returns joined string with newlines
- `RenderWithCursor()` renders with cursor position, highlights, insert mode, and word wrapping
- `RenderPlainText()` renders without cursor or highlights
- `GetLetterWidth()` returns character width in runes
- `GetTotalWidth()` returns word width including spacing
- `CalculateTotalWidth()` calculates width including cursor position

### FR-003: Buffer Operations
- Multi-line text support
- Cursor position tracking (X, Y)
- Insert at cursor
- Delete from cursor (n characters)
- Delete line
- Delete to end of line
- Replace character

### FR-004: Layout Functions
- Center lines within width
- Left-justify with margin
- Right-justify within width
- Word wrapping on space boundaries
- Width calculation excluding ANSI codes

### FR-005: Styling
- ANSI color code support
- Lipgloss style integration
- Theme support with predefined themes
- Per-character coloring
- Highlight states (correct, incorrect, cursor, etc.)

### FR-006: Animation
- Spring physics using harmonica library
- Position animation (slide)
- Opacity animation (fade)
- Scale animation (grow/shrink)
- Configurable timing

## Technical Requirements

### TR-001: Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/harmonica` - Spring animations

### TR-002: Performance
- No allocations in hot paths where possible
- Efficient string building with strings.Builder
- Lazy rendering (only when View() called)

### TR-003: Testing
- Unit tests for all public functions
- Table-driven tests for multiple inputs
- Test coverage for edge cases (empty strings, unicode)

### TR-004: Documentation
- Godoc comments on all public types and functions
- README with examples
- SPECIFICATION.md (this document)
- PACKAGES.md describing package structure

## API Specification

### Constants

```go
const LetterHeight = 6      // Height of all block letters
const LetterSpacing = 1     // Space between letters
const AnimationInterval = 50 * time.Millisecond
```

### Types

```go
type Alignment int          // AlignLeft, AlignCenter, AlignRight
type Mode int               // ModeNormal, ModeInsert, ModeVisual, etc.
type CharHighlight int      // HighlightNone, HighlightCorrect, etc.
type CursorStyle int        // CursorBlock, CursorLine, CursorUnderline
type TransitionType int     // TransitionSlideUp, TransitionFadeIn, etc.
```

### Primary Interfaces

```go
// Widget provides a high-level tea.Model for block text
type Widget struct { ... }
func NewWidget(opts WidgetOptions) *Widget
func (w *Widget) SetText(text string)
func (w *Widget) Text() string
func (w *Widget) View() string
func (w *Widget) Update(msg tea.Msg) (*Widget, tea.Cmd)

// Buffer provides vim-style text editing
type Buffer struct { ... }
func NewBuffer(text string) *Buffer
func (b *Buffer) Text() string
func (b *Buffer) Insert(text string)
func (b *Buffer) Delete(n int) string

// Animator provides spring-based animations
type Animator struct { ... }
func NewAnimator() *Animator
func (a *Animator) TriggerTransition(t TransitionType)
func (a *Animator) Update() bool
```

## Version History

### v0.1.0 (Initial Release)
- Core font rendering
- Widget with bubbletea integration
- Vim buffer editing
- Layout utilities
- Styling with themes
- Spring animations
- Character highlighting

---

Made with :heart: by [Kartoza](https://kartoza.com) | [Donate!](https://github.com/sponsors/timlinux) | [GitHub](https://github.com/timlinux/blockfont)
