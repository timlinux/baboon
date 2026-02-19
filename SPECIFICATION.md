# Baboon - Typing Practice Application

## Overview

Baboon is a cross-platform terminal-based typing practice application built with Go and Bubble Tea. It helps users improve their typing speed and accuracy by presenting common English words in large block letter format. Words are displayed using Unicode block characters (█) and change colour in real-time as the user types.

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

## Functional Requirements

### FR-001: Word Display
- The application SHALL display words from a dictionary of common English words (British English spelling)
- Each word SHALL be rendered in large block characters using Unicode block elements (█)
- Words SHALL be displayed centered horizontally and vertically on the terminal screen
- The word display SHALL show progress indicator: "Word X/30"
- Letters SHALL change colour in-place as the user types (no separate input display line)
- All words SHALL be lowercase only (the font only supports a-z)
- The previous word SHALL be displayed in gray (colour 240) in the top left corner
- The next word SHALL be displayed in gray (colour 240) in the top right corner

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
- Total of 14 animated rows (WPM, Time, Accuracy stats + sessions + letter matrix)

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

## File Structure

```
baboon/
├── flake.nix           # Nix flake for cross-platform builds
├── flake.lock          # Nix flake lock file
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── main.go             # Application entry point, TUI logic, rendering
├── font/
│   └── font.go         # Block letter font definitions (a-z + space)
├── words/
│   └── words.go        # Dictionary of common words (British English)
├── stats/
│   └── stats.go        # Statistics types, persistence, validation
├── SPECIFICATION.md    # This file
├── README.md           # User documentation
├── LICENSE             # MIT license
└── .gitignore          # Git ignore patterns
```

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
    "b": {"presented": 45, "correct": 43},
    "c": {"presented": 62, "correct": 60}
  },
  "letter_seek_time": {
    "a": {"total_time_ms": 15000, "count": 100},
    "b": {"total_time_ms": 9000, "count": 45},
    "c": {"total_time_ms": 12400, "count": 62}
  },
  "bigram_seek_time": {
    "th": {"total_time_ms": 8500, "count": 50},
    "he": {"total_time_ms": 7200, "count": 48},
    "in": {"total_time_ms": 9100, "count": 35}
  }
}
```

## Colour Palette Reference

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

## Version History

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
