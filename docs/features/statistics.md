# Statistics

Baboon tracks comprehensive typing statistics to help you understand your strengths and weaknesses. Here's everything you need to know about the metrics.

## Core Metrics

### Words Per Minute (WPM)

**Formula**: `WPM = (correct_characters / 5) / minutes`

- Standard word length is defined as 5 characters
- Only correctly typed characters count
- Time measured from first correct keystroke to round completion

<span class="badge wpm">WPM</span> Higher is better!

| Rating | WPM Range | Description |
|--------|-----------|-------------|
| Beginner | < 30 | Just starting out |
| Average | 30-50 | Typical typing speed |
| Good | 50-70 | Above average |
| Fast | 70-90 | Professional level |
| Expert | 90+ | Elite typist |

### Accuracy

**Formula**: `Accuracy = (correct_characters / total_characters) × 100`

- Every keystroke counts toward total
- Backspace removes the last character from consideration
- Extra characters beyond word length count as incorrect

<span class="badge accuracy">Accuracy</span> Higher is better!

| Rating | Accuracy | Description |
|--------|----------|-------------|
| Needs work | < 90% | Focus on accuracy first |
| Good | 90-95% | Solid foundation |
| Excellent | 95-98% | Highly accurate |
| Perfect | 99%+ | Elite accuracy |

### Time

**Formula**: Time from first correct keystroke to space after final word

<span class="badge time">Time</span> Lower is better!

Each round is exactly 150 characters, making times directly comparable.

## Statistical Comparisons

For each core metric, Baboon displays:

| Value | Meaning |
|-------|---------|
| **This run** | Your current session result |
| **Best** | Your personal record |
| **Average** | Mean of all your sessions |

### New Personal Best

When you beat a record, you'll see a star ⭐ indicator:

- **WPM**: New best if current ≥ historical best
- **Accuracy**: New best if current ≥ historical best
- **Time**: New best if current ≤ historical best (lower is better)

## Per-Letter Statistics

### Letter Accuracy

Each letter (A-Z) has its own accuracy tracking:

- **Presented**: Times this letter appeared
- **Correct**: Times you typed it correctly
- **Accuracy**: `correct / presented × 100`

The results screen shows a colour-coded heatmap:

| Accuracy | Colour |
|----------|--------|
| 95-100% | :material-circle:{ style="color: #4CAF50" } Bright green |
| 90-94% | :material-circle:{ style="color: #8BC34A" } Light green |
| 85-89% | :material-circle:{ style="color: #CDDC39" } Lime |
| 80-84% | :material-circle:{ style="color: #FFEB3B" } Yellow |
| 75-79% | :material-circle:{ style="color: #FFC107" } Amber |
| 70-74% | :material-circle:{ style="color: #FF9800" } Orange |
| 65-69% | :material-circle:{ style="color: #FF5722" } Deep orange |
| 60-64% | :material-circle:{ style="color: #f44336" } Red |
| < 60% | :material-circle:{ style="color: #d32f2f" } Dark red |

### Letter Frequency

Shows how often each letter has been presented relative to others. Baboon tries to balance letter frequency so you practice all letters equally.

### Letter Seek Time

Measures how quickly you reach each letter:

- **Seek time**: Milliseconds between previous keystroke and this one
- Only recorded for **correct** keystrokes
- First letter of each word is **excluded** (includes word-reading time)
- Times > 5000ms are filtered out (assumed pauses)

!!! tip "Interpreting Seek Time"
    Faster seek times indicate better muscle memory for that key position.

## Typing Theory Metrics

### Finger Accuracy

Each finger's performance is tracked separately:

| Finger | Code | Keys |
|--------|------|------|
| Left Pinky | LP | q, a, z |
| Left Ring | LR | w, s, x |
| Left Middle | LM | e, d, c |
| Left Index | LI | r, f, v, t, g, b |
| Right Index | RI | y, h, n, u, j, m |
| Right Middle | RM | i, k |
| Right Ring | RR | o, l |
| Right Pinky | RP | p |

This uses standard QWERTY touch typing positions.

### Keyboard Row Statistics

Performance by row:

| Row | Keys |
|-----|------|
| **Top** | q, w, e, r, t, y, u, i, o, p |
| **Home** | a, s, d, f, g, h, j, k, l |
| **Bottom** | z, x, c, v, b, n, m |

!!! info "Home Row Advantage"
    Your home row should generally be fastest - that's where your fingers rest!

### Hand Balance

Tracks left vs right hand usage:

- **Left hand**: q-t, a-g, z-b
- **Right hand**: y-p, h-l, n-m

Displayed as: `L:48% R:52%`

Ideal balance depends on the language and word selection.

### Alternation Rate

**Formula**: `alternations / (alternations + same_hand_runs) × 100`

- **Alternation**: Switching hands between keystrokes
- **Same-hand run**: Consecutive keystrokes with same hand

Higher alternation rate generally indicates smoother typing flow.

### Same-Finger Bigrams (SFB)

An SFB occurs when consecutive letters use the same finger:

- Example: "un" (both typed with right index)
- SFBs are inherently slower than alternating fingers
- Baboon tracks count and average time

Common SFBs to watch for: `ed`, `de`, `un`, `nu`, `ec`, `ce`

### Rhythm Consistency

Measures typing evenness using standard deviation of seek times:

- **Lower StdDev** = More consistent rhythm
- **Higher StdDev** = More variable timing

Professional typists tend to have very consistent rhythm.

## Error Pattern Tracking

### Error Substitutions

Baboon records which letters get confused:

```
Top errors: e→r(5) a→s(3) i→o(2) t→y(2) n→m(1)
```

This shows:

- You typed 'r' when meaning to type 'e' five times
- Adjacent keys are commonly confused
- This data persists across sessions

!!! tip "Using Error Patterns"
    Focus practice on your common substitutions. If you frequently type 'a' as 's', slow down on words containing 'a'.

## Statistics Persistence

All statistics are saved to:

```
~/.config/baboon/stats.json
```

### What's Saved

```json
{
  "best_wpm": 65.5,
  "best_accuracy": 98.2,
  "best_time": 45.3,
  "total_wpm": 850.5,
  "total_accuracy": 1420.8,
  "total_time": 725.0,
  "total_sessions": 15,
  "letter_accuracy": { ... },
  "letter_seek_time": { ... },
  "bigram_seek_time": { ... },
  "finger_stats": { ... },
  "error_substitution": { ... }
}
```

### Data Validation

On load, Baboon validates statistics for corruption:

- If totals are 0 but bests exist → data is reset
- If average WPM < half of best WPM → data is reset

## Using Statistics Effectively

### Weekly Review

1. Check your per-letter accuracy heatmap
2. Identify your weakest letters (red/orange)
3. Note any recurring error patterns
4. Compare current average to personal bests

### Focus Areas

| If you struggle with... | Try... |
|------------------------|--------|
| Left pinky (q, a, z) | Practice words with "question", "amazing" |
| Bottom row | Words with "excellent", "amazing" |
| Speed but not accuracy | Slow down deliberately |
| Accuracy but not speed | Push yourself faster |
| Rhythm consistency | Use a metronome while practicing |

## Next Steps

- [Understanding Stats](../guide/understanding-stats.md) - Practical interpretation
- [Improving Speed](../guide/improving-speed.md) - Tips for faster typing
- [Adaptive Learning](adaptive.md) - How Baboon uses your stats
