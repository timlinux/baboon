# Quick Start

Get typing with Baboon in under a minute!

## Launch Baboon

=== "Terminal UI"

    ```bash
    baboon
    ```

=== "Web UI"

    ```bash
    # Start backend + web frontend
    make web-start

    # Or manually:
    ./baboon -server &
    cd web && npm start
    ```

    Then open http://localhost:3000

## Your First Round

When Baboon starts, you'll see a large word displayed in block letters:

```
    ██   ██ ███████ ██      ██       ██████
    ██   ██ ██      ██      ██      ██    ██
    ███████ █████   ██      ██      ██    ██
    ██   ██ ██      ██      ██      ██    ██
    ██   ██ ███████ ███████ ███████  ██████
```

### How to Play

1. **Start typing** - The timer begins when you type the first correct character
2. **Watch the colours**:
   - <span class="key correct">Green</span> = Correct
   - <span class="key incorrect">Red</span> = Incorrect
   - <span class="key">Gray</span> = Not yet typed
3. **Press SPACE** - Move to the next word when you've typed all letters
4. **Complete 30 words** - View your statistics
5. **Press ENTER** - Start a new round
6. **Press ESC** - Quit at any time

## Basic Controls

| Key | Action |
|-----|--------|
| Any letter | Type the next character |
| ++backspace++ | Remove the last typed character |
| ++space++ | Move to the next word (when current word is complete) |
| ++enter++ | Start a new round (on results screen) |
| ++escape++ | Exit Baboon |

## Command Line Options

```bash
# Standard mode
baboon

# With punctuation practice
baboon -p

# Server mode (for multiple clients)
baboon -server

# Connect to existing server
baboon -client

# Custom port
baboon -port 9000
```

## Understanding the Display

### During Typing

```
                    Word 5/30

    ██████  ██████  ██████  █████  ████████
    ██   ██ ██   ██ ██      ██   ██    ██
    ██████  ██████  ██████  ███████    ██
    ██      ██   ██ ██      ██   ██    ██
    ██      ██   ██ ██████  ██   ██    ██

    ▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂░░░░░░░░░░  52 WPM
    0                60              120
```

- **Word counter** - Shows your progress through the 30 words
- **Block letters** - The current word to type
- **WPM bar** - Real-time typing speed indicator

### Results Screen

After completing 30 words, you'll see:

- **WPM** - Words per minute (higher is better)
- **Accuracy** - Percentage of correct keystrokes
- **Time** - How long the round took
- **Letter statistics** - Accuracy heatmap for each letter
- **Finger accuracy** - How each finger performed
- **Hand balance** - Left vs right hand usage

## Pro Tips

!!! tip "Baboon's Typing Tips"

    1. **Don't look at the keyboard** - Trust your muscle memory
    2. **Aim for accuracy first** - Speed will follow naturally
    3. **Use all your fingers** - Proper touch typing technique matters
    4. **Practice daily** - Even 10 minutes a day helps
    5. **Check your weak letters** - The heatmap shows where to focus

## What's Next?

- [Your First Session](first-session.md) - Detailed walkthrough of a complete session
- [Understanding Stats](../guide/understanding-stats.md) - Learn what all the numbers mean
- [Improving Speed](../guide/improving-speed.md) - Tips to type faster
