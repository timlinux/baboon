# Terminal Interface

The Terminal UI (TUI) is Baboon's original interface - a beautiful, efficient way to practice typing right in your terminal.

![Terminal Typing Screen](../screenshots/console1.png)

## Technology

The TUI is built with the excellent [Charm](https://charm.sh/) libraries:

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - The Elm-inspired TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - CSS-like styling for terminals
- **[Harmonica](https://github.com/charmbracelet/harmonica)** - Spring physics for smooth animations

## Block Letter Font

Words are displayed using a custom 6-line tall font built from Unicode block characters:

| Character | Usage |
|-----------|-------|
| `█` | Full block (solid areas) |
| `▀` | Upper half block (rounded tops) |
| `▄` | Lower half block (rounded bottoms) |

This creates smooth, readable letters even in a terminal:

```
    ▄▄▄▄▄▄▄▄  █████▄  █████▄  ▄█████  ▄█████  █     █
    █     ▄▀  █    █  █    █  █    █  █    █  █▄    █
    █████▀    █████▀  █████▀  █    █  █    █  █ ▀▄  █
    █    ▀▄   █    █  █    █  █    █  █    █  █   ▀▄█
    █████▀▀   █████▀  █████▀  ▀█████  ▀█████  █     █
```

The font supports:

- Lowercase letters a-z
- Punctuation: `, . ; : ! ?`
- Unknown characters render as spaces

## Colour Scheme

The TUI uses 256-colour mode for rich visual feedback:

| Element | Colour Code | Description |
|---------|-------------|-------------|
| Correct letter | 10 | Bright green |
| Incorrect letter | 9 | Bright red |
| Untyped letter | 8 | Gray |
| Title | 14 | Cyan |
| Labels | 7 | Light gray |
| Values | 15 | White |
| New best star | 226 | Yellow |

### Gradient Colours

The WPM bar and statistics bars use a gradient from red to green:

```
196 → 202 → 208 → 214 → 220 → 226 → 190 → 154 → 118 → 82 → 46 → 47
```

## Word Carousel Animation

When you press space to advance to the next word, a smooth animation plays:

1. **Previous word** fades in above with animated greyscale opacity
2. **Current word** slides up with spring physics
3. **Next word** fades in from below

Animation parameters:

- Frame rate: 60 FPS
- Spring frequency: 6.0
- Spring damping: 0.5
- Stagger delay: 3 frames

## Layout

### Typing Screen

```
┌──────────────────────────────────────────────────────┐
│                                                      │
│                     previous                         │ ← Dimmed
│                                                      │
│     ██████ ██   ██ ██████ ██████ ███████ ██   ██     │
│     ██     ██   ██ ██   █ ██   █ ██      ████ ██     │ ← Current
│     ██     ██   ██ ██████ ██████ █████   ██ ████     │    word
│     ██     ██   ██ ██ ██  ██ ██  ██      ██  ███     │
│     ██████ ███████ ██  ██ ██  ██ ███████ ██   ██     │
│                                                      │
│                      ▼ next ▼                        │ ← Dimmed
│                       word2                          │
│                       word3                          │
│                                                      │
│     Word 15/30                                       │
│                                                      │
│     ▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂░░░░░░░░░░░░░░   52 WPM      │
│     0              60              120               │
│                                                      │
└──────────────────────────────────────────────────────┘
```

### Results Screen

The results screen displays statistics with animated slide-in:

- Each row slides in from the right
- 50ms interval between frames
- 3-frame stagger between rows
- 25 total animated rows

## Running the Terminal UI

### Combined Mode (Default)

```bash
baboon
```

Runs backend and frontend together in one process.

### Server/Client Mode

For running the backend separately (useful for multiple sessions):

```bash
# Terminal 1: Start server
baboon -server

# Terminal 2: Connect client
baboon -client
```

### With Punctuation

```bash
baboon -p
```

Words will have random punctuation between them.

## Keyboard Controls

| Key | Action |
|-----|--------|
| a-z | Type the next character |
| ++backspace++ | Remove last typed character |
| ++space++ | Advance to next word |
| ++enter++ | Start new round (results screen) |
| ++escape++ or ++ctrl+c++ | Exit |

## Terminal Requirements

- **256-colour support** - Most modern terminals work
- **Unicode support** - For block characters
- **Minimum width**: ~80 columns
- **Minimum height**: ~24 rows

Recommended terminals:

- **Linux**: kitty, Alacritty, GNOME Terminal
- **macOS**: iTerm2, kitty, Terminal.app
- **Windows**: Windows Terminal, ConEmu

## Fullscreen Mode

Baboon uses `tea.WithAltScreen()` for fullscreen mode:

- Enters alternate screen buffer on start
- Restores original screen on exit
- No scrollback pollution

## Next Steps

- [Web Interface](web.md) - Try the web UI
- [Statistics](statistics.md) - Understanding your results
- [How to Play](../guide/how-to-play.md) - Detailed playing guide
