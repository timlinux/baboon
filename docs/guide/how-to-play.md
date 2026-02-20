# How to Play

Master Baboon with this comprehensive guide to typing practice.

## The Basics

### Starting a Round

When Baboon launches, you'll see:

1. A large word displayed in block letters
2. A word counter (Word 1/30)
3. An empty WPM bar

**The timer hasn't started yet!** Take a moment to:

- Position your hands on the home row
- Read the first word
- Prepare mentally

### Typing

Begin typing the displayed word:

- The timer starts on your **first correct keystroke**
- Letters turn <span class="key correct">green</span> when correct
- Letters turn <span class="key incorrect">red</span> when incorrect
- Keep typing until the word is complete

!!! tip "Timer Tip"
    If you mistype the first character, the timer won't start. This prevents accidental early starts.

### Advancing to the Next Word

Press ++space++ to move to the next word when:

- You've typed all letters (correct or not)
- You're ready to proceed

**Important**: If you press space before completing the word, it counts as an error!

### Completing a Round

After 30 words:

1. The timer stops automatically
2. Statistics are calculated
3. The results screen appears

Press ++enter++ to start a new round.

## Controls Reference

### During Typing

| Key | Action |
|-----|--------|
| ++a++ - ++z++ | Type the corresponding letter |
| ++backspace++ | Delete the last typed character |
| ++space++ | Move to next word (when complete) |
| ++escape++ | Exit Baboon |

### On Results Screen

| Key | Action |
|-----|--------|
| ++enter++ | Start a new round |
| ++escape++ | Exit Baboon |

## Proper Technique

### Hand Position

Use the standard touch typing home row position:

```
Left hand:  A S D F (index on F)
Right hand: J K L ; (index on J)
```

Your fingers should rest lightly on these keys, ready to reach for others.

### Finger Assignments

| Finger | Keys |
|--------|------|
| Left Pinky | Q A Z |
| Left Ring | W S X |
| Left Middle | E D C |
| Left Index | R F V T G B |
| Right Index | Y H N U J M |
| Right Middle | I K |
| Right Ring | O L |
| Right Pinky | P |

!!! warning "Don't Look!"
    Looking at the keyboard slows you down. Trust your muscle memory!

### Posture

- Sit up straight
- Elbows at 90-degree angle
- Wrists floating above keyboard
- Eyes on the screen

## Strategies

### For Beginners

**Goal**: Build accuracy first

1. Focus on hitting the correct keys
2. Don't worry about speed initially
3. Use proper finger positioning
4. Take your time

**Target metrics**:

- Accuracy: > 95%
- WPM: Any

### For Intermediate

**Goal**: Build speed while maintaining accuracy

1. Push yourself a little faster
2. Accept some accuracy drop (temporarily)
3. Identify weak letters in the heatmap
4. Practice problem areas

**Target metrics**:

- Accuracy: > 90%
- WPM: 40-60

### For Advanced

**Goal**: Optimise rhythm and consistency

1. Work on typing flow
2. Minimise pauses between words
3. Develop consistent rhythm
4. Target specific weak bigrams

**Target metrics**:

- Accuracy: > 95%
- WPM: 60+
- Low StdDev (rhythm)

## Common Mistakes

### Looking at the Keyboard

**Problem**: Slows you down, creates dependency

**Solution**: Cover the keyboard or use a blank keyboard if needed. Trust your training!

### Incorrect Finger Usage

**Problem**: "Hunt and peck" or using wrong fingers

**Solution**: Consciously use the correct finger for each key, even if it feels slower initially

### Rushing

**Problem**: High speed, low accuracy

**Solution**: Slow down! Accuracy first. Speed comes from muscle memory, not conscious effort.

### Tensing Up

**Problem**: Tired hands, slower typing

**Solution**: Keep hands relaxed. Light touches are faster than heavy presses.

## Practice Tips

### Daily Practice

Even 10-15 minutes daily is more effective than occasional long sessions.

### Progressive Goals

Set achievable milestones:

1. Complete 5 rounds with >90% accuracy
2. Reach 50 WPM average
3. Get all letters to green in the heatmap
4. Achieve <100ms average seek time

### Use Punctuation Mode

Once comfortable with letters, add punctuation:

```bash
baboon -p
```

This adds `, . ; : ! ?` between words.

### Track Your Progress

Pay attention to:

- Which letters are red/orange in the heatmap
- Your error substitution patterns
- Your average vs best scores
- Rhythm consistency (StdDev)

## Understanding Feedback

### The WPM Bar

The live WPM bar shows your current pace:

```
▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂▂░░░░░░░░░░░░░░░░  45 WPM
0              60              120
```

- **Red zone** (0-40): Below average
- **Yellow zone** (40-60): Average
- **Green zone** (60+): Above average

Don't obsess over it - it's meant for general awareness.

### Colour Changes

Watch how letters change as you type:

- Immediate green = great muscle memory
- Red followed by correction = that's learning!
- Consistent reds = focus on that letter

### The Word Carousel

See what's coming:

- **Above**: Word you just finished (dimmed)
- **Center**: Current word (bright)
- **Below**: Next 3 words (dimmed)

Use peripheral vision to prepare for upcoming words.

## Troubleshooting

### "I keep making the same mistake"

Check your error patterns on the results screen. Common causes:

- Adjacent key confusion (e→r, i→o)
- Finger position drift
- Rushing through familiar sequences

### "My speed plateaued"

Try:

- Focus on weak letters (red in heatmap)
- Practice punctuation mode
- Work on rhythm consistency
- Take a break and return fresh

### "My accuracy dropped"

You might be:

- Pushing speed too fast
- Getting tired
- Using incorrect fingers

Slow down and rebuild proper technique.

## Next Steps

- [Understanding Stats](understanding-stats.md) - Interpret your results
- [Improving Speed](improving-speed.md) - Advanced techniques
- [Punctuation Mode](punctuation-mode.md) - Add punctuation practice
