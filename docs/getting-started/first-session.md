# Your First Session

Let's walk through a complete Baboon session from start to finish. By the end, you'll understand exactly how Baboon helps you improve your typing.

## Starting Up

When you launch Baboon, you'll be greeted with your first word displayed in large block letters:

<div class="terminal-window">
  <div class="terminal-header">
    <span class="terminal-button red"></span>
    <span class="terminal-button yellow"></span>
    <span class="terminal-button green"></span>
  </div>
  <div class="terminal-content">
    <pre style="color: #666; text-align: center;">
                              Word 1/30

  ▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄  ▄▄▄     ▄▄▄▄▄▄  ▄▄  ▄▄  ▄▄▄▄▄▄
  ██    ██  ██  ██  ██      ██  ██  ██  ██  ██  ██
  ██        ██  ██  ██      ██  ██  ██  ██  ████▀▀
  ██    ▄▄  ██  ██  ██      ██  ██  ██  ██  ██
  ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀  ▀▀▀▀▀▀  ▀▀▀▀▀▀  ▀▀▀▀▀▀  ██  ██
    </pre>
  </div>
</div>

Notice:

- The word counter shows "Word 1/30"
- The timer hasn't started yet
- All letters are in gray (untyped)

## Typing Your First Word

!!! info "Timer Behaviour"
    The timer only starts when you type the **first correct character**. If you mistype, the timer won't start until you get it right!

As you type each letter:

- **Correct letters** turn <span class="key correct">green</span>
- **Incorrect letters** turn <span class="key incorrect">red</span>
- Letters change colour instantly as you type

### Example: Typing "colour"

| You type | Display | Timer |
|----------|---------|-------|
| (nothing) | `colour` (all gray) | Not started |
| `c` | <span class="key correct">c</span>`olour` | Started! |
| `co` | <span class="key correct">co</span>`lour` | Running |
| `col` | <span class="key correct">col</span>`our` | Running |
| `colk` (typo!) | <span class="key correct">col</span><span class="key incorrect">k</span>`ur` | Running |
| ++backspace++ | <span class="key correct">col</span>`our` | Running |
| `colo` | <span class="key correct">colo</span>`ur` | Running |
| `colou` | <span class="key correct">colou</span>`r` | Running |
| `colour` | <span class="key correct">colour</span> | Running |
| ++space++ | Next word! | Running |

## The Word Carousel

As you type, you'll see context around your current word:

- **Above**: The previous word you just completed (dimmed)
- **Center**: The current word you're typing (bright)
- **Below**: The next 3 upcoming words (dimmed)

This helps you prepare for what's coming next without losing focus on the current word.

## The Live WPM Bar

At the bottom of the screen, a colourful bar shows your current typing speed:

```
▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂░░░░░░░░░░  52 WPM
0                60              120
```

The bar uses a gradient colour scheme:

- **Red** (left): Slow speeds (0-40 WPM)
- **Yellow** (middle): Moderate speeds (40-60 WPM)
- **Green** (right): Fast speeds (60+ WPM)

The bar updates every 100ms, giving you instant feedback on your pace.

## Completing the Round

After typing all 30 words, the results screen appears with a wealth of statistics.

### Core Statistics

| Metric | Description | Goal |
|--------|-------------|------|
| **WPM** | Words per minute | Higher is better |
| **Accuracy** | Correct keystrokes / Total keystrokes | Higher is better |
| **Time** | Total round duration | Lower is better |

Each metric shows:

- Your current session result
- Your personal best (with a star if you beat it!)
- Your historical average

### Letter Accuracy Heatmap

```
  A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
  ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ●
```

Each circle is coloured by accuracy:

- <span style="color: #4CAF50">**Green**</span> (95%+): Excellent
- <span style="color: #FFC107">**Yellow**</span> (70-94%): Needs work
- <span style="color: #f44336">**Red**</span> (<70%): Focus here!

### Typing Theory Statistics

Baboon tracks advanced metrics:

- **Finger accuracy**: How each finger performs (LP, LR, LM, LI, RI, RM, RR, RP)
- **Hand balance**: Left vs right hand distribution
- **Alternation rate**: How often you switch hands (higher is better!)
- **Same-finger bigrams**: Letter pairs typed with the same finger (slower)
- **Rhythm consistency**: Standard deviation of your typing speed

### Common Errors

The results show your top 5 most frequent mistakes:

```
Top errors: e→r(5) a→s(3) i→o(2) t→y(2) n→m(1)
```

This tells you that you typed 'r' when you meant 'e' five times!

## Starting a New Round

Press ++enter++ on the results screen to start a fresh round.

!!! tip "Adaptive Learning"
    Baboon remembers your mistakes! The next round will prioritise words containing letters you struggle with, giving you targeted practice.

## Ending Your Session

Press ++escape++ at any time to exit. Your statistics are automatically saved to:

```
~/.config/baboon/stats.json
```

This file persists between sessions, tracking your long-term progress.

## What You've Learned

After your first session, you now know:

- [x] How to start and stop Baboon
- [x] How the timer works
- [x] What the colour feedback means
- [x] How to read the WPM bar
- [x] How to interpret your statistics
- [x] Where your data is saved

## Next Steps

Ready to dive deeper?

- [Understanding Stats](../guide/understanding-stats.md) - Deep dive into all metrics
- [Improving Speed](../guide/improving-speed.md) - Tips from typing experts
- [Punctuation Mode](../guide/punctuation-mode.md) - Practice with punctuation
