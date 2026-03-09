---
title: "Feature Overview"
description: "A comprehensive look at Baboon's capabilities"
weight: 1
icon: "✨"
---

Baboon is more than just a typing test - it's a comprehensive typing improvement tool with features designed to help you become a faster, more accurate typist.

## Core Features

### Beautiful Block Letters

Words are displayed using large Unicode block characters that change colour in real-time as you type:

- **Green** for correct keystrokes
- **Red** for mistakes
- **Orange** highlighting the current position

### Real-time WPM Tracking

Watch your speed as you type with a live WPM bar that updates continuously, giving you immediate feedback on your pace.

### Comprehensive Statistics

Track everything:
- Words per minute
- Accuracy percentage
- Time per round
- Personal bests and averages

### Per-Letter Tracking

Detailed statistics for each letter (a-z):
- Accuracy rate
- Presentation frequency
- Average seek time (speed)

## Advanced Features

### Adaptive Word Selection

Baboon learns your weaknesses! Words containing letters you frequently mistype are prioritised, giving you more practice where you need it most.

### Finger-Specific Stats

Track performance by finger with standard touch-typing assignments. See which fingers are fastest and most accurate.

### Row Statistics

Compare your performance across keyboard rows:
- Top row (qwerty...)
- Home row (asdf...)
- Bottom row (zxcv...)

### Hand Balance

Monitor left vs. right hand usage and alternation rate for better typing flow.

### Same-Finger Bigram Detection

Identify letter combinations typed with the same finger (SFBs) - a common source of slowdowns.

### Rhythm Analysis

Track your typing consistency with standard deviation measurements.

### Error Pattern Detection

See which letters you commonly confuse (e.g., typing 's' instead of 'a').

## Interface Options

### Terminal UI

A beautiful TUI built with Bubble Tea:
- Works in any terminal
- Smooth animations
- Responsive design

### Web Interface

A modern React-based web app:
- Physics-based animations
- Kartoza brand styling
- Works in any browser

### Punctuation Mode

Add punctuation between words for extra challenge:
```bash
baboon -p
```

## Data & Privacy

### Local Storage

All your data stays on your machine:
- Statistics saved to `~/.config/baboon/`
- No accounts required
- No cloud sync
- No tracking

### Session Persistence

Your progress accumulates across sessions, building a comprehensive picture of your typing abilities over time.

## Platform Support

Baboon runs everywhere:
- **Linux**: Native binaries, Deb, RPM, Flatpak
- **macOS**: Intel and Apple Silicon
- **Windows**: x86_64 binaries
- **Nix**: Flake-based installation
