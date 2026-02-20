# Understanding Your Statistics

After each round, Baboon presents a wealth of data. Here's how to interpret it and use it to improve.

## The Results Screen

The results screen is divided into several sections:

1. Core Statistics (WPM, Time, Accuracy)
2. Letter Statistics Matrix
3. Typing Theory Metrics
4. Error Patterns

Let's break down each section.

## Core Statistics

### Reading the Bars

Each metric shows three values with gradient bars:

```
      WPM this run:    52.3  ▂▂▂▂▂▂▂▂▂▂▂▂▂▂░░░░░░░░░░░
          WPM best:    65.5  ▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂░░░░░░░░ *
       WPM average:    48.2  ▂▂▂▂▂▂▂▂▂▂▂▂░░░░░░░░░░░░░
```

- **This run**: Your current session
- **Best**: Your personal record (⭐ if you beat it)
- **Average**: Mean of all sessions

### WPM (Words Per Minute)

**What it measures**: Raw typing speed

**Scale**: 0-120 WPM

**How to interpret**:

| WPM | Interpretation |
|-----|----------------|
| < 30 | Beginner - focus on technique |
| 30-50 | Average - building skill |
| 50-70 | Good - above average typist |
| 70-90 | Fast - professional level |
| 90+ | Expert - elite speed |

### Time

**What it measures**: Round completion duration

**Scale**: 0-180 seconds (inverted - lower is better)

**Note**: All rounds are exactly 150 characters, so times are directly comparable.

**How to interpret**:

| Time | For 150 chars | Implied WPM |
|------|---------------|-------------|
| 60s | 1 minute | 30 WPM |
| 45s | Fast | 40 WPM |
| 30s | Very fast | 60 WPM |
| 24s | Expert | 75 WPM |

### Accuracy

**What it measures**: Percentage of correct keystrokes

**Scale**: 0-100%

**How to interpret**:

| Accuracy | Interpretation |
|----------|----------------|
| < 85% | Needs work - slow down |
| 85-90% | Learning - acceptable while building speed |
| 90-95% | Good - solid foundation |
| 95-98% | Excellent - professional level |
| 99%+ | Elite - minimal errors |

!!! tip "The Accuracy/Speed Trade-off"
    It's normal for accuracy to dip when pushing speed. Aim for 95%+ when practicing technique, allow 90%+ when speed training.

## Letter Statistics Matrix

### The Display

```
  A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
  ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ●  Accuracy
  ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ●  Frequency
  ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ● ●  Seek Time
```

### Row Meanings

**Accuracy Row**: How often you type each letter correctly

- Green = 95%+ accuracy (great!)
- Yellow = 70-94% (needs work)
- Red = < 70% (focus here!)

**Frequency Row**: How often each letter has appeared

- Green = well-represented
- Red = underrepresented (Baboon will show it more)

**Seek Time Row**: How quickly you reach each letter

- Green = fast (<150ms average)
- Yellow = moderate (150-250ms)
- Red = slow (>250ms) - needs practice

### What to Look For

1. **Red accuracy letters**: Your weak spots
2. **Patterns**: Adjacent reds might indicate hand position issues
3. **Rare letters**: Q, Z, X often need extra attention

## Typing Theory Metrics

### Finger Accuracy

```
Fingers: LP LR LM LI | RI RM RR RP
         ●  ●  ●  ●    ●  ●  ●  ●
```

**Finger codes**:

| Code | Finger | Common Issues |
|------|--------|---------------|
| LP | Left Pinky | Often weakest, Q/A/Z |
| LR | Left Ring | W/S/X |
| LM | Left Middle | Usually strong |
| LI | Left Index | Reaches far, R/T/G/B |
| RI | Right Index | Reaches far, Y/U/H/N |
| RM | Right Middle | Usually strong |
| RR | Right Ring | O/L |
| RP | Right Pinky | P, often weak |

!!! note "Pinky Problems"
    Pinkies are typically the weakest fingers. If yours are red, consider pinky-specific exercises.

### Row Accuracy

```
Rows: Top Home Bot
      ●   ●    ●
```

**Expected pattern**:

- **Home row** should be strongest (fingers rest here)
- **Top row** often second
- **Bottom row** can be tricky (reaching down)

**If home row is weak**: Your hand positioning may be off.

### Hand Balance

```
Hands: L:47% R:53%
```

Shows the distribution of keystrokes between hands.

**Interpretation**:

- 45-55% split is normal for English text
- Heavily imbalanced? Check word content or finger technique

### Alternation Rate

```
Alternation: 68%
```

**What it measures**: How often you switch hands between keystrokes

**Why it matters**:

- Higher alternation = smoother typing flow
- 60-70% is typical for English
- Low alternation may indicate same-hand bigram struggles

### Same-Finger Bigrams (SFB)

```
Same-finger: 23 (avg 245ms)
```

**What it measures**: Letter pairs typed with the same finger

**Examples**: "de", "ed", "un", "nu"

**Why it matters**:

- SFBs are inherently slower
- High average time = that finger is struggling
- Practice words with common SFBs

### Rhythm Consistency

```
Rhythm: StdDev 85ms (avg: 78ms)
```

**What it measures**: Standard deviation of your seek times

**Interpretation**:

| StdDev | Meaning |
|--------|---------|
| < 50ms | Very consistent - professional rhythm |
| 50-100ms | Good consistency |
| 100-150ms | Moderate variability |
| > 150ms | Inconsistent - work on rhythm |

**Lower is better** - indicates even, predictable typing.

## Error Patterns

### The Display

```
Top errors: e→r(5) a→s(3) i→o(2) t→y(2) n→m(1)
```

### Reading Error Patterns

Format: `expected→typed(count)`

- `e→r(5)`: You typed 'r' when meaning to type 'e', 5 times
- These are your most common mistakes

### What Patterns Tell You

**Adjacent key errors** (e→r, i→o): Finger precision issue

**Same-finger errors** (e→d): Finger confusion

**Mirror errors** (f→j): Hand confusion

### Using Error Data

1. Note your top 3 error patterns
2. Practice words containing those letters
3. Slow down when approaching those keys
4. Consciously think about correct finger placement

## Practical Analysis

### Example Session Analysis

```
Session Results:
- WPM: 52 (best: 58, avg: 48)
- Accuracy: 94% (best: 98%, avg: 95%)
- Time: 57s

Letter Heatmap: Red on Q, Z, X
Finger Accuracy: LP weak (red)
Error Pattern: q→w(3), z→x(2)
```

**Analysis**:

1. Speed is above average - good progress
2. Accuracy slightly below average - pushed too hard
3. Left pinky (LP) is the issue
4. Q, Z, X are all left pinky keys!

**Action plan**:

1. Slow down slightly for accuracy
2. Focus on left pinky exercises
3. Practice words with Q, Z, X

### Weekly Review Checklist

- [ ] Compare average WPM to last week
- [ ] Check for new personal bests
- [ ] Identify any declining letter accuracies
- [ ] Note persistent error patterns
- [ ] Review finger/hand balance

## Next Steps

- [Improving Speed](improving-speed.md) - Techniques for faster typing
- [Punctuation Mode](punctuation-mode.md) - Add punctuation practice
- [Adaptive Learning](../features/adaptive.md) - How Baboon targets your weaknesses
