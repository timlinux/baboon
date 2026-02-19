# Baboon

A terminal-based typing practice application that helps you improve your typing speed and accuracy.

## Features

- Practice typing with the 1000 most common English words
- Large ASCII art display for each word
- Real-time color feedback (green for correct, red for incorrect)
- WPM (Words Per Minute) and accuracy tracking
- Historical best comparison across sessions
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Using Nix Flakes

```bash
# Run directly
nix run github:timlinux/baboon

# Or install to your profile
nix profile install github:timlinux/baboon
```

### From Source

```bash
git clone https://github.com/timlinux/baboon.git
cd baboon
nix build
./result/bin/baboon
```

Or with Go directly:

```bash
go build -o baboon .
./baboon
```

## Usage

1. Launch the application
2. Type the displayed word character by character
3. Characters turn **green** when correct, **red** when incorrect
4. Press **SPACE** to move to the next word
5. After 30 words, view your statistics
6. Press **ENTER** to start a new round
7. Press **ESC** or **Ctrl+C** to quit

## Statistics

After each round, you'll see:
- Words Per Minute (WPM)
- Accuracy percentage
- Time elapsed
- Total characters typed
- Comparison to your historical best

Your best scores are saved to `~/.config/baboon/stats.json` and persist between sessions.

## Development

```bash
# Enter development shell
nix develop

# Run tests
go test ./...

# Build
go build -o baboon .
```

## License

MIT
