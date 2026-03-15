# blockfont Package Structure

## Overview

blockfont is organized as a single Go package providing Unicode block letter rendering for terminal applications.

## Package: `blockfont`

**Import Path:** `github.com/timlinux/blockfont`

### Files

| File | Description |
|------|-------------|
| `font.go` | Block letter definitions and core rendering functions |
| `render.go` | Full-featured rendering with cursor, highlights, and word wrapping |
| `widget.go` | High-level tea.Model widget for bubbletea integration |
| `buffer.go` | Vim-style text buffer with editing operations |
| `layout.go` | Text alignment and layout utilities |
| `style.go` | Styling, themes, and character highlighting |
| `animation.go` | Spring-based animation system |
| `ansi.go` | ANSI escape code utilities |

### Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework for widget integration |
| `github.com/charmbracelet/lipgloss` | Terminal styling and colors |
| `github.com/charmbracelet/harmonica` | Spring physics for animations |

## Component Breakdown

### font.go
The core font rendering system.

**Exports:**
- `BlockLetters` - Map of rune to 6-line string arrays
- `LetterHeight` - Constant (6)
- `LetterSpacing` - Constant (1)
- `RenderWord()` - Render word to 2D array
- `RenderText()` - Render word to string with newlines
- `GetLetterWidth()` - Get character width
- `GetTotalWidth()` - Get word width

### widget.go
High-level widget for bubbletea applications.

**Exports:**
- `Widget` - Main widget struct (implements tea.Model)
- `WidgetOptions` - Configuration options
- `DefaultWidgetOptions()` - Sensible defaults
- `NewWidget()` - Constructor
- `AnimationTickMsg` - Message type for animation updates

### buffer.go
Vim-style text editing buffer.

**Exports:**
- `Buffer` - Text buffer with cursor
- `Mode` - Editing mode enum (Normal, Insert, Visual, etc.)
- `NewBuffer()` - Constructor
- Movement methods: `MoveLeft()`, `MoveRight()`, `MoveUp()`, `MoveDown()`
- Edit methods: `Insert()`, `Delete()`, `DeleteLine()`, `ReplaceChar()`

### layout.go
Text layout and alignment utilities.

**Exports:**
- `Alignment` - Enum (Left, Center, Right)
- `CenterLines()` - Center text lines
- `LeftJustify()` - Left-justify with margin
- `RightJustify()` - Right-justify within width
- `AlignLines()` - Apply alignment
- `WrapOnWordBoundaries()` - Word wrap text
- `PadToWidth()` - Pad string to width
- `MaxLineWidth()` - Find max line width
- `JoinBlockLines()` - Join block letter arrays

### style.go
Styling and theming system.

**Exports:**
- `CharHighlight` - Highlight type enum
- `CursorStyle` - Cursor appearance enum
- `Theme` - Color theme struct
- `StyleConfig` - Configuration struct
- `DefaultTheme` - Default colors
- `KartozaTheme` - Kartoza-branded colors
- `GradientColors` - Red-to-green gradient array
- `GetGradientColor()` - Get color from gradient
- `GetWPMColor()` - WPM-based color
- `GetProgressColor()` - Progress-based color

### animation.go
Spring-based animation system.

**Exports:**
- `AnimationInterval` - Default tick interval
- `TransitionType` - Transition enum
- `Animator` - Main animator struct
- `NewAnimator()` - Constructor
- `WordCarouselAnimator` - Specialized animator for word carousels
- `GetAnimationInterval()` - Get recommended tick interval

### ansi.go
ANSI escape code utilities.

**Exports:**
- ANSI constants: `ANSIReset`, `ANSIRed`, `ANSIGreen`, etc.
- `RemoveANSI()` - Strip ANSI codes
- `VisibleStringWidth()` - Width excluding ANSI
- `InsertAt()` - Insert text at position
- `WrapWithColor()` - Wrap text with color codes
- `InvertLine()` - Apply inverse video

## Usage Flow

```
Application
    │
    ▼
┌──────────┐
│ NewWidget │ ◄── Configure options
└────┬─────┘
     │
     ▼
┌──────────┐
│ SetText  │ ◄── Set content
└────┬─────┘
     │
     ├──────────────────────────────────────┐
     ▼                                      ▼
┌──────────┐                          ┌──────────┐
│ Update   │ ◄── Handle keys         │  View    │ ◄── Render output
└────┬─────┘                          └────┬─────┘
     │                                     │
     ▼                                     ▼
┌──────────┐                          ┌──────────┐
│  Buffer  │ ◄── Edit text            │  Render  │ ◄── Generate lines
└──────────┘                          └────┬─────┘
                                           │
                                           ▼
                                      ┌──────────┐
                                      │ RenderWord │ ◄── Core font
                                      └──────────┘
```

---

Made with :heart: by [Kartoza](https://kartoza.com) | [Donate!](https://github.com/sponsors/timlinux) | [GitHub](https://github.com/timlinux/blockfont)
