# Punctuation Mode

Master punctuation for real-world typing skills with Baboon's punctuation mode.

## Enabling Punctuation Mode

### Command Line

```bash
baboon -p
```

### What Changes

In punctuation mode, words are separated by random punctuation followed by a space:

**Normal mode**:

```
hello world typing practice
```

**Punctuation mode**:

```
hello, world. typing; practice:
```

## Supported Punctuation

| Character | Name | Keyboard Position |
|-----------|------|-------------------|
| `,` | Comma | Right of M |
| `.` | Period/Full stop | Right of comma |
| `;` | Semicolon | Right of L |
| `:` | Colon | Shift + semicolon |
| `!` | Exclamation mark | Shift + 1 |
| `?` | Question mark | Shift + / |

## Typing Punctuation

### How It Works

1. Type the word normally
2. Type the punctuation character
3. Press ++space++ to advance

**Example**: Typing "hello,"

| Step | Input | Display |
|------|-------|---------|
| 1 | h | <span class="key correct">h</span>ello, |
| 2 | e | <span class="key correct">he</span>llo, |
| 3 | l | <span class="key correct">hel</span>lo, |
| 4 | l | <span class="key correct">hell</span>o, |
| 5 | o | <span class="key correct">hello</span>, |
| 6 | , | <span class="key correct">hello,</span> |
| 7 | ++space++ | Next word! |

### Punctuation on Keyboard

```
  1 2 3 4 5 6 7 8 9 0 - =
  ! @ # $ % ^ & * ( ) _ +     ← Shift row
   Q W E R T Y U I O P [ ]
    A S D F G H J K L ; '     ← ; here
     Z X C V B N M , . /      ← , . here
```

## Statistics in Punctuation Mode

### What's Tracked

Letter statistics **only count a-z**:

- Punctuation doesn't affect letter accuracy
- Punctuation doesn't affect letter seek time
- Word scores still based on letters only

### Why?

Punctuation appears randomly, so tracking it wouldn't provide useful adaptive data. The focus remains on letter proficiency.

### What's Still Measured

- Overall WPM (includes punctuation typing time)
- Overall accuracy (punctuation errors count!)
- Time to complete round

## Why Practice Punctuation?

### Real-World Typing

Most text includes punctuation:

- Emails and messages
- Code and programming
- Documents and reports

### Reach Training

Punctuation keys require reaching:

- Right pinky for `;` and `'`
- Shift key combinations
- Right hand stretching

### Flow Disruption Practice

Punctuation breaks normal flow:

- Learn to handle interruptions
- Build recovery speed
- Improve transitions

## Punctuation Techniques

### Comma and Period

These are the most common:

- Right ring finger for `.`
- Right middle finger for `,`
- Quick tap, don't pause

**Practice words**:

```
hello, world. typing, test. good, practice.
```

### Semicolon and Colon

Right pinky territory:

- `;` - pinky, no shift
- `:` - pinky + left shift

**Practice words**:

```
example; testing: another; more:
```

### Exclamation and Question

Require shift:

- `!` - left pinky (shift) + left pinky (1)
- `?` - left pinky (shift) + right pinky (/)

!!! tip "Shift Technique"
    Use the opposite hand's shift key when possible. For `?`, use left shift with right hand's `/`.

## Exercises

### Exercise 1: Comma Flow

Focus on smooth comma typing:

```
one, two, three, four, five,
red, blue, green, yellow, orange,
```

### Exercise 2: Period Practice

End every sentence:

```
done. next. typing. practice. better.
```

### Exercise 3: Mixed Punctuation

Combine different marks:

```
hello! how? are; you: today,
```

### Exercise 4: Real Sentences

Practise sentence-like patterns:

```
the quick brown fox jumps. over the lazy dog!
```

## Common Challenges

### Forgetting Punctuation

**Problem**: Pressing space before typing punctuation

**Solution**: Read the full word including punctuation before typing

### Slow Shift Key

**Problem**: Pause before shifted characters (!?)

**Solution**: Practice shift combinations separately:

```
!! ?? !! ?? :: ::
```

### Wrong Punctuation

**Problem**: Typing `,` instead of `.`

**Solution**: Read ahead more carefully, slow down initially

## Progression Path

### Phase 1: Accuracy

1. Enable punctuation mode
2. Type slowly and deliberately
3. Focus on 100% punctuation accuracy
4. Don't worry about WPM

### Phase 2: Recognition

1. Practise reading punctuation in words
2. Build mental preparation
3. Reduce pause before punctuation

### Phase 3: Speed

1. Increase overall pace
2. Integrate punctuation smoothly
3. Maintain 95%+ accuracy

### Phase 4: Flow

1. Punctuation becomes automatic
2. No noticeable pause
3. Treat word+punctuation as unit

## Statistics Interpretation

### With Punctuation

Your WPM might be:

- 5-10 WPM lower initially
- Different character distribution
- More shift key usage

This is normal! Punctuation adds complexity.

### Comparing Modes

Track separately:

- Normal mode personal bests
- Punctuation mode personal bests

Don't compare directly - they're different skills.

## Tips for Success

1. **Start slow** - Accuracy first
2. **Read ahead** - See the punctuation coming
3. **Practise shifts** - Get comfortable with shift combos
4. **Be patient** - Speed will come
5. **Mix modes** - Alternate between normal and punctuation

## Command Reference

```bash
# Punctuation mode
baboon -p

# Punctuation + server mode
baboon -p -server

# Punctuation + custom port
baboon -p -port 9000
```

Punctuation mode persists for the entire session (all rounds).

## Next Steps

- [How to Play](how-to-play.md) - Basic techniques
- [Improving Speed](improving-speed.md) - Speed tips
- [Understanding Stats](understanding-stats.md) - Interpreting results
