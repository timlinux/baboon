# Baboon - Typing Practice Application

## Overview

Baboon is a cross-platform typing practice application built with Go. It helps users improve their typing speed and accuracy by presenting common English words in large block letter format. The application features two frontends:

1. **Terminal UI (TUI)**: Built with Bubble Tea and Lipgloss, displaying words using Unicode block characters (█) that change colour in real-time as the user types.

2. **Web UI**: Built with React and Chakra UI, featuring physics-based animations with Framer Motion, large chunky block letters, and a beautiful dark theme.

Both frontends communicate with the same Go backend via REST API, ensuring 100% feature parity.

## User Stories

### US-001: Practice Typing Common Words
**As a** user wanting to improve my typing skills
**I want to** practice typing common English words
**So that** I can increase my typing speed and accuracy

### US-002: Visual Feedback During Typing
**As a** user practicing typing
**I want to** see immediate visual feedback on my keystrokes
**So that** I can identify and correct errors quickly

### US-003: Track My Progress
**As a** regular user
**I want to** see my typing statistics after each round
**So that** I can measure my improvement over time

### US-004: Compare to Personal Best
**As a** competitive user
**I want to** compare my current performance to my historical best and average
**So that** I can see how I'm improving

### US-005: Monitor Speed in Real-Time
**As a** user during a typing session
**I want to** see my current WPM as I type
**So that** I can adjust my pace accordingly

### US-006: Access via Web Browser
**As a** user who prefers a web-based interface
**I want to** practice typing through a beautiful web application
**So that** I can use any device with a modern browser

## Functional Requirements

### FR-001: Word Display
- The application SHALL display words from a dictionary of common English words (British English spelling)
- Each word SHALL be rendered in large block characters using Unicode block elements (█)
- Words SHALL be displayed centered horizontally and vertically on the terminal screen
- The word display SHALL show progress indicator: "Word X/30"
- Letters SHALL change colour in-place as the user types (no separate input display line)
- All words SHALL be lowercase only (the font only supports a-z)
- Words SHALL be displayed in a carousel layout:
  - The previous word SHALL be displayed ABOVE the current word in dimmed text
  - The next 3 upcoming words SHALL be displayed BELOW the current word in dimmed text
  - When advancing to the next word, smooth carousel animation SHALL scroll words upward
- Console: Previous word uses greyscale colour (240), next words use decreasing greyscale (from 245), with decorative markers on first upcoming word
- Web: Previous/next words displayed at 50% scale with blur effects

### FR-002: Block Letter Font
- Each letter SHALL be 6 lines tall
- Letters SHALL be constructed using Unicode block elements for smooth edges:
  - █ (full block) for solid areas
  - ▀ (upper half block) for rounded tops
  - ▄ (lower half block) for rounded bottoms
- The font SHALL support lowercase letters a-z and punctuation: , . ; : ! ?
- Unknown characters SHALL render as spaces
- Letters SHALL have 1 character spacing between them

### FR-003: Typing Input and Colour Feedback
- The application SHALL accept keyboard input character by character
- Each character typed SHALL immediately change the corresponding letter's colour:
  - **Green (colour 10)**: Character matches the expected character at that position
  - **Red (colour 9)**: Character does not match the expected character
  - **Gray (colour 8)**: Characters not yet typed
- The backspace key SHALL remove the last typed character (undoing its colour)
- The space key SHALL only advance to the next word when ALL letters have been typed
- If space is pressed before the word is complete, it SHALL be treated as an incorrect character (red)
- Extra characters beyond word length SHALL count as incorrect (red)

### FR-004: Round Structure
- Each round SHALL consist of exactly 30 words totalling exactly 150 characters
- Words SHALL be randomly selected using stratified selection to meet both constraints
- Word selection algorithm SHALL:
  - Calculate ideal word length based on remaining characters and words
  - Allow variance of ±2 characters from ideal to maintain variety
  - Ensure feasibility by checking remaining capacity
  - Retry up to 100 times if constraints cannot be met
- After completing all 30 words, the application SHALL display the results screen
- The user SHALL be able to start a new round by pressing Enter on results screen

### FR-015: Adaptive Word Selection
- Word selection SHALL be weighted based on two factors:
  1. **Frequency balancing**: Favour words with underrepresented letters
  2. **Accuracy practice**: Favour words with letters the user frequently mistypes
- Each word SHALL be scored using a combined algorithm:
  - Frequency score = 1 - (letter_presented / max_letter_presented)
  - Accuracy score = 1 - (letter_correct / letter_presented)
  - Letter score = (frequency_score + accuracy_score) / 2
  - Word score = average letter score across all letters in the word
- Words with higher scores SHALL have higher selection probability
- This adaptive selection helps users practice their weakest letters
- Frequency balancing aims to achieve spread within 10% from highest to lowest

### FR-005: Timer Behaviour
- The timer SHALL NOT start when the application launches
- The timer SHALL start when the user types the first CORRECT character of the first word
- If the first character typed is incorrect, the timer SHALL NOT start
- The timer SHALL stop when the 30th word is completed (space pressed)

### FR-006: Live WPM Bar (During Typing)
- The application SHALL display a gradient WPM bar at the bottom of the screen during typing
- The bar SHALL be 50 characters wide
- The bar SHALL update every 100ms to show current typing speed
- The bar SHALL use a gradient colour scheme from red (slow) through yellow to green (fast):
  - Colours: 196, 202, 208, 214, 220, 226, 190, 154, 118, 82, 46, 47
- The bar SHALL scale from 0 to 120 WPM maximum
- The bar SHALL display numeric WPM value with colour coding:
  - Red (colour 196): Below 40 WPM
  - Yellow (colour 226): 40-60 WPM
  - Green (colour 46): Above 60 WPM
- The bar SHALL show scale markers: "0", "60", "120"
- Empty portion of bar SHALL use character ░ in colour 236

### FR-007: Results Screen Layout
- The results screen SHALL display "Round Complete!" title in cyan (colour 14), bold
- Statistics SHALL be displayed in a grid layout with three columns:
  - Label column: 18 characters wide, right-aligned, gray (colour 7)
  - Value column: 8 characters wide, right-aligned, white (colour 15)
  - Bar column: 30 characters wide gradient bar + 2 character star column

### FR-008: Results Statistics Display
- WPM section:
  - "WPM this run:" with current session WPM and bar
  - "WPM best:" with historical best WPM and bar
  - "WPM average:" with calculated average WPM and bar
- Time section (blank line before):
  - "Time this run:" with session time in seconds (e.g., "147.2s") and bar
  - "Time best:" with historical best (fastest) time and bar
  - "Time average:" with calculated average time and bar
- Accuracy section (blank line before):
  - "Accuracy this run:" with percentage (e.g., "95.5%") and bar
  - "Accuracy best:" with historical best accuracy and bar
  - "Accuracy average:" with calculated average accuracy and bar
- Sessions section (blank line before):
  - "Total sessions:" label in cyan (colour 6) with count
- Legend (blank line before):
  - "* = New personal best!" in yellow (colour 226), bold

### FR-009: Results Bar Rendering
- WPM bars: Scale 0-120, higher is better (more fill = better)
- Time bars: Scale 0-180 seconds, INVERTED (lower time = more fill = better)
- Accuracy bars: Scale 0-100%, higher is better (more fill = better)
- New personal best SHALL show " *" after the bar
- Non-best bars SHALL show "  " (two spaces) to maintain alignment
- All bars SHALL use same gradient colours as live WPM bar

### FR-010: Historical Statistics Persistence
- Historical data SHALL be stored in `~/.config/baboon/stats.json`
- The config directory SHALL be created automatically if it doesn't exist
- The application SHALL track:
  - `best_wpm`: Highest WPM achieved (float64)
  - `best_accuracy`: Highest accuracy percentage achieved (float64)
  - `best_time`: Fastest (lowest) completion time in seconds (float64)
  - `total_wpm`: Sum of all session WPMs for averaging (float64)
  - `total_accuracy`: Sum of all session accuracies for averaging (float64)
  - `total_time`: Sum of all session times for averaging (float64)
  - `total_sessions`: Count of completed sessions (int)
  - `last_session_date`: Timestamp of last session (RFC3339)
  - `letter_accuracy`: Per-letter accuracy tracking (map of letter to stats)
  - `letter_seek_time`: Per-letter seek time tracking (map of letter to timing stats)
  - `bigram_seek_time`: Per-bigram (letter pair) seek time tracking (map of bigram to timing stats)

### FR-013: Per-Letter Accuracy Tracking
- When a round starts, all letters in all 30 words SHALL be recorded as "presented"
- When a user types a correct letter, that letter SHALL be recorded as "correct"
- Letter statistics SHALL be tracked per individual letter (a-z)
- For each letter, the application SHALL track:
  - `presented`: Number of times this letter was presented to the user
  - `correct`: Number of times the user typed this letter correctly
- Letter accuracy data SHALL persist across sessions (cumulative)
- Letter accuracy SHALL be calculated as: (correct / presented) × 100

### FR-016: Per-Letter Seek Time Tracking
- The application SHALL track the time between keystrokes (seek time)
- Seek time SHALL only be recorded for CORRECT keystrokes
- Seek time SHALL be recorded against the EXPECTED letter (not the typed character)
- The FIRST letter of each word SHALL be excluded from seek time tracking (includes word-reading time)
- Seek times > 5000ms SHALL be filtered out (assumed user pauses)
- For each letter, the application SHALL track:
  - `total_time_ms`: Total seek time in milliseconds
  - `count`: Number of measurements
- Average seek time = total_time_ms / count
- Seek time data SHALL persist across sessions (cumulative)

### FR-017: Bigram (Letter Pair) Seek Time Tracking
- The application SHALL track seek time for letter pairs (bigrams)
- A bigram is formed from the previous correctly typed letter + current correctly typed letter
- Bigrams SHALL only be recorded for consecutive correct keystrokes
- Bigrams SHALL reset at word boundaries (first letter of new word has no preceding letter)
- For each bigram (e.g., "th", "he", "in"), the application SHALL track:
  - `total_time_ms`: Total seek time in milliseconds
  - `count`: Number of measurements
- Bigram data SHALL persist across sessions (cumulative)
- Common slow bigrams indicate letter combinations the user struggles with

### FR-014: Letter Statistics Display
- The results screen SHALL display a letter statistics matrix:
  1. **Header row**: 26 uppercase letters (A-Z) as column labels, white bold text
  2. **Accuracy row**: Filled circles (●) coloured by typing accuracy
  3. **Frequency row**: Filled circles (●) coloured by presentation count
  4. **Seek time row**: Filled circles (●) coloured by average typing speed
- Each circle SHALL be coloured on a red-to-green gradient
- Letters in header row are spaced to align with circles below
- Seek time is measured as milliseconds between keystrokes
- Seek times > 5 seconds are filtered out (assumed pauses)
- Gradient colours (accuracy/speed percentage → colour code):
  - 95-100%: 46 (bright green)
  - 90-94%: 82
  - 85-89%: 118
  - 80-84%: 154
  - 75-79%: 190
  - 70-74%: 226 (yellow)
  - 65-69%: 220
  - 60-64%: 214
  - 50-59%: 208
  - 40-49%: 202
  - 0-39%: 196 (red)

### FR-018: Results Screen Animation
- Results screen elements SHALL animate in sequentially using spring physics
- The harmonica library SHALL be used for smooth spring-based animations
- Each stat row SHALL slide in from the right with staggered timing
- Animation interval SHALL be 50ms per frame
- Stagger delay SHALL be 3 frames between each row starting
- Spring parameters: 60 FPS, frequency 6.0, damping 0.5
- Total of 25 animated rows (core stats + typing theory stats + letter matrix)

### FR-020: Finger-Specific Statistics
- The application SHALL track per-finger typing accuracy and speed
- Standard touch typing finger assignments SHALL be used:
  - Left pinky (0): q, a, z
  - Left ring (1): w, s, x
  - Left middle (2): e, d, c
  - Left index (3): r, f, v, t, g, b
  - Right index (6): y, h, n, u, j, m
  - Right middle (7): i, k
  - Right ring (8): o, l
  - Right pinky (9): p
- For each finger, the application SHALL track:
  - `presented`: Times a key for this finger was needed
  - `correct`: Times the correct key was pressed
  - `total_time_ms`: Total seek time for keys typed with this finger
  - `count`: Number of timed keypresses
- Results screen SHALL display finger accuracy row with colour-coded indicators
- Finger labels: LP, LR, LM, LI (left hand), RI, RM, RR, RP (right hand)

### FR-021: Keyboard Row Statistics
- The application SHALL track per-row typing accuracy and speed
- Row assignments:
  - Top row (0): q, w, e, r, t, y, u, i, o, p
  - Home row (1): a, s, d, f, g, h, j, k, l
  - Bottom row (2): z, x, c, v, b, n, m
- For each row, the application SHALL track:
  - `presented`: Times a key on this row was needed
  - `correct`: Times the correct key was pressed
  - `total_time_ms`: Total seek time for keys on this row
  - `count`: Number of timed keypresses
- Results screen SHALL display row accuracy with labels: Top, Home, Bot

### FR-022: Hand Balance and Alternation Tracking
- The application SHALL track hand usage balance (left vs right)
- Hand assignments: Left (q-t, a-g, z-b), Right (y-p, h-l, n-m)
- The application SHALL track hand alternations vs same-hand runs:
  - `hand_alternations`: Count of transitions between hands
  - `same_hand_runs`: Count of consecutive same-hand keypresses
- Alternation rate = hand_alternations / (hand_alternations + same_hand_runs) × 100
- Higher alternation rate indicates better typing flow
- Results screen SHALL display hand balance (L:X% R:Y%) and alternation rate

### FR-023: Same-Finger Bigram (SFB) Tracking
- The application SHALL detect and track same-finger bigrams
- An SFB occurs when consecutive letters use the same finger
- For each SFB occurrence, the application SHALL track:
  - Count of SFBs encountered
  - Total seek time for SFBs
  - Average seek time = total_time / count
- SFBs are inherently slower than alternating-finger bigrams
- Results screen SHALL display SFB count and average time per session

### FR-024: Rhythm Consistency (Variance) Tracking
- The application SHALL track typing rhythm consistency
- Rhythm is measured as the standard deviation of seek times
- For rhythm calculation, the application SHALL track:
  - All seek times during the session
  - Sum of seek times
  - Sum of squared seek times (for variance calculation)
- Variance = (sum_of_squares / count) - (mean²)
- Standard deviation = √variance
- Lower standard deviation indicates more consistent rhythm
- Results screen SHALL display session StdDev and historical average

### FR-025: Error Substitution Pattern Tracking
- The application SHALL track which letters are commonly confused
- When a letter is mistyped, the application SHALL record:
  - The expected letter
  - The typed letter
  - Increment the count for this (expected → typed) pair
- Error substitution data SHALL persist across sessions
- Results screen SHALL display top 5 most common error patterns
- Format: "a→s(12)" means 'a' was typed as 's' 12 times

### FR-019: Punctuation Mode
- The application SHALL support a `-p` command line flag for punctuation mode
- When enabled, words SHALL be separated by random punctuation followed by space
- Supported punctuation characters: , . ; : ! ?
- Punctuation SHALL be appended to each word except the last word in the round
- The user SHALL type the punctuation character before pressing space to advance
- Letter accuracy tracking SHALL only count letters (a-z), not punctuation
- Letter seek time tracking SHALL only measure letters (a-z), not punctuation
- Punctuation mode persists for subsequent rounds until the application exits

### FR-011: Statistics Validation
- On load, the application SHALL validate historical statistics for corruption
- If totals are 0 but bests exist, data SHALL be reset using best values as estimates
- If average WPM is less than half of best WPM, data SHALL be reset
- Reset formula: total = best × total_sessions

### FR-012: Navigation
- ESC or Ctrl+C SHALL exit the application at any time
- SPACE SHALL advance to the next word during typing (when input length > 0)
- ENTER SHALL start a new round when viewing results screen
- The application SHALL use alternate screen buffer (fullscreen mode)

## Technical Requirements

### TR-001: Cross-Platform Compatibility
- The application SHALL run on Linux, macOS, and Windows
- The application SHALL be buildable using Nix flakes for reproducible builds
- The terminal SHALL support 256-colour mode for proper gradient display

### TR-002: Terminal Interface
- The application SHALL use the Bubble Tea framework (github.com/charmbracelet/bubbletea)
- The application SHALL use lipgloss for styling (github.com/charmbracelet/lipgloss)
- The application SHALL use custom block font rendering (no external font libraries)
- The application SHALL use tea.WithAltScreen() for fullscreen mode
- The application SHALL handle tea.WindowSizeMsg for responsive centering

### TR-003: Update Loop
- The application SHALL use tea.Tick with 100ms interval for WPM bar updates
- Tick messages SHALL continue throughout the typing session
- Window resize messages SHALL update width/height for centering calculations

## Business Rules

### BR-001: Word Selection
- Words are selected randomly with replacement (same word may appear multiple times)
- All words in dictionary have equal probability of selection
- Empty words or whitespace-only words SHALL be skipped
- All words SHALL be converted to lowercase before use

### BR-002: WPM Calculation
- Formula: WPM = (correct_characters / 5) / minutes_elapsed
- Standard word length is defined as 5 characters
- Only correctly typed characters contribute to WPM
- Time measured from first correct keystroke to round completion

### BR-003: Accuracy Calculation
- Formula: Accuracy = (correct_characters / total_characters) × 100
- Every keystroke counts toward total_characters (including errors)
- Backspace removes the last character from consideration
- Extra characters beyond word length count as incorrect

### BR-004: Best Time Logic
- Best time is the LOWEST (fastest) completion time
- On first session, current time becomes best time
- Subsequent sessions only update best if time < current best

### BR-005: New Best Detection
- WPM: New best if current >= historical best
- Accuracy: New best if current >= historical best
- Time: New best if current <= historical best (lower is better)
- First session always counts as "new best" for all metrics

## Word Dictionary

### British English Spellings
The dictionary SHALL use British English spellings:
- colour (not color)
- behaviour (not behavior)
- centre (not center)
- defence (not defense)
- favour (not favor)
- realise (not realize)
- organisation (not organization)
- recognise (not recognize)
- programme (not program)
- labour, honour, neighbour
- travelling
- theatre, metre, litre, fibre
- Words ending in -ise (apologise, capitalise, emphasise, etc.)

## Architecture

The application follows a clean backend/frontend separation with a well-defined API:

### Backend Package (`backend/`)
The backend handles all game logic, statistics tracking, and state management. The frontend communicates exclusively through the `GameAPI` interface.

**Key Components:**
- `api.go` - Defines the `GameAPI` interface and data types
- `engine.go` - Implements the game engine

**GameAPI Interface:**
```go
type GameAPI interface {
    // Game Lifecycle
    StartRound()

    // Input Handling
    ProcessKeystroke(char string) KeystrokeResult
    ProcessBackspace() bool
    ProcessSpace() SpaceResult

    // State Queries
    GetGameState() GameState
    GetSessionStats() *stats.Stats
    GetHistoricalStats() *stats.HistoricalStats

    // Persistence
    SaveStats() error
}
```

### Frontend Package (`frontend/`)
The TUI frontend handles all rendering, user input, and visual presentation. It communicates with the backend exclusively through the `GameAPI` interface.

**Key Components:**
- `model.go` - Bubble Tea model (Init, Update, View)
- `views.go` - All rendering functions (typing screen, results screen)
- `styles.go` - Lipgloss styles and colour definitions
- `animations.go` - Spring-based animation logic

### Web Frontend (`web/`)
The web frontend is a React application that provides the same functionality as the TUI but with a beautiful, modern web interface.

**Technology Stack:**
- **React 18**: Modern React with hooks for state management
- **Chakra UI 2.x**: Component library with dark theme support
- **Framer Motion**: Physics-based animations with spring dynamics
- **Custom Theme**: Dark theme with Kartoza brand colours (orange #D4922A and blue #4A90A4)

**Key Components:**
- `App.js` - Main application component with state management and screen routing
- `api.js` - REST API client for backend communication
- `theme.js` - Custom Chakra UI theme with dark mode and chunky button styles
- `components/WelcomeScreen.js` - Landing page with animated logo and game options
- `components/TypingScreen.js` - Main typing interface with BlockLetter physics and live WPM bar
- `components/ResultsScreen.js` - Statistics display with all typing theory metrics

**Features:**
- Large chunky block letters with spring-based physics animations
- Real-time colour feedback (green for correct, red for incorrect)
- Live WPM bar with gradient colouring during typing
- Animated transitions between screens
- Letter statistics grid with accuracy and speed indicators
- Finger accuracy display and hand balance statistics
- Common error pattern display
- Responsive design for various screen sizes

**Physics Animations:**
- Letters bounce and scale on correct/incorrect keystrokes using Framer Motion springs
- Stat cards slide in with staggered spring animations
- Progress bars animate smoothly with spring dynamics
- UI elements have hover/tap scaling effects

### Data Flow
1. User input → Frontend Model → REST Client → HTTP → REST Server → Game Engine
2. Game state changes → REST Client queries API → HTTP Response → Render updated view
3. Statistics persist through REST Server → Stats package → JSON file

### REST API

The backend exposes a RESTful API that the frontend communicates with via HTTP. This enables:
- Clean separation between frontend and backend processes
- Multiple concurrent frontend sessions connecting to a single backend
- Potential for alternative clients (web UI, mobile app, etc.)
- Network-based or remote gameplay

**Base URL:** `http://127.0.0.1:8787` (configurable via `-port` flag)

#### Session Management

The REST API uses session-based routing. Each frontend client creates a session on startup and receives a unique session ID. All game operations are then scoped to that session.

**Session Lifecycle:**
1. Client calls `POST /api/sessions` to create a new session
2. Server responds with a unique `session_id` (32-character hex string)
3. Client includes session ID in all subsequent requests: `/api/sessions/{id}/...`
4. Client calls `DELETE /api/sessions/{id}` on exit to clean up

#### Endpoints

**Session Management:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/sessions` | Create a new session |
| DELETE | `/api/sessions/{id}` | Delete a session |
| GET | `/api/sessions` | List all active sessions |
| GET | `/api/health` | Health check (includes active session count) |

**Game Operations (session-scoped):**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/sessions/{id}/round` | Start a new round |
| POST | `/api/sessions/{id}/keystroke` | Process a keystroke (with timing) |
| POST | `/api/sessions/{id}/backspace` | Process backspace |
| POST | `/api/sessions/{id}/space` | Process space key (with timing) |
| POST | `/api/sessions/{id}/timing` | Submit final round timing data |
| GET | `/api/sessions/{id}/state` | Get current game state |
| GET | `/api/sessions/{id}/stats/session` | Get session statistics |
| GET | `/api/sessions/{id}/stats/historical` | Get historical statistics |
| POST | `/api/sessions/{id}/save` | Save statistics to disk |

#### Frontend Timing

All timing-critical measurements are performed on the frontend to avoid network latency affecting accuracy:

1. **Timer tracking**: The frontend tracks when the timer starts (first correct keystroke) and ends (round complete)
2. **Seek time measurement**: Time between keystrokes is measured locally and sent with each request
3. **Live WPM calculation**: Computed on the frontend using local timing data
4. **Duration calculation**: Total round duration is calculated on the frontend and submitted at round end

**How it works:**
- Each keystroke/space request includes a `seek_time_ms` field with the frontend-measured time since the previous keystroke
- When a round completes, the frontend calls `POST /api/sessions/{id}/timing` with:
  - `start_time_unix_ms`: Unix milliseconds when timer started
  - `end_time_unix_ms`: Unix milliseconds when round ended
  - `duration_ms`: Total duration in milliseconds
- The backend uses this timing data for WPM and accuracy calculations, ensuring network latency doesn't affect statistics

#### Request/Response Examples

**POST /api/sessions**
```json
// Request
{"punctuation_mode": false}

// Response (201 Created)
{"session_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"}
```

**GET /api/sessions**
```json
// Response
{
  "sessions": [
    {
      "id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
      "created_at": "2024-01-15T10:30:00Z",
      "last_used": "2024-01-15T10:35:00Z"
    }
  ]
}
```

**GET /api/health**
```json
// Response
{"status": "healthy", "active_sessions": 3}
```

**POST /api/sessions/{id}/keystroke**
```json
// Request
{"char": "a"}

// Response
{"is_correct": true, "timer_started": true, "char_index": 0}
```

**POST /api/sessions/{id}/space**
```json
// Response
{"advanced": true, "round_complete": false, "treated_as_error": false}
```

**GET /api/sessions/{id}/state**
```json
// Response
{
  "words": ["hello", "world", ...],
  "current_word_idx": 0,
  "current_input": "hel",
  "timer_started": true,
  "punctuation_mode": false,
  "word_number": 1,
  "total_words": 30,
  "live_wpm": 45.2,
  "current_word": "hello",
  "previous_word": "",
  "next_word": "world",
  "next_words": ["world", "typing", "test"]
}
```

## File Structure

```
baboon/
├── flake.nix           # Nix flake for cross-platform builds
├── flake.lock          # Nix flake lock file
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── main.go             # Entry point - supports server, client, and combined modes
├── backend/
│   ├── api.go          # GameAPI interface and types
│   ├── engine.go       # Game engine implementation
│   └── server.go       # REST API server with session management
├── frontend/
│   ├── model.go        # Bubble Tea model with local timing
│   ├── views.go        # Rendering functions
│   ├── styles.go       # Lipgloss styles
│   ├── animations.go   # Spring animation logic
│   └── client.go       # REST API client (implements GameAPI)
├── font/
│   └── font.go         # Block letter font definitions (a-z + punctuation)
├── words/
│   └── words.go        # Dictionary of common words (British English)
├── stats/
│   ├── stats.go        # Statistics types, persistence, validation
│   └── keyboard.go     # Keyboard layout mappings (finger, hand, row)
├── scripts/
│   ├── start-backend.sh   # Start backend server in background
│   ├── stop-backend.sh    # Stop backend server
│   ├── status-backend.sh  # Check backend status and health
│   └── launch-frontend.sh # Launch frontend client
├── web/                   # React web frontend
│   ├── package.json       # NPM dependencies
│   ├── package-lock.json  # NPM lockfile
│   ├── public/
│   │   └── index.html     # HTML template with fonts
│   └── src/
│       ├── index.js       # React entry point
│       ├── App.js         # Main application component
│       ├── api.js         # REST API client
│       ├── theme.js       # Chakra UI custom theme
│       └── components/
│           ├── WelcomeScreen.js   # Landing screen
│           ├── TypingScreen.js    # Typing practice screen
│           └── ResultsScreen.js   # Statistics display
├── Makefile            # Build and run targets
├── SPECIFICATION.md    # This file
├── README.md           # User documentation
├── LICENSE             # MIT license
└── .gitignore          # Git ignore patterns
```

## Running Modes

The application supports three running modes:

### Combined Mode (Default)
```bash
baboon              # Start backend + frontend together
baboon -p           # With punctuation mode
baboon -port 9000   # On custom port
```
Both backend and frontend run in the same process. When you exit, everything stops.

### Server-Only Mode
```bash
baboon -server              # Run backend only (blocking)
baboon -server -port 9000   # On custom port
```
Runs the REST API server in the foreground. Useful for running as a service or allowing multiple frontend connections. Writes PID to `$XDG_RUNTIME_DIR/baboon.pid`.

### Client-Only Mode
```bash
baboon -client              # Connect to existing backend
baboon -client -p           # With punctuation mode
baboon -client -port 9000   # Connect to custom port
```
Connects to an already-running backend server. Multiple clients can connect simultaneously, each with their own session.

### Web Frontend Mode
```bash
make web-start              # Start backend + web frontend
make web-dev                # Start web dev server only (needs backend running)
make web-build              # Build for production
```
Starts the React web frontend on port 3000. The frontend proxies API requests to the backend on port 8787.

## Management Scripts

Scripts are provided in the `scripts/` directory for managing the backend as a background service:

### start-backend.sh
Starts the backend server in the background.
```bash
./scripts/start-backend.sh           # Start on default port
./scripts/start-backend.sh -port 9000  # Custom port
./scripts/start-backend.sh -p        # With punctuation mode
```
- Checks if backend is already running
- Writes PID file for management
- Logs output to `$XDG_RUNTIME_DIR/baboon.log`

### stop-backend.sh
Stops the backend server gracefully.
```bash
./scripts/stop-backend.sh      # Graceful shutdown
./scripts/stop-backend.sh -f   # Force kill
```

### status-backend.sh
Checks the backend server status and health.
```bash
./scripts/status-backend.sh              # Check default port
./scripts/status-backend.sh -port 9000   # Check custom port
```
Shows: process status, health endpoint response, active session count.

### launch-frontend.sh
Launches a frontend client connected to the backend.
```bash
./scripts/launch-frontend.sh           # Connect to default port
./scripts/launch-frontend.sh -p        # With punctuation mode
./scripts/launch-frontend.sh -port 9000  # Connect to custom port
```
Checks that backend is running before launching.

## Stats File Format

Location: `~/.config/baboon/stats.json`

```json
{
  "best_wpm": 65.5,
  "best_accuracy": 98.2,
  "best_time": 45.3,
  "total_wpm": 850.5,
  "total_accuracy": 1420.8,
  "total_time": 725.0,
  "total_sessions": 15,
  "last_session_date": "2024-01-15T10:30:00Z",
  "letter_accuracy": {
    "a": {"presented": 100, "correct": 99},
    "b": {"presented": 45, "correct": 43}
  },
  "letter_seek_time": {
    "a": {"total_time_ms": 15000, "count": 100},
    "b": {"total_time_ms": 9000, "count": 45}
  },
  "bigram_seek_time": {
    "th": {"total_time_ms": 8500, "count": 50},
    "he": {"total_time_ms": 7200, "count": 48}
  },
  "finger_stats": {
    "0": {"presented": 200, "correct": 198, "total_time_ms": 30000, "count": 198},
    "1": {"presented": 180, "correct": 175, "total_time_ms": 27000, "count": 175}
  },
  "hand_stats": {
    "0": {"presented": 800, "correct": 790, "total_time_ms": 120000, "count": 790},
    "1": {"presented": 750, "correct": 740, "total_time_ms": 112500, "count": 740}
  },
  "row_stats": {
    "0": {"presented": 400, "correct": 395, "total_time_ms": 60000, "count": 395},
    "1": {"presented": 600, "correct": 595, "total_time_ms": 90000, "count": 595},
    "2": {"presented": 350, "correct": 340, "total_time_ms": 52500, "count": 340}
  },
  "error_substitution": {
    "a": {"s": 5, "q": 2},
    "e": {"r": 3, "w": 1}
  },
  "sfb_stats": {"count": 150, "total_time_ms": 45000},
  "hand_alternations": 1200,
  "same_hand_runs": 800,
  "rhythm_stats": {
    "total_seek_time_ms": 300000,
    "total_seek_time_sq": 75000000.0,
    "count": 2000,
    "last_variance": 0
  }
}
```

## Colour Palette Reference

### TUI Colour Codes

| Usage | Colour Code | Description |
|-------|-------------|-------------|
| Correct letter | 10 | Bright green |
| Incorrect letter | 9 | Bright red |
| Untyped letter | 8 | Gray |
| Title | 14 | Cyan |
| Labels | 7 | Light gray |
| Values | 15 | White |
| Session label | 6 | Cyan |
| New best star | 226 | Yellow |
| Help text | 8 | Gray |
| Empty bar | 236 | Dark gray |
| Gradient | 196→47 | Red through yellow to green |

### Web UI Kartoza Brand Colours

The web frontend uses Kartoza's brand colour scheme derived from their wallpaper artwork.

| Colour | Hex Code | Usage |
|--------|----------|-------|
| Kartoza Orange (Primary) | #D4922A | Primary brand colour, buttons, current letter highlight, WPM stat |
| Kartoza Blue (Secondary) | #4A90A4 | Secondary brand colour, hover states, alternation, time stat |
| Kartoza Gray Light | #C4C4C4 | Light gray accents |
| Kartoza Gray Medium | #9A9A9A | Medium gray accents |
| Kartoza Gray Dark | #6A6A6A | Dark gray, pending letters |
| Green (Correct) | #4CAF50 | Correct keystrokes, high accuracy |
| Red (Incorrect) | #E53935 | Incorrect keystrokes, errors |
| Background Primary | #1a2833 | Main background |
| Background Secondary | #243442 | Secondary background |
| Background Card | #1f3040 | Card backgrounds |

**Brand Colour Palette (Orange)**
- 50: #fef6e9
- 100: #fce8c7
- 200: #f9d9a5
- 300: #f5c983
- 400: #e8a93d
- 500: #D4922A (primary)
- 600: #b87a22
- 700: #9c631a
- 800: #804c12
- 900: #64350a

**Brand Colour Palette (Blue)**
- 50: #e9f4f7
- 100: #c7e3ea
- 200: #a5d2dd
- 300: #83c1d0
- 400: #61b0c3
- 500: #4A90A4 (primary)
- 600: #3d7688
- 700: #305c6c
- 800: #234250
- 900: #162834

The Kartoza wallpaper (`web/public/kartoza-wallpaper.png`) is included in the project assets for reference.

## Version History

### v1.2.0
- Show next 3 upcoming words in carousel display
  - Words displayed below the current word with decreasing opacity
  - First upcoming word has decorative arrows (▼), subsequent words shown plain
  - All upcoming words centered horizontally using lipgloss alignment
  - Provides better look-ahead for typing preparation
- Updated GameState API to include `NextWords` slice (array of up to 3 words)
- Backwards compatible: falls back to `NextWord` if `NextWords` is empty

### v1.1.0
- Beautiful carousel animation for word transitions
  - **Console (TUI)**: Smooth harmonica spring-based animations
    - Previous word fades in with animated greyscale opacity as it scrolls up
    - Current word slides up from below with spring physics
    - Next word fades in from below with staggered timing
    - Animation triggered on space key when advancing to next word
  - **Web**: Framer Motion spring animations
    - Previous word floats above at 50% scale with blur and fade
    - Current word displays large block letters with spring transitions
    - Next word floats below at 50% scale with blur
    - Decorative glow effect behind current word
- Fixed accuracy statistics exceeding 100%
  - Bug: When backspacing and retyping a character, "Correct" was counted multiple times while "Presented" was only counted once at round start
  - Fix: Track which character positions have been recorded as correct using a position map
  - Accuracy stats (letter, finger, hand, row) only recorded on first correct keystroke per position
  - Timing stats still recorded for all keystrokes (useful data regardless of retypes)
- Updated FR-001 to describe carousel word display layout

### v1.0.0
- First stable release
- Beautiful README with screenshots and badges
- GitHub Actions CI/CD workflows:
  - Test workflow: runs on push/PR, executes go test and go vet
  - Build workflow: cross-platform build verification (Linux, macOS, Windows)
  - Release workflow: automated builds on tag push with all package formats
- Pre-built binaries for multiple platforms:
  - Linux AMD64 and ARM64
  - macOS Intel and Apple Silicon
  - Windows AMD64
  - DEB package for Debian/Ubuntu
  - RPM package for Fedora/RHEL
  - Flatpak package
- macOS unsigned binary instructions in README
- Nix flake integration for system configurations

### v0.9.1
- Kartoza brand colour scheme applied to web frontend
  - Primary colour: Kartoza Orange (#D4922A)
  - Secondary colour: Kartoza Blue (#4A90A4)
  - Updated theme.js with full Kartoza colour palettes
  - Updated WelcomeScreen gradient title with brand colours
  - Updated TypingScreen with orange current letter highlight and blue/orange progress bars
  - Updated ResultsScreen with brand colours for stats, hand balance, and heatmaps
  - Kartoza wallpaper added to project assets (web/public/kartoza-wallpaper.png)

### v0.9.0
- React web frontend with 100% feature parity to TUI
  - Built with React 18, Chakra UI 2.x, and Framer Motion
  - Physics-based animations using Framer Motion spring dynamics
  - Large chunky block letters with bounce effects on keystrokes
  - Dark theme with Kartoza brand colours (orange and blue)
  - Custom theme with chunky button styles ("glow" and "chunky" variants)
- Web frontend components:
  - WelcomeScreen: Animated logo, connection status, game options
  - TypingScreen: Block letters with physics, live WPM bar, progress indicator
  - ResultsScreen: Full statistics display with animated stat cards
- Letter statistics grid with colour-coded accuracy and speed indicators
- Finger accuracy display and hand balance statistics
- Common error pattern tracking and display
- Responsive design for various screen sizes
- Makefile targets for web development:
  - `make web-install` - Install NPM dependencies
  - `make web-dev` - Start development server
  - `make web-build` - Build for production
  - `make web-start` - Start backend + web frontend together
- Proxy configuration for development (port 3000 → 8787)

### v0.8.0
- Management scripts for backend server lifecycle
  - `start-backend.sh` - Start backend in background with PID tracking
  - `stop-backend.sh` - Graceful or forced shutdown
  - `status-backend.sh` - Health check and session monitoring
  - `launch-frontend.sh` - Launch frontend against running backend
- Three running modes added:
  - Combined mode (default): Backend + frontend in same process
  - Server-only mode (`-server`): Run backend only, blocking
  - Client-only mode (`-client`): Connect to existing backend
- PID file written to `$XDG_RUNTIME_DIR/baboon.pid` in server mode
- Graceful shutdown handling with SIGINT/SIGTERM

### v0.7.0
- Frontend timing implementation to eliminate network latency effects
  - All timing-critical measurements now performed on the frontend
  - Seek times measured locally and sent with keystroke/space requests
  - Live WPM calculated on frontend using local timing data
  - Round duration submitted via dedicated timing endpoint
- New API methods:
  - `ProcessKeystrokeWithTiming(char, seekTimeMs)` - keystroke with timing
  - `ProcessSpaceWithTiming(seekTimeMs)` - space with timing
  - `SubmitTiming(startTime, endTime, durationMs)` - submit final timing
- New REST endpoint: `POST /api/sessions/{id}/timing`
- Backend no longer calls `time.Now()` for timing-critical operations
- Ensures accurate WPM and seek time statistics regardless of network conditions

### v0.6.0
- Multi-client session management for REST API
  - Each frontend client creates a unique session on startup
  - Session IDs are 32-character hex strings generated cryptographically
  - All game operations are scoped to sessions: `/api/sessions/{id}/...`
  - Sessions are automatically cleaned up when clients disconnect
- New session management endpoints:
  - `POST /api/sessions` - Create a new session
  - `DELETE /api/sessions/{id}` - Delete a session
  - `GET /api/sessions` - List all active sessions
- Health endpoint now reports active session count
- Thread-safe session storage with mutex protection
- Enables multiple concurrent players on a single backend server

### v0.5.0
- RESTful API implementation for frontend-backend communication
  - All frontend interactions routed through HTTP REST API
  - Backend server (`backend/server.go`) exposes REST endpoints
  - Frontend client (`frontend/client.go`) implements GameAPI via HTTP
- REST API endpoints:
  - `POST /api/round` - Start new round
  - `POST /api/keystroke` - Process keystroke
  - `POST /api/backspace` - Process backspace
  - `POST /api/space` - Process space
  - `GET /api/state` - Get game state
  - `GET /api/stats/session` - Get session statistics
  - `GET /api/stats/historical` - Get historical statistics
  - `POST /api/save` - Save statistics
  - `GET /api/health` - Health check
- Configurable port via `-port` flag (default: 8787)
- Thread-safe server with mutex protection

### v0.4.0
- Major architecture refactoring: clean backend/frontend separation
  - Backend package (`backend/`): Game engine with `GameAPI` interface
  - Frontend package (`frontend/`): TUI model, views, styles, animations
  - Main.go reduced to simple entry point wiring backend and frontend
- Clear API boundary: frontend communicates only through `GameAPI` interface
- Improved code organisation and maintainability

### v0.3.0
- Advanced typing theory statistics for effective touch typing practice:
  - Finger-specific accuracy and speed tracking (8 fingers mapped to QWERTY layout)
  - Keyboard row tracking (top, home, bottom row performance)
  - Hand balance and alternation rate tracking
  - Same-finger bigram (SFB) detection and timing
  - Rhythm consistency tracking (standard deviation of seek times)
  - Error substitution pattern tracking (which letters get confused)
- Results screen displays new typing theory metrics:
  - Finger accuracy row: LP LR LM LI | RI RM RR RP with colour-coded indicators
  - Row accuracy: Top, Home, Bot performance indicators
  - Hand balance: L:X% R:Y% distribution with alternation rate
  - Rhythm: Session StdDev vs historical average
  - Same-finger: SFB count and average timing
  - Common errors: Top 5 letter substitution patterns (e.g., a→s(12))
- Increased animated rows from 14 to 25 for new stats sections

### v0.2.0
- Per-letter accuracy tracking (tracks how often each letter a-z is presented and typed correctly)
- Letter statistics persist across sessions for cumulative tracking
- Results screen displays 26-letter accuracy row with red-to-green gradient
- Results screen displays 26-letter frequency row showing relative letter presentation counts
- Results screen displays 26-letter seek time row showing typing speed per letter
- Fixed 30 words with exactly 150 characters per round for consistent timing comparisons
- Smooth font with Unicode half-block characters (▀, ▄) for rounded edges
- Previous word displayed in top left, next word in top right during typing
- Adaptive word selection: weights by both letter frequency AND accuracy
- Words with low-accuracy letters are favoured to give more practice on weak letters
- Letter seek time tracking: measures time between keystrokes for each letter
- Improved seek time calculation:
  - Only records for correct keystrokes (not errors)
  - Records against expected letter (not what user typed)
  - Excludes first letter of each word (avoids word-reading time)
- Bigram (letter pair) seek time tracking: measures transition speed between letter pairs
- Letter statistics display redesigned: header row with letters, filled circles (●) for data
- Results screen animation: rows slide in sequentially using harmonica spring physics
- Punctuation mode (-p flag): words separated by random punctuation (, . ; : ! ?)

### v1.3.0
- MkDocs documentation site with Material theme
  - Comprehensive user documentation
  - Developer guides and API reference
  - Beautiful baboon mascot typing at keyboard
  - GitHub Pages deployment via GitHub Actions
  - Dark/light theme support
  - Full-text search
  - Responsive mobile design

### v0.1.0 (Initial Release)
- Basic typing practice with 30-word rounds
- Block letter word display with real-time colour feedback
- Live WPM bar during typing
- Results screen with paired comparison bars
- WPM, time, and accuracy tracking with averages
- Historical best comparison with star indicators
- British English word dictionary
- Cross-platform support via Nix flakes
- Statistics persistence with corruption detection
