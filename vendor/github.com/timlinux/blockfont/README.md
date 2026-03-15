# blockfont

Unicode block letter rendering for terminal applications.

![blockfont demo](demo/blockfont-demo.gif)

```
██████ ██     ◢████◣ ◢████◣ ██ ◢█◤ ██████ ◢████◣ ██◣ ██ ◥█  █◤
██  ██ ██     ██  ██ ██  ██ ██◢█◤  ██     ██  ██ ███◣██  ◥██◤
██████ ██     ██  ██ ██     ███◣   ████   ██  ██ ██◥███   ██
██  ██ ██     ██  ██ ██  ◢█ ██◥█◣  ██     ██  ██ ██ ◥██   ██
██████ ██████ ◥████◤ ◥████◤ ██ ◥█◣ ██     ◥████◤ ██  ██   ██
```

## Features

- **Block Letter Rendering**: Convert text to beautiful Unicode block letters using `█ ◢ ◣ ◤ ◥`
- **Full Character Set**: Lowercase, uppercase, numbers, and common punctuation
- **Vim-Style Editing**: Full vim buffer with normal/insert modes and common operations
- **Character Highlighting**: Color individual characters for typing games, diffs, etc.
- **Spring Animations**: Smooth physics-based transitions using harmonica
- **Word Wrapping**: Intelligent word-boundary wrapping with cursor tracking
- **Lipgloss Integration**: Full charmbracelet/lipgloss styling support
- **Bubbletea Widget**: Ready-to-use `tea.Model` for easy integration
- **Layout Utilities**: Centering, alignment, and positioning helpers

## Installation

```bash
go get github.com/timlinux/blockfont
```

## Quick Start

### Simple Rendering

```go
package main

import (
    "fmt"
    "github.com/timlinux/blockfont"
)

func main() {
    // Render text as block letters
    text := blockfont.RenderText("hello")
    fmt.Println(text)
}
```

### Widget with Bubbletea

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/timlinux/blockfont"
)

type model struct {
    widget *blockfont.Widget
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    w, cmd := m.widget.Update(msg)
    m.widget = w
    return m, cmd
}

func (m model) View() string {
    return m.widget.View()
}

func main() {
    opts := blockfont.DefaultWidgetOptions()
    opts.Width = 80
    opts.VimMode = true
    opts.Theme = blockfont.KartozaTheme

    widget := blockfont.NewWidget(opts)
    widget.SetText("edit me")
    widget.Focus()

    tea.NewProgram(model{widget: widget}).Run()
}
```

### Character Highlighting

Perfect for typing games and diffs:

```go
text := "typing"
highlights := []blockfont.CharHighlight{
    blockfont.HighlightCorrect,   // t - green
    blockfont.HighlightCorrect,   // y - green
    blockfont.HighlightIncorrect, // p - red
    blockfont.HighlightPending,   // i - dim
    blockfont.HighlightPending,   // n - dim
    blockfont.HighlightPending,   // g - dim
}

lines := blockfont.RenderWithCursor(
    text, -1, highlights, false, 0, blockfont.DefaultTheme)
```

## Run the Demo

Try out all features interactively:

```bash
# Clone the repository
git clone https://github.com/timlinux/blockfont
cd blockfont

# Run with make
make demo

# Or with nix
nix run
```

The demo includes 8 screens showcasing:
- Basic character rendering
- Character highlighting
- Vim-style editing
- Spring animations
- Word wrapping
- Theme options

## Development

```bash
# Enter nix development shell
nix develop

# Build and test
make build
make test

# Run demo
make demo

# Start documentation server
make docs-dev

# Record a demo
nix run .#demo-record

# Play recorded demo
nix run .#demo-play
```

## API Reference

### Core Functions

| Function | Description |
|----------|-------------|
| `RenderText(text string) string` | Render text with newlines |
| `RenderWord(word string) [][]string` | Render as 2D array |
| `RenderWithCursor(...)` | Full rendering with cursor and highlights |
| `RenderPlainText(text, color string)` | Simple colored rendering |
| `GetLetterWidth(char rune) int` | Character width |
| `GetTotalWidth(word string) int` | Word width |

### Widget

| Method | Description |
|--------|-------------|
| `NewWidget(opts WidgetOptions) *Widget` | Create widget |
| `SetText(text string)` | Set content |
| `Text() string` | Get content |
| `View() string` | Render (tea.Model) |
| `Update(msg tea.Msg)` | Handle input (tea.Model) |
| `Focus() / Blur()` | Manage keyboard focus |
| `Mode() Mode` | Get vim mode |

### Layout

| Function | Description |
|----------|-------------|
| `CenterLines(lines, width)` | Center text |
| `LeftJustify(lines, margin)` | Left-align |
| `RightJustify(lines, width)` | Right-align |
| `WrapOnWordBoundaries(text, width)` | Word wrap |

### Animation

| Method | Description |
|--------|-------------|
| `NewAnimator()` | Create animator |
| `TriggerTransition(type)` | Start animation |
| `Update() bool` | Advance frame |
| `GetOpacity() float64` | Current opacity |
| `GetOffset() int` | Current offset |

## Supported Characters

- **Lowercase**: a-z
- **Uppercase**: A-Z
- **Digits**: 0-9
- **Punctuation**: `. , ; : ! ? - ' " ( ) [ ] { } < > / \ | _ + = * & @ # $ % ^ ` ~`
- **Space**

## Themes

```go
// Default theme (green/red/cyan)
opts.Theme = blockfont.DefaultTheme

// Kartoza theme (orange/blue)
opts.Theme = blockfont.KartozaTheme

// Custom theme
opts.Theme = blockfont.Theme{
    Correct:   blockfont.ANSIGreen,
    Incorrect: blockfont.ANSIRed,
    Cursor:    blockfont.ANSICyan,
    Pending:   blockfont.ANSIDim,
}
```

## Documentation

Full documentation is available at the [documentation site](https://timlinux.github.io/blockfont).

- [Getting Started](https://timlinux.github.io/blockfont/user/getting-started/)
- [Widget Guide](https://timlinux.github.io/blockfont/user/widget/)
- [Vim Editing](https://timlinux.github.io/blockfont/user/vim-editing/)
- [Animations](https://timlinux.github.io/blockfont/user/animations/)
- [API Reference](https://timlinux.github.io/blockfont/developer/api/)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

---

Made with :heart: by [Kartoza](https://kartoza.com) | [Donate!](https://github.com/sponsors/timlinux) | [GitHub](https://github.com/timlinux/blockfont)
