# Changelog

All notable changes to Baboon are documented here.

## [v1.2.0] - 2024

### Added

- **Extended word carousel**: Now shows next 3 upcoming words below the current word
  - Words displayed with decreasing opacity for visual hierarchy
  - First upcoming word has decorative arrows (▼)
  - All words centered horizontally
  - Provides better look-ahead for typing preparation

### Changed

- Updated `GameState` API to include `NextWords` slice (array of up to 3 words)
- Backwards compatible: falls back to `NextWord` if `NextWords` is empty

---

## [v1.1.0] - 2024

### Added

- **Beautiful carousel animation for word transitions**
  - **Console (TUI)**: Smooth harmonica spring-based animations
    - Previous word fades in with animated greyscale opacity
    - Current word slides up with spring physics
    - Next word fades in from below with staggered timing
  - **Web**: Framer Motion spring animations
    - Previous/next words at 50% scale with blur
    - Decorative glow effect behind current word

### Fixed

- **Accuracy statistics exceeding 100%**: Fixed bug where backspacing and retyping a character counted "Correct" multiple times
  - Now tracks which character positions have been recorded
  - Accuracy stats only recorded on first correct keystroke per position
  - Timing stats still recorded for all keystrokes

---

## [v1.0.0] - 2024

### Added

- **First stable release**
- Beautiful README with screenshots and badges
- **GitHub Actions CI/CD workflows**:
  - Test workflow: runs on push/PR
  - Build workflow: cross-platform verification
  - Release workflow: automated builds on tag push
- **Pre-built binaries**:
  - Linux AMD64 and ARM64
  - macOS Intel and Apple Silicon
  - Windows AMD64
  - DEB package for Debian/Ubuntu
  - RPM package for Fedora/RHEL
  - Flatpak package
- macOS unsigned binary instructions
- Nix flake integration for system configurations

---

## [v0.9.1] - 2024

### Changed

- **Kartoza brand colour scheme** applied to web frontend
  - Primary: Kartoza Orange (#D4922A)
  - Secondary: Kartoza Blue (#4A90A4)
  - Updated all UI components with brand colours

### Added

- Kartoza wallpaper asset for reference

---

## [v0.9.0] - 2024

### Added

- **React web frontend** with 100% feature parity to TUI
  - Built with React 18, Chakra UI 2.x, Framer Motion
  - Physics-based animations using spring dynamics
  - Large chunky block letters with bounce effects
  - Dark theme with custom colours
- **Web components**:
  - WelcomeScreen: Animated logo, connection status, game options
  - TypingScreen: Block letters with physics, live WPM bar
  - ResultsScreen: Full statistics display
- Letter statistics grid with colour-coded indicators
- Finger accuracy and hand balance displays
- Makefile targets for web development

---

## [v0.8.0] - 2024

### Added

- **Management scripts** for backend server lifecycle
  - `start-backend.sh`: Start backend in background
  - `stop-backend.sh`: Graceful or forced shutdown
  - `status-backend.sh`: Health check and monitoring
  - `launch-frontend.sh`: Launch frontend client
- **Three running modes**:
  - Combined mode (default): Backend + frontend together
  - Server-only mode (`-server`)
  - Client-only mode (`-client`)
- PID file tracking for server management
- Graceful shutdown handling

---

## [v0.7.0] - 2024

### Added

- **Frontend timing implementation** to eliminate network latency
  - Seek times measured locally on frontend
  - Live WPM calculated on frontend
  - Round duration submitted via timing endpoint
- New API methods for timing data
- New REST endpoint: `POST /api/sessions/{id}/timing`

### Changed

- Backend no longer calls `time.Now()` for timing-critical operations

---

## [v0.6.0] - 2024

### Added

- **Multi-client session management**
  - Each client creates unique session on startup
  - 32-character hex session IDs
  - All operations scoped to sessions
  - Automatic session cleanup
- New session management endpoints
- Health endpoint reports active session count
- Thread-safe session storage

---

## [v0.5.0] - 2024

### Added

- **RESTful API** for frontend-backend communication
  - Backend server in `backend/server.go`
  - Frontend client in `frontend/client.go`
  - All game operations via HTTP
- REST endpoints for game operations
- Configurable port via `-port` flag
- Thread-safe server implementation

---

## [v0.4.0] - 2024

### Changed

- **Major architecture refactoring**
  - Clean backend/frontend separation
  - Backend package with `GameAPI` interface
  - Frontend package for TUI
  - Clear API boundaries
- Improved code organisation

---

## [v0.3.0] - 2024

### Added

- **Typing theory statistics**:
  - Finger-specific accuracy and speed (8 fingers)
  - Keyboard row tracking (top, home, bottom)
  - Hand balance and alternation rate
  - Same-finger bigram (SFB) detection
  - Rhythm consistency (standard deviation)
  - Error substitution patterns
- Results screen displays new metrics
- Increased animated rows from 14 to 25

---

## [v0.2.0] - 2024

### Added

- **Per-letter accuracy tracking** (a-z)
- Letter statistics persistence across sessions
- Results screen 26-letter heatmaps:
  - Accuracy row
  - Frequency row
  - Seek time row
- Fixed 30 words / 150 characters per round
- Smooth font with Unicode half-blocks
- **Adaptive word selection** based on accuracy
- Letter seek time tracking
- Bigram seek time tracking
- Results screen spring animations
- **Punctuation mode** (`-p` flag)

### Changed

- Improved seek time calculation:
  - Only records for correct keystrokes
  - Records against expected letter
  - Excludes first letter of each word
- Redesigned letter statistics display

---

## [v0.1.0] - Initial Release

### Added

- Basic typing practice with 30-word rounds
- Block letter word display
- Real-time colour feedback
- Live WPM bar during typing
- Results screen with comparison bars
- WPM, time, and accuracy tracking
- Historical best comparison with stars
- British English word dictionary
- Cross-platform support via Nix
- Statistics persistence with corruption detection

---

## Roadmap

### Planned Features

- [ ] Multiple word lists (programming, languages)
- [ ] Timed mode options
- [ ] Leaderboards
- [ ] Custom themes
- [ ] Sound effects
- [ ] Accessibility improvements
- [ ] Mobile-friendly web interface
- [ ] Practice specific letter combinations

### Under Consideration

- Multi-player racing mode
- Integration with typing test websites
- Keyboard heatmap visualisation
- Custom word list import
