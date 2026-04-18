# N-gram Training Mode Design

**Date:** 2026-04-18
**Status:** Approved
**Author:** Claude Code

## Overview

Add a new practice mode to Baboon for focused bigram and trigram training. Users type random n-gram sequences (like "th er ing on") with adaptive selection prioritising their weakest combinations based on historical seek time data.

## Requirements Summary

| Requirement | Decision |
|-------------|----------|
| Content | Random n-gram sequences without forming words |
| Selection | Adaptive - prioritise user's slowest n-grams |
| Mode type | Combined bigrams + trigrams in single mode |
| Round structure | Fixed 150 characters (matches word mode) |
| Separator | Space between n-grams |
| Statistics | Separate tracking from word mode |
| Access | Options screen (Ctrl+O) |

## Architecture

### Practice Mode Enum

```go
type PracticeMode int
const (
    ModeWords PracticeMode = iota  // Default: 30 words, 150 chars
    ModeNgrams                      // N-gram training: bigrams + trigrams
)
```

### New Package: `ngrams/`

Create `ngrams/ngrams.go` containing:
- N-gram generation and selection logic
- Takes historical `BigramSeekTime` and `TrigramSeekTime` data
- Generates sequences targeting user's slowest n-grams
- Returns space-separated strings like "th er ing on the"

### Engine Changes

- `Config` struct gains `PracticeMode` field
- `StartRound()` branches based on mode:
  - `ModeWords`: calls `words.GetRandomWordsFixedCount()`
  - `ModeNgrams`: calls `ngrams.GetTrainingSequence()`
- All keystroke/space processing remains identical (n-grams separated by spaces work like words)

### Stats Structure Additions

```go
type HistoricalStats struct {
    // ... existing fields ...

    // Trigram tracking (new - shared across modes)
    TrigramSeekTime map[string]TrigramSeekStats `json:"trigram_seek_time"`

    // N-gram mode separate tracking
    NgramBestWPM       float64 `json:"ngram_best_wpm"`
    NgramBestAccuracy  float64 `json:"ngram_best_accuracy"`
    NgramBestTime      float64 `json:"ngram_best_time"`
    NgramTotalWPM      float64 `json:"ngram_total_wpm"`
    NgramTotalAccuracy float64 `json:"ngram_total_accuracy"`
    NgramTotalTime     float64 `json:"ngram_total_time"`
    NgramTotalSessions int     `json:"ngram_total_sessions"`
}
```

## N-gram Generation & Selection

### Adaptive Selection Algorithm

1. **Collect candidates**: All observed bigrams plus common trigrams
2. **Score each n-gram** using historical seek time data:
   - Score = average seek time (slower = higher priority)
   - N-grams with no data get medium priority (need baseline measurement)
   - Filter out n-grams with < 3 measurements (insufficient data)
3. **Weighted random selection**: Higher scores = more likely to appear
4. **Mix ratio**: ~70% bigrams, ~30% trigrams

### Trigram Tracking

Add trigram recording alongside existing bigram tracking:
- Record during both word mode and n-gram mode
- Trigram = 3 consecutive correct letters within a word/n-gram
- Same seek time filtering rules as bigrams (< 5000ms)

### Sequence Generation Function

```go
func GetTrainingSequence(
    targetChars int,
    bigramData map[string]BigramSeekStats,
    trigramData map[string]TrigramSeekStats,
) []string
```

- Returns slice of n-grams totalling ~150 characters (including spaces)
- Each n-gram appears as a separate "word" for the engine
- Example output: `["th", "er", "ing", "on", "the", "an", "ion", ...]`

## User Interface

### Options Screen (Ctrl+O)

Add new setting below "Perfect Mode":

```
Practice Mode
  ● Words (default)
  ○ N-grams
```

Navigation unchanged: Arrow keys, Enter/Space to select, number keys for quick select.

### Typing Screen Adaptations

| Element | Word Mode | N-gram Mode |
|---------|-----------|-------------|
| Header | "Typing Practice" | "N-gram Training" |
| Progress | "Word 12/30" | "N-gram 12/50" |
| Display | Block letter words | Block letter n-grams |
| Controls | Identical | Identical |

### Results Screen Adaptations

- Title: "N-gram Training Complete!"
- Stats section shows n-gram-specific bests/averages
- Additional insight: "Slowest n-grams this session:" (top 5)
- Letter matrix, finger/hand/row stats remain (still applicable)

### Web Frontend

- WelcomeScreen: Add mode toggle/selector
- TypingScreen: Same adaptations as TUI
- ResultsScreen: Same adaptations as TUI
- Settings stored in `baboon_settings` localStorage

## Data Flow & Persistence

### Settings Persistence

`~/.config/baboon/settings.json`:
```json
{
  "advance_key": "space",
  "perfect_mode": false,
  "practice_mode": "words"
}
```

### Statistics Flow

```
Round Start
    │
    ├─ Mode = Words? ──→ words.GetRandomWordsFixedCount()
    │
    └─ Mode = Ngrams? ─→ ngrams.GetTrainingSequence()
                              │
                              ├─ Reads: BigramSeekTime, TrigramSeekTime
                              └─ Returns: ["th", "er", "ing", ...]

During Typing (both modes)
    │
    └─ RecordBigramSeekTime()  ─→ updates shared BigramSeekTime
    └─ RecordTrigramSeekTime() ─→ updates shared TrigramSeekTime (NEW)

Round Complete
    │
    ├─ Mode = Words? ──→ UpdateHistorical() (word stats)
    │
    └─ Mode = Ngrams? ─→ UpdateNgramHistorical() (ngram stats)
```

### Shared vs Separate Data

| Data | Shared | Rationale |
|------|--------|-----------|
| BigramSeekTime | Yes | Informs adaptive selection in both modes |
| TrigramSeekTime | Yes | Informs adaptive selection in both modes |
| LetterAccuracy | Yes | Per-letter data useful across modes |
| LetterSeekTime | Yes | Per-letter timing useful across modes |
| Best WPM/Accuracy/Time | No | Separate for fair comparison |
| Total Sessions | No | Track mode usage independently |

## Testing

### Unit Tests (`ngrams/ngrams_test.go`)

- `TestGetTrainingSequence_ReturnsCorrectCharCount`: Verifies ~150 chars total
- `TestGetTrainingSequence_PrioritisesSlowNgrams`: Slowest n-grams appear more frequently
- `TestGetTrainingSequence_HandlesEmptyHistory`: Returns common n-grams when no data
- `TestGetTrainingSequence_MixesBigramsAndTrigrams`: Validates ~70/30 ratio

### Integration Tests

- `TestEngine_NgramMode_StartsRound`: Engine generates valid sequences
- `TestEngine_NgramMode_TracksTrigramSeekTime`: Trigrams recorded during typing
- `TestEngine_NgramMode_UpdatesNgramStats`: Separate stats updated on complete
- `TestEngine_NgramMode_SharedBigramData`: Bigram times contribute to shared pool

### Settings Tests

- `TestSettings_PracticeModePerisists`: Mode saved/loaded correctly
- `TestSettings_DefaultsToWords`: New users get word mode

### Manual Testing Checklist

- [ ] Options screen shows practice mode toggle
- [ ] N-gram mode displays sequences correctly in block letters
- [ ] Space advances between n-grams
- [ ] Results show n-gram-specific stats
- [ ] Switching modes preserves separate statistics
- [ ] Web frontend mirrors TUI behaviour

## Files to Create/Modify

### New Files
- `ngrams/ngrams.go` - N-gram generation and selection
- `ngrams/ngrams_test.go` - Unit tests

### Modified Files
- `backend/api.go` - Add `PracticeMode` to `Config`
- `backend/engine.go` - Branch on mode in `StartRound()`, add trigram tracking
- `stats/stats.go` - Add `TrigramSeekStats`, n-gram historical fields
- `settings/settings.go` - Add `PracticeMode` field
- `frontend/model.go` - Pass mode to engine, adapt header/progress text
- `frontend/views.go` - Adapt results screen for n-gram mode
- `web/src/components/TypingScreen.jsx` - Adapt for n-gram mode
- `web/src/components/ResultsScreen.jsx` - Show n-gram stats
- `web/src/components/OptionsScreen.jsx` - Add mode selector (or create if needed)
- `SPECIFICATION.md` - Document new feature

## User Stories

### US-013: Practice N-gram Typing
**As a** user wanting to improve specific letter combinations
**I want to** practice typing bigrams and trigrams in isolation
**So that** I can focus on my weakest transitions without full words

### US-014: Adaptive N-gram Selection
**As a** user practicing n-grams
**I want to** automatically practice my slowest combinations more often
**So that** I can efficiently improve my weak areas

## Functional Requirements

### FR-041: N-gram Practice Mode
- The application SHALL support an n-gram practice mode selectable via options screen
- N-gram mode SHALL present sequences of bigrams and trigrams separated by spaces
- Each round SHALL target 150 characters total (matching word mode)
- N-grams SHALL be selected adaptively based on user's historical seek times

### FR-042: Trigram Seek Time Tracking
- The application SHALL track seek time for letter triplets (trigrams)
- Trigrams SHALL be recorded for 3 consecutive correct keystrokes within a word
- Trigram data SHALL persist across sessions and inform adaptive selection

### FR-043: Separate N-gram Statistics
- N-gram mode SHALL track WPM, accuracy, and time separately from word mode
- Historical bests and averages SHALL be maintained independently
- Results screen SHALL display n-gram-specific statistics when in n-gram mode

### FR-044: N-gram Results Insights
- Results screen in n-gram mode SHALL display the 5 slowest n-grams from the session
- This helps users identify specific combinations needing more practice
