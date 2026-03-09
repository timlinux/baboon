---
title: "How to Play"
description: "Master the basics of Baboon"
weight: 1
icon: "🎮"
---

Baboon is designed to be intuitive, but here's a complete guide to playing.

## The Basics

### Starting a Round

1. Launch Baboon: `baboon`
2. A word appears in large block letters
3. The timer starts when you type the first correct letter

### Typing Words

1. Type each letter as you see it
2. Letters change colour based on correctness:
   - **Gray**: Pending (not yet typed)
   - **Orange**: Current position
   - **Green**: Correct
   - **Red**: Incorrect

3. Press {{< kbd >}}Space{{< /kbd >}} when you've typed the entire word
4. The next word appears

### Completing a Round

- Each round has 30 words totalling 150 characters
- After the 30th word, you see the results screen
- Press {{< kbd >}}Enter{{< /kbd >}} to start a new round

## Keyboard Controls

### During Typing

| Key | Action |
|-----|--------|
| Any letter | Type that character |
| {{< kbd >}}Space{{< /kbd >}} | Advance to next word (when complete) |
| {{< kbd >}}Backspace{{< /kbd >}} | Delete last character |
| {{< kbd >}}Ctrl+W{{< /kbd >}} | Clear entire word |
| {{< kbd >}}Tab{{< /kbd >}} | Restart current round |
| {{< kbd >}}Ctrl+O{{< /kbd >}} | Open options (before timer starts) |
| {{< kbd >}}Esc{{< /kbd >}} | Exit application |

### On Results Screen

| Key | Action |
|-----|--------|
| {{< kbd >}}Enter{{< /kbd >}} or {{< kbd >}}Tab{{< /kbd >}} | Start new round |
| {{< kbd >}}Ctrl+O{{< /kbd >}} | Open options |
| {{< kbd >}}Esc{{< /kbd >}} | Exit application |

## The Display

### Header Section

Shows:
- Application title
- Current word progress (e.g., "Word 15/30")

### Main Word Area

The carousel display shows:
- Previous word (dimmed, above)
- Current word (large, centred)
- Next word (dimmed, below)

### Footer Section

Shows:
- Live WPM bar with gradient colours
- Current WPM value
- Keyboard hints

## Understanding the WPM Bar

The WPM bar uses colour to indicate speed:

| Colour | WPM Range |
|--------|-----------|
| Red | Below 30 |
| Orange | 30-50 |
| Yellow | 50-70 |
| Green | 70-90 |
| Blue | Above 90 |

## Tips for Success

### Focus on Accuracy

Speed follows accuracy. If you're making lots of mistakes, slow down until you can type accurately, then gradually increase speed.

### Don't Look Down

Keep your eyes on the screen. Looking at the keyboard breaks your flow and slows you down.

### Use Proper Technique

- Keep fingers on the home row (ASDF JKL;)
- Use the correct finger for each key
- Type with a light touch

### Practice Regularly

Short, frequent sessions (10-15 minutes) are better than long, exhausting ones. Baboon tracks your progress across sessions.

## Command Line Options

```bash
# Normal mode
baboon

# With punctuation
baboon -p

# Custom port
baboon -port 9000

# Server only (for multiple clients)
baboon -server

# Connect to existing server
baboon -client
```
