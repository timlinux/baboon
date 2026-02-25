# Baboon Package Architecture

This document provides an annotated overview of all packages in the Baboon application.

## Core Application Packages

### `main.go` (Entry Point)
The application entry point that wires together the backend and frontend components.

**Responsibilities:**
- Parse command-line flags (`-p`, `-port`, `-server`, `-client`)
- Start the appropriate mode (combined, server-only, or client-only)
- Handle graceful shutdown and PID file management

### `backend/` - Game Engine and REST Server

#### `api.go`
Defines the `GameAPI` interface that abstracts all game operations.

**Key Types:**
- `GameAPI` - Interface for game operations (start round, process keystrokes, get state)
- `GameState` - Current game state (words, input, timer status)
- `KeystrokeResult` - Result of processing a keystroke
- `SpaceResult` - Result of pressing space/advance key
- `Config` - Game configuration (punctuation mode, word count)

#### `engine.go`
Implements the game engine that tracks typing sessions.

**Responsibilities:**
- Word selection with adaptive weighting based on user's weak letters
- Keystroke processing with timing and accuracy tracking
- Session statistics calculation (WPM, accuracy, per-letter stats)
- Integration with typing theory metrics (finger, hand, row stats)

#### `server.go`
REST API server for multi-client support.

**Endpoints:**
- Session management (`POST/DELETE/GET /api/sessions`)
- Game operations (`/api/sessions/{id}/keystroke`, `/space`, `/round`)
- Statistics retrieval (`/api/sessions/{id}/stats/session`, `/historical`)
- Health check (`/api/health`)

### `frontend/` - Terminal User Interface

#### `model.go`
Bubble Tea model implementing the main event loop.

**States:**
- `StateTyping` - Active typing session
- `StateResults` - Viewing round results
- `StateOptions` - Configuring settings

**Responsibilities:**
- Handle keyboard input for each state
- Manage local timing for accurate WPM calculation
- Coordinate with backend via GameAPI interface

#### `views.go`
All rendering functions for the TUI.

**Functions:**
- `RenderTypingScreenAnimated()` - Main typing interface with carousel
- `RenderResultsScreen()` - Statistics display after round completion
- `RenderOptionsScreen()` - Settings configuration UI

#### `styles.go`
Lipgloss style definitions and colour constants.

**Includes:**
- Colour codes for correct/incorrect/untyped letters
- Gradient colours for progress bars
- Styles for labels, values, titles, and help text

#### `animations.go`
Spring-based animation system using harmonica library.

**Components:**
- `CarouselAnimator` - Word transition animations (previous/current/next)
- `Animator` - Results screen row animations with stagger

#### `client.go`
REST API client implementing the `GameAPI` interface.

**Responsibilities:**
- HTTP communication with backend server
- Session creation and cleanup
- Request/response serialization

### `font/` - Block Letter Rendering

#### `font.go`
Custom block font using Unicode characters.

**Features:**
- 6-line tall letters using `█`, `▀`, `▄` characters
- Supports a-z lowercase and punctuation (, . ; : ! ?)
- `RenderWord()` function for multi-letter rendering

### `words/` - Dictionary

#### `words.go`
Collection of common English words with British spellings.

**Features:**
- 250+ words for varied practice
- British English spellings (colour, behaviour, centre)
- Words filtered for lowercase letters only

### `stats/` - Statistics and Persistence

#### `stats.go`
Statistics types and persistence logic.

**Key Types:**
- `Stats` - Session statistics (WPM, accuracy, letter stats)
- `HistoricalStats` - Cumulative statistics across sessions
- Recording functions for letters, fingers, hands, rows, errors

**Persistence:**
- Stats saved to `~/.config/baboon/stats.json`
- Validation and corruption detection on load

#### `keyboard.go`
QWERTY keyboard layout mappings.

**Mappings:**
- Letter to finger assignment (8 fingers)
- Letter to hand (left/right)
- Letter to row (top/home/bottom)
- Same-finger bigram detection

### `settings/` - User Preferences

#### `settings.go`
User configuration persistence.

**Settings:**
- `AdvanceKey` - Which key advances to next word (Space, Enter, Either)

**Persistence:**
- Settings saved to `~/.config/baboon/settings.json`
- Loaded on startup with sensible defaults

## Web Frontend Packages (`web/`)

### `src/App.js`
Main React component with state management.

**States:**
- Welcome screen (connection check, game options)
- Typing screen (active session)
- Results screen (post-round statistics)

### `src/api.js`
REST API client for web frontend.

**Functions:**
- Session management
- Game operations (keystrokes, round control)
- Statistics retrieval

### `src/theme.js`
Chakra UI custom theme with Kartoza brand colours.

**Includes:**
- Dark mode configuration
- Custom button variants ("glow", "chunky")
- Kartoza orange (#D4922A) and blue (#4A90A4) colour palettes

### `src/components/`

#### `WelcomeScreen.js`
Landing page with animated logo and game options.

#### `TypingScreen.js`
Main typing interface with physics-based block letters.

#### `ResultsScreen.js`
Statistics display with animated stat cards.

## External Dependencies

### Go Dependencies
- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/lipgloss` - TUI styling
- `github.com/charmbracelet/harmonica` - Spring physics animations

### Web Dependencies
- `react` (18.x) - UI framework
- `@chakra-ui/react` (2.x) - Component library
- `framer-motion` - Physics-based animations

## Package Dependency Graph

```
main.go
├── backend/
│   ├── api.go (types)
│   ├── engine.go → stats/, words/
│   └── server.go → engine.go
└── frontend/
    ├── model.go → backend/api, settings/
    ├── views.go → backend/, stats/, settings/, font/
    ├── styles.go
    ├── animations.go
    └── client.go → backend/api

stats/
├── stats.go
└── keyboard.go

settings/
└── settings.go

font/
└── font.go

words/
└── words.go
```
