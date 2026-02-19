# Baboon - Typing Practice Application

## Overview

Baboon is a cross-platform terminal-based typing practice application built with Go and Bubble Tea. It helps users improve their typing speed and accuracy by presenting common English words in large ASCII art format.

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
**I want to** compare my current performance to my historical best
**So that** I can see how I'm improving

## Functional Requirements

### FR-001: Word Display
- The application SHALL display words from a dictionary of the 1000 most common English words
- Each word SHALL be rendered in large block characters using Unicode block elements (█)
- Words SHALL be displayed centered on the terminal screen
- The word display SHALL show progress (current word number / total words)
- Letters SHALL change color in-place as the user types (no separate input display)

### FR-002: Typing Input
- The application SHALL accept keyboard input character by character
- Each character typed SHALL be immediately reflected in the display
- Characters SHALL be colored based on correctness:
  - **Green**: Character matches the expected character at that position
  - **Red**: Character does not match the expected character
  - **Gray/White**: Characters not yet typed
- The backspace key SHALL remove the last typed character
- The space key SHALL advance to the next word (regardless of correctness)

### FR-003: Round Structure
- Each round SHALL consist of exactly 30 words
- Words SHALL be randomly selected from the dictionary
- After completing 30 words, the application SHALL display statistics
- The user SHALL be able to start a new round by pressing Enter

### FR-004: Statistics Display
- At the end of each round, the application SHALL display:
  - Words Per Minute (WPM)
  - Accuracy percentage
  - Total time elapsed
  - Total characters typed
- WPM SHALL be calculated as: (correct characters / 5) / minutes elapsed
- Accuracy SHALL be calculated as: (correct characters / total characters) * 100

### FR-007: Live WPM Bar
- The application SHALL display a gradient WPM bar at the bottom of the screen during typing
- The bar SHALL update in real-time (every 100ms) to show current typing speed
- The bar SHALL use a gradient colour scheme:
  - Red (0-40 WPM) - needs improvement
  - Yellow (40-60 WPM) - average
  - Green (60+ WPM) - good speed
- The bar SHALL scale from 0 to 120 WPM
- The bar SHALL display the current WPM value numerically alongside the bar

### FR-005: Historical Statistics
- The application SHALL persist best performance data between sessions
- Historical data SHALL be stored in `~/.config/baboon/stats.json`
- The application SHALL track:
  - Best WPM achieved
  - Best accuracy achieved
  - Total number of sessions completed
  - Date of last session
- After each round, the application SHALL compare current performance to historical best:
  - Display "NEW BEST!" in green when a new record is achieved
  - Display difference from best in red when below the record

### FR-006: Navigation
- ESC or Ctrl+C SHALL exit the application at any time
- SPACE SHALL advance to the next word during typing
- ENTER SHALL start a new round when viewing results

## Technical Requirements

### TR-001: Cross-Platform Compatibility
- The application SHALL run on Linux, macOS, and Windows
- The application SHALL be buildable using Nix flakes for reproducible builds

### TR-002: Terminal Interface
- The application SHALL use the Bubble Tea framework for the TUI
- The application SHALL use lipgloss for styling
- The application SHALL use custom block font rendering with Unicode block characters
- The application SHALL use the alternate screen buffer (fullscreen mode)

### TR-003: Dependencies
- github.com/charmbracelet/bubbletea - TUI framework
- github.com/charmbracelet/lipgloss - Styling

## Business Rules

### BR-001: Word Selection
- Words are selected randomly with replacement (same word may appear multiple times)
- All 1000 words have equal probability of selection

### BR-002: Accuracy Calculation
- Only characters that have been typed are counted toward accuracy
- Extra characters beyond the word length count as incorrect
- Backspace removes characters from consideration (they don't count as errors)

### BR-003: WPM Calculation
- Standard word length is defined as 5 characters
- Only correctly typed characters contribute to WPM
- Time starts when the first correct character of the first word is typed (not when the word is displayed)

### BR-004: Historical Best
- Best WPM and Best Accuracy are tracked independently
- A session can set a new best in one category without affecting the other
- First session establishes the baseline for comparison

## File Structure

```
baboon/
├── flake.nix           # Nix flake for building
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── main.go             # Application entry point and TUI logic
├── font/
│   └── font.go         # Custom block letter font definitions
├── words/
│   └── words.go        # Dictionary of 1000 common words
├── stats/
│   └── stats.go        # Statistics tracking and persistence
├── SPECIFICATION.md    # This file
└── README.md           # User documentation
```

## Configuration

### Storage Location
- Linux/macOS: `~/.config/baboon/stats.json`
- The directory is created automatically if it doesn't exist

### Stats File Format (JSON)
```json
{
  "best_wpm": 65.5,
  "best_accuracy": 98.2,
  "total_sessions": 15,
  "last_session_date": "2024-01-15T10:30:00Z"
}
```

## Version History

### v0.1.0 (Initial Release)
- Basic typing practice with 30-word rounds
- ASCII art word display with colored feedback
- WPM and accuracy tracking
- Historical best comparison
- Cross-platform support via Nix flakes
