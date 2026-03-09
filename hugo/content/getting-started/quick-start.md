---
title: "Quick Start"
description: "Start practicing in under a minute"
weight: 2
icon: "🚀"
---

Get typing in under a minute with this quick start guide!

## Terminal Mode

The simplest way to start:

```bash
baboon
```

That's it! You'll see a word displayed in large block letters. Start typing!

### What You'll See

1. **The Word** - Displayed in large block letters in the center
2. **Progress** - Shows "Word X/30" in the header
3. **WPM Bar** - Live words-per-minute displayed at the bottom

### How to Play

1. Type the letters you see
2. Letters turn **green** when correct, **red** when wrong
3. Press {{< kbd >}}Space{{< /kbd >}} when you've typed the whole word
4. Complete 30 words to finish the round
5. View your statistics on the results screen

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| {{< kbd >}}Space{{< /kbd >}} | Advance to next word |
| {{< kbd >}}Backspace{{< /kbd >}} | Delete last character |
| {{< kbd >}}Ctrl+W{{< /kbd >}} | Clear entire word |
| {{< kbd >}}Tab{{< /kbd >}} | Restart round |
| {{< kbd >}}Ctrl+O{{< /kbd >}} | Open options |
| {{< kbd >}}Esc{{< /kbd >}} | Exit |

## Web Mode

For a beautiful browser-based experience:

```bash
# Build the web frontend first
cd web && npm install && npm run build && cd ..

# Start the web server
baboon web -port 8080
```

Then open http://localhost:8080 in your browser.

## Options

### Punctuation Mode

Add punctuation between words for extra challenge:

```bash
baboon -p
```

### Custom Port

Run on a different port:

```bash
baboon -port 9000
```

## Next Steps

- Learn about [Your First Session](/getting-started/first-session/)
- Explore the [Features](/features/)
- Read the [User Guide](/guide/)
