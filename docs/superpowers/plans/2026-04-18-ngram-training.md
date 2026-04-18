# N-gram Training Mode Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a practice mode for focused bigram and trigram training with adaptive selection based on user's slowest n-grams.

**Architecture:** Unified engine with mode flag. New `ngrams/` package handles n-gram generation and selection. Stats extended with trigram tracking and separate n-gram mode statistics. Settings extended with practice mode preference.

**Tech Stack:** Go, Bubble Tea (TUI), React/Chakra UI (Web), JSON persistence

---

## File Structure

### New Files
| File | Responsibility |
|------|----------------|
| `ngrams/ngrams.go` | N-gram generation, adaptive selection algorithm, common n-gram data |
| `ngrams/ngrams_test.go` | Unit tests for n-gram generation |

### Modified Files
| File | Changes |
|------|---------|
| `settings/settings.go` | Add `PracticeMode` enum and field |
| `stats/stats.go` | Add `TrigramSeekStats`, trigram tracking methods, n-gram historical fields |
| `backend/api.go` | Add `PracticeMode` to `Config` |
| `backend/engine.go` | Branch on mode in `StartRound()`, add trigram tracking |
| `frontend/model.go` | Pass mode to views, handle options screen changes |
| `frontend/views.go` | Adapt typing/results screens for n-gram mode |
| `web/src/components/WelcomeScreen.jsx` | Add mode selector |
| `web/src/components/TypingScreen.jsx` | Adapt progress/header for n-gram mode |
| `web/src/components/ResultsScreen.jsx` | Show n-gram stats, slowest n-grams insight |
| `SPECIFICATION.md` | Document FR-041 through FR-044 |

---

## Task 1: Add PracticeMode to Settings

**Files:**
- Modify: `settings/settings.go:10-60`
- Test: `settings/settings_test.go` (create)

- [ ] **Step 1: Write the failing test**

Create `settings/settings_test.go`:

```go
package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPracticeMode_String(t *testing.T) {
	tests := []struct {
		mode PracticeMode
		want string
	}{
		{ModeWords, "Words"},
		{ModeNgrams, "N-grams"},
	}
	for _, tt := range tests {
		if got := tt.mode.String(); got != tt.want {
			t.Errorf("PracticeMode.String() = %v, want %v", got, tt.want)
		}
	}
}

func TestSettings_PracticeModeDefaultsToWords(t *testing.T) {
	s := DefaultSettings()
	if s.PracticeMode != ModeWords {
		t.Errorf("DefaultSettings().PracticeMode = %v, want ModeWords", s.PracticeMode)
	}
}

func TestSettings_PracticeModePersists(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create config dir
	configDir := filepath.Join(tmpDir, ".config", "baboon")
	os.MkdirAll(configDir, 0755)

	// Save settings with N-gram mode
	s := DefaultSettings()
	s.PracticeMode = ModeNgrams
	if err := s.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load and verify
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.PracticeMode != ModeNgrams {
		t.Errorf("Loaded PracticeMode = %v, want ModeNgrams", loaded.PracticeMode)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./settings/... -v`
Expected: FAIL with "undefined: PracticeMode" or "undefined: ModeWords"

- [ ] **Step 3: Write minimal implementation**

Edit `settings/settings.go` to add after line 16 (after AdvanceKeyEither):

```go
// PracticeMode defines the type of typing practice
type PracticeMode int

const (
	ModeWords  PracticeMode = iota // Default: word-based practice
	ModeNgrams                     // N-gram (bigram/trigram) training
)

// String returns the display name for the practice mode
func (p PracticeMode) String() string {
	switch p {
	case ModeWords:
		return "Words"
	case ModeNgrams:
		return "N-grams"
	default:
		return "Words"
	}
}
```

Edit the `Settings` struct (around line 48) to add `PracticeMode`:

```go
// Settings holds user preferences
type Settings struct {
	AdvanceKey   AdvanceKey   `json:"advance_key"`
	PerfectMode  bool         `json:"perfect_mode"`
	PracticeMode PracticeMode `json:"practice_mode"`
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./settings/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add settings/settings.go settings/settings_test.go
git commit -m "feat: add PracticeMode setting for n-gram training"
```

---

## Task 2: Add TrigramSeekStats to Stats Package

**Files:**
- Modify: `stats/stats.go:59-72` and `stats/stats.go:197-220`
- Test: `stats/stats_test.go` (create or extend)

- [ ] **Step 1: Write the failing test**

Create `stats/trigram_test.go`:

```go
package stats

import "testing"

func TestTrigramSeekStats_AverageMs(t *testing.T) {
	tests := []struct {
		name  string
		stats TrigramSeekStats
		want  float64
	}{
		{"empty", TrigramSeekStats{}, 0},
		{"single", TrigramSeekStats{TotalTimeMs: 100, Count: 1}, 100},
		{"multiple", TrigramSeekStats{TotalTimeMs: 300, Count: 3}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stats.AverageMs(); got != tt.want {
				t.Errorf("AverageMs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStats_RecordTrigramSeekTime(t *testing.T) {
	s := &Stats{
		TrigramSeekTime: make(map[string]TrigramSeekStats),
	}

	s.RecordTrigramSeekTime("the", 150)
	s.RecordTrigramSeekTime("the", 200)
	s.RecordTrigramSeekTime("ing", 100)

	if got := s.TrigramSeekTime["the"].Count; got != 2 {
		t.Errorf("the count = %v, want 2", got)
	}
	if got := s.TrigramSeekTime["the"].TotalTimeMs; got != 350 {
		t.Errorf("the total = %v, want 350", got)
	}
	if got := s.TrigramSeekTime["ing"].Count; got != 1 {
		t.Errorf("ing count = %v, want 1", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./stats/... -v -run Trigram`
Expected: FAIL with "undefined: TrigramSeekStats"

- [ ] **Step 3: Write minimal implementation**

Add after `BigramSeekStats` (around line 63) in `stats/stats.go`:

```go
// TrigramSeekStats tracks seek time for letter triplets (e.g., "the", "ing", "tion")
type TrigramSeekStats struct {
	TotalTimeMs int64 `json:"total_time_ms"` // Total time in milliseconds
	Count       int   `json:"count"`         // Number of measurements
}

// AverageMs returns the average seek time in milliseconds for the trigram
func (s TrigramSeekStats) AverageMs() float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.TotalTimeMs) / float64(s.Count)
}
```

Add `TrigramSeekTime` field to `Stats` struct (around line 23):

```go
TrigramSeekTime map[string]TrigramSeekStats `json:"-"` // Per-trigram seek time for this session
```

Add `RecordTrigramSeekTime` method after `RecordBigramSeekTime` (around line 262):

```go
// RecordTrigramSeekTime records the time taken to type a letter triplet (trigram)
func (s *Stats) RecordTrigramSeekTime(trigram string, durationMs int64) {
	if s.TrigramSeekTime == nil {
		s.TrigramSeekTime = make(map[string]TrigramSeekStats)
	}
	stats := s.TrigramSeekTime[trigram]
	stats.TotalTimeMs += durationMs
	stats.Count++
	s.TrigramSeekTime[trigram] = stats
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./stats/... -v -run Trigram`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add stats/stats.go stats/trigram_test.go
git commit -m "feat: add TrigramSeekStats for trigram seek time tracking"
```

---

## Task 3: Add N-gram Historical Stats Fields

**Files:**
- Modify: `stats/stats.go:197-220` (HistoricalStats struct)
- Modify: `stats/stats.go:422-477` (LoadHistoricalStats)
- Modify: `stats/stats.go:534-654` (UpdateHistorical)

- [ ] **Step 1: Write the failing test**

Add to `stats/trigram_test.go`:

```go
func TestHistoricalStats_NgramAverages(t *testing.T) {
	h := &HistoricalStats{
		NgramTotalWPM:      200,
		NgramTotalAccuracy: 190,
		NgramTotalTime:     300,
		NgramTotalSessions: 2,
	}

	if got := h.NgramAverageWPM(); got != 100 {
		t.Errorf("NgramAverageWPM() = %v, want 100", got)
	}
	if got := h.NgramAverageAccuracy(); got != 95 {
		t.Errorf("NgramAverageAccuracy() = %v, want 95", got)
	}
	if got := h.NgramAverageTime(); got != 150 {
		t.Errorf("NgramAverageTime() = %v, want 150", got)
	}
}

func TestHistoricalStats_UpdateNgramHistorical(t *testing.T) {
	h := &HistoricalStats{
		TrigramSeekTime: make(map[string]TrigramSeekStats),
	}
	session := &Stats{
		WPM:      80,
		Accuracy: 95,
		Duration: 60 * 1000000000, // 60 seconds
		TrigramSeekTime: map[string]TrigramSeekStats{
			"the": {TotalTimeMs: 200, Count: 2},
		},
	}

	h.UpdateNgramHistorical(session)

	if h.NgramTotalSessions != 1 {
		t.Errorf("NgramTotalSessions = %v, want 1", h.NgramTotalSessions)
	}
	if h.NgramBestWPM != 80 {
		t.Errorf("NgramBestWPM = %v, want 80", h.NgramBestWPM)
	}
	if h.TrigramSeekTime["the"].Count != 2 {
		t.Errorf("TrigramSeekTime[the].Count = %v, want 2", h.TrigramSeekTime["the"].Count)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./stats/... -v -run Ngram`
Expected: FAIL with "undefined: NgramTotalWPM" or similar

- [ ] **Step 3: Write minimal implementation**

Add to `HistoricalStats` struct (around line 209):

```go
	// Trigram tracking (shared across modes)
	TrigramSeekTime map[string]TrigramSeekStats `json:"trigram_seek_time"`

	// N-gram mode separate tracking
	NgramBestWPM       float64 `json:"ngram_best_wpm"`
	NgramBestAccuracy  float64 `json:"ngram_best_accuracy"`
	NgramBestTime      float64 `json:"ngram_best_time"`
	NgramTotalWPM      float64 `json:"ngram_total_wpm"`
	NgramTotalAccuracy float64 `json:"ngram_total_accuracy"`
	NgramTotalTime     float64 `json:"ngram_total_time"`
	NgramTotalSessions int     `json:"ngram_total_sessions"`
```

Add average methods after `AverageTime()` (around line 678):

```go
// NgramAverageWPM returns the average WPM across all n-gram sessions
func (h *HistoricalStats) NgramAverageWPM() float64 {
	if h.NgramTotalSessions == 0 {
		return 0
	}
	return h.NgramTotalWPM / float64(h.NgramTotalSessions)
}

// NgramAverageAccuracy returns the average accuracy across all n-gram sessions
func (h *HistoricalStats) NgramAverageAccuracy() float64 {
	if h.NgramTotalSessions == 0 {
		return 0
	}
	return h.NgramTotalAccuracy / float64(h.NgramTotalSessions)
}

// NgramAverageTime returns the average time across all n-gram sessions in seconds
func (h *HistoricalStats) NgramAverageTime() float64 {
	if h.NgramTotalSessions == 0 {
		return 0
	}
	return h.NgramTotalTime / float64(h.NgramTotalSessions)
}

// UpdateNgramHistorical updates n-gram-specific stats with new session data
func (h *HistoricalStats) UpdateNgramHistorical(session *Stats) {
	h.NgramTotalSessions++

	// Update totals for averages
	h.NgramTotalWPM += session.WPM
	h.NgramTotalAccuracy += session.Accuracy
	h.NgramTotalTime += session.Duration.Seconds()

	// Update bests
	if session.WPM > h.NgramBestWPM {
		h.NgramBestWPM = session.WPM
	}
	if session.Accuracy > h.NgramBestAccuracy {
		h.NgramBestAccuracy = session.Accuracy
	}
	if h.NgramBestTime == 0 || session.Duration.Seconds() < h.NgramBestTime {
		h.NgramBestTime = session.Duration.Seconds()
	}

	// Merge trigram seek time (shared across modes)
	if h.TrigramSeekTime == nil {
		h.TrigramSeekTime = make(map[string]TrigramSeekStats)
	}
	for trigram, sessionStats := range session.TrigramSeekTime {
		histStats := h.TrigramSeekTime[trigram]
		histStats.TotalTimeMs += sessionStats.TotalTimeMs
		histStats.Count += sessionStats.Count
		h.TrigramSeekTime[trigram] = histStats
	}
}
```

Update `LoadHistoricalStats()` to initialize `TrigramSeekTime` (add after BigramSeekTime init around line 459):

```go
	if stats.TrigramSeekTime == nil {
		stats.TrigramSeekTime = make(map[string]TrigramSeekStats)
	}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./stats/... -v -run Ngram`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add stats/stats.go stats/trigram_test.go
git commit -m "feat: add n-gram historical stats fields and UpdateNgramHistorical"
```

---

## Task 4: Create N-gram Generation Package

**Files:**
- Create: `ngrams/ngrams.go`
- Create: `ngrams/ngrams_test.go`

- [ ] **Step 1: Write the failing test**

Create `ngrams/ngrams_test.go`:

```go
package ngrams

import (
	"testing"

	"github.com/timlinux/baboon/stats"
)

func TestGetTrainingSequence_ReturnsCorrectCharCount(t *testing.T) {
	bigramData := map[string]stats.BigramSeekStats{}
	trigramData := map[string]stats.TrigramSeekStats{}

	result := GetTrainingSequence(150, bigramData, trigramData, nil)

	// Count total characters (including spaces between n-grams)
	totalChars := 0
	for i, ngram := range result {
		totalChars += len(ngram)
		if i < len(result)-1 {
			totalChars++ // space between n-grams
		}
	}

	// Should be within 10 chars of target (allow for n-gram boundaries)
	if totalChars < 140 || totalChars > 160 {
		t.Errorf("Total chars = %d, want 140-160", totalChars)
	}
}

func TestGetTrainingSequence_HandlesEmptyHistory(t *testing.T) {
	result := GetTrainingSequence(150, nil, nil, nil)

	if len(result) == 0 {
		t.Error("Expected non-empty result with empty history")
	}

	// Verify all results are valid n-grams (2-3 lowercase letters)
	for _, ngram := range result {
		if len(ngram) < 2 || len(ngram) > 3 {
			t.Errorf("Invalid n-gram length: %q", ngram)
		}
		for _, c := range ngram {
			if c < 'a' || c > 'z' {
				t.Errorf("Invalid character in n-gram: %q", ngram)
			}
		}
	}
}

func TestGetTrainingSequence_PrioritisesSlowNgrams(t *testing.T) {
	// Create data where "zz" is very slow
	bigramData := map[string]stats.BigramSeekStats{
		"zz": {TotalTimeMs: 1000, Count: 5}, // 200ms average - slow
		"th": {TotalTimeMs: 250, Count: 5},  // 50ms average - fast
		"er": {TotalTimeMs: 250, Count: 5},  // 50ms average - fast
	}

	// Generate multiple sequences and count "zz" occurrences
	zzCount := 0
	thCount := 0
	for i := 0; i < 100; i++ {
		result := GetTrainingSequence(150, bigramData, nil, func(n int) int { return i % n })
		for _, ngram := range result {
			if ngram == "zz" {
				zzCount++
			}
			if ngram == "th" {
				thCount++
			}
		}
	}

	// "zz" should appear more often than "th" due to higher priority
	if zzCount <= thCount {
		t.Errorf("Slow n-gram 'zz' count (%d) should exceed fast 'th' count (%d)", zzCount, thCount)
	}
}

func TestGetTrainingSequence_MixesBigramsAndTrigrams(t *testing.T) {
	result := GetTrainingSequence(150, nil, nil, nil)

	bigramCount := 0
	trigramCount := 0
	for _, ngram := range result {
		if len(ngram) == 2 {
			bigramCount++
		} else if len(ngram) == 3 {
			trigramCount++
		}
	}

	// Should have mix of both
	if bigramCount == 0 {
		t.Error("Expected some bigrams in result")
	}
	if trigramCount == 0 {
		t.Error("Expected some trigrams in result")
	}

	// Trigrams should be minority (~30%)
	total := bigramCount + trigramCount
	trigramRatio := float64(trigramCount) / float64(total)
	if trigramRatio > 0.5 {
		t.Errorf("Trigram ratio %.2f is too high, expected ~0.3", trigramRatio)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./ngrams/... -v`
Expected: FAIL with "no Go files in ngrams" or "undefined: GetTrainingSequence"

- [ ] **Step 3: Write minimal implementation**

Create `ngrams/ngrams.go`:

```go
// Package ngrams provides n-gram generation for typing practice.
package ngrams

import (
	"sort"

	"github.com/timlinux/baboon/stats"
)

// CommonBigrams contains frequently occurring English bigrams
var CommonBigrams = []string{
	"th", "he", "in", "er", "an", "re", "on", "at", "en", "nd",
	"ti", "es", "or", "te", "of", "ed", "is", "it", "al", "ar",
	"st", "to", "nt", "ng", "se", "ha", "as", "ou", "io", "le",
	"ve", "co", "me", "de", "hi", "ri", "ro", "ic", "ne", "ea",
	"ra", "ce", "li", "ch", "ll", "be", "ma", "si", "om", "ur",
}

// CommonTrigrams contains frequently occurring English trigrams
var CommonTrigrams = []string{
	"the", "and", "ing", "ion", "tio", "ent", "ati", "for", "her", "ter",
	"hat", "tha", "ere", "ate", "his", "con", "res", "ver", "all", "ons",
	"nce", "men", "ith", "ted", "ers", "pro", "thi", "wit", "are", "ess",
	"not", "ive", "was", "ect", "rea", "com", "eve", "per", "int", "est",
	"sta", "cti", "ica", "ist", "ear", "ain", "one", "our", "iti", "rat",
}

// scoredNgram holds an n-gram with its selection score
type scoredNgram struct {
	ngram string
	score float64
}

// GetTrainingSequence generates a sequence of n-grams for training.
// It prioritises n-grams where the user has slower seek times.
// Returns a slice of n-grams that together total approximately targetChars characters.
func GetTrainingSequence(
	targetChars int,
	bigramData map[string]stats.BigramSeekStats,
	trigramData map[string]stats.TrigramSeekStats,
	randFunc func(n int) int,
) []string {
	if randFunc == nil {
		randFunc = defaultRand
	}

	// Build scored list of candidates
	candidates := buildCandidates(bigramData, trigramData)

	// Generate sequence
	result := []string{}
	currentChars := 0

	for currentChars < targetChars {
		// Decide bigram vs trigram (~70/30 split)
		useTrigram := randFunc(10) < 3

		// Select from appropriate candidates
		var selected string
		if useTrigram {
			selected = selectWeighted(filterByLength(candidates, 3), randFunc)
		} else {
			selected = selectWeighted(filterByLength(candidates, 2), randFunc)
		}

		// Fallback if no candidates of desired length
		if selected == "" {
			if useTrigram && len(CommonTrigrams) > 0 {
				selected = CommonTrigrams[randFunc(len(CommonTrigrams))]
			} else if len(CommonBigrams) > 0 {
				selected = CommonBigrams[randFunc(len(CommonBigrams))]
			}
		}

		if selected == "" {
			break
		}

		// Check if adding this would exceed target (with space)
		addedChars := len(selected)
		if len(result) > 0 {
			addedChars++ // space
		}
		if currentChars+addedChars > targetChars+5 {
			break
		}

		result = append(result, selected)
		currentChars += addedChars
	}

	return result
}

// buildCandidates creates a scored list of n-gram candidates
func buildCandidates(
	bigramData map[string]stats.BigramSeekStats,
	trigramData map[string]stats.TrigramSeekStats,
) []scoredNgram {
	candidates := []scoredNgram{}

	// Add bigrams from history
	for ngram, data := range bigramData {
		if data.Count >= 3 {
			candidates = append(candidates, scoredNgram{
				ngram: ngram,
				score: data.AverageMs(),
			})
		}
	}

	// Add trigrams from history
	for ngram, data := range trigramData {
		if data.Count >= 3 {
			candidates = append(candidates, scoredNgram{
				ngram: ngram,
				score: data.AverageMs(),
			})
		}
	}

	// Add common n-grams with medium priority if not already present
	existingNgrams := make(map[string]bool)
	for _, c := range candidates {
		existingNgrams[c.ngram] = true
	}

	mediumScore := 150.0 // Default score for unknown n-grams
	for _, ngram := range CommonBigrams {
		if !existingNgrams[ngram] {
			candidates = append(candidates, scoredNgram{ngram: ngram, score: mediumScore})
		}
	}
	for _, ngram := range CommonTrigrams {
		if !existingNgrams[ngram] {
			candidates = append(candidates, scoredNgram{ngram: ngram, score: mediumScore})
		}
	}

	// Sort by score descending (slower = higher priority)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	return candidates
}

// filterByLength returns only n-grams of the specified length
func filterByLength(candidates []scoredNgram, length int) []scoredNgram {
	filtered := []scoredNgram{}
	for _, c := range candidates {
		if len(c.ngram) == length {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// selectWeighted selects an n-gram with probability proportional to score
func selectWeighted(candidates []scoredNgram, randFunc func(n int) int) string {
	if len(candidates) == 0 {
		return ""
	}

	// Calculate total score
	totalScore := 0.0
	for _, c := range candidates {
		totalScore += c.score
	}

	if totalScore == 0 {
		return candidates[randFunc(len(candidates))].ngram
	}

	// Select weighted random
	target := float64(randFunc(int(totalScore*100))) / 100.0
	cumulative := 0.0
	for _, c := range candidates {
		cumulative += c.score
		if cumulative >= target {
			return c.ngram
		}
	}

	return candidates[len(candidates)-1].ngram
}

// defaultRand is a simple deterministic pseudo-random for testing
var randState int = 42

func defaultRand(n int) int {
	if n <= 0 {
		return 0
	}
	randState = (randState*1103515245 + 12345) & 0x7fffffff
	return randState % n
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./ngrams/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add ngrams/ngrams.go ngrams/ngrams_test.go
git commit -m "feat: add ngrams package with adaptive n-gram generation"
```

---

## Task 5: Add PracticeMode to Backend Config

**Files:**
- Modify: `backend/api.go:139-201`
- Modify: `settings/settings.go` (import)

- [ ] **Step 1: Write the failing test**

Create `backend/api_test.go`:

```go
package backend

import (
	"testing"

	"github.com/timlinux/baboon/settings"
)

func TestDefaultConfig_PracticeModeIsWords(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.PracticeMode != settings.ModeWords {
		t.Errorf("DefaultConfig().PracticeMode = %v, want ModeWords", cfg.PracticeMode)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/... -v -run TestDefaultConfig_PracticeMode`
Expected: FAIL with "cfg.PracticeMode undefined"

- [ ] **Step 3: Write minimal implementation**

Edit `backend/api.go`:

Add import at top:
```go
import (
	"time"

	"github.com/timlinux/baboon/settings"
	"github.com/timlinux/baboon/stats"
)
```

Add to `Config` struct (around line 142):
```go
	// PracticeMode selects between word and n-gram training
	PracticeMode settings.PracticeMode
```

Update `DefaultConfig()` (around line 197):
```go
func DefaultConfig() Config {
	return Config{
		PunctuationMode:    false,
		WordsPerRound:      30,
		CharactersPerRound: 150,
		PracticeMode:       settings.ModeWords,
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/... -v -run TestDefaultConfig_PracticeMode`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/api.go backend/api_test.go
git commit -m "feat: add PracticeMode to backend Config"
```

---

## Task 6: Engine Branches on PracticeMode in StartRound

**Files:**
- Modify: `backend/engine.go:52-108`

- [ ] **Step 1: Write the failing test**

Add to `backend/engine_test.go` (create if needed):

```go
package backend

import (
	"testing"

	"github.com/timlinux/baboon/settings"
)

func TestEngine_NgramMode_StartsRound(t *testing.T) {
	cfg := DefaultConfig()
	cfg.PracticeMode = settings.ModeNgrams

	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	state := engine.GetGameState()

	// Should have n-grams (2-3 letter sequences)
	if len(state.Words) == 0 {
		t.Error("Expected non-empty word list")
	}

	for _, ngram := range state.Words {
		if len(ngram) < 2 || len(ngram) > 3 {
			t.Errorf("Expected n-gram (2-3 chars), got %q", ngram)
		}
	}
}

func TestEngine_WordMode_StartsRound(t *testing.T) {
	cfg := DefaultConfig()
	cfg.PracticeMode = settings.ModeWords

	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	state := engine.GetGameState()

	// Should have words (variable length)
	if len(state.Words) != 30 {
		t.Errorf("Expected 30 words, got %d", len(state.Words))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/... -v -run TestEngine_NgramMode`
Expected: FAIL (n-grams not generated, still using words)

- [ ] **Step 3: Write minimal implementation**

Edit `backend/engine.go`:

Add import:
```go
import (
	"fmt"
	"math/rand"
	"time"

	"github.com/timlinux/baboon/ngrams"
	"github.com/timlinux/baboon/settings"
	"github.com/timlinux/baboon/stats"
	"github.com/timlinux/baboon/words"
)
```

Modify `StartRound()` to branch on mode (around line 53-62):

```go
// StartRound initialises a new round with fresh words and resets session stats.
func (e *Engine) StartRound() {
	// Get letter data for weighted word selection
	letterData := e.getLetterData()

	// Generate words or n-grams based on practice mode
	if e.config.PracticeMode == settings.ModeNgrams {
		e.words = ngrams.GetTrainingSequence(
			e.config.CharactersPerRound,
			e.historical.BigramSeekTime,
			e.historical.TrigramSeekTime,
			e.rng.Intn,
		)
	} else {
		e.words = words.GetRandomWordsFixedCount(
			e.config.WordsPerRound,
			e.config.CharactersPerRound,
			e.rng.Intn,
			letterData,
		)
	}

	// Create new session stats
	e.session = &stats.Stats{
		LetterAccuracy:    make(map[string]stats.LetterStats),
		LetterSeekTime:    make(map[string]stats.LetterSeekStats),
		BigramSeekTime:    make(map[string]stats.BigramSeekStats),
		TrigramSeekTime:   make(map[string]stats.TrigramSeekStats),
		FingerStats:       make(map[int]stats.FingerStat),
		HandStats:         make(map[int]stats.HandStat),
		RowStats:          make(map[int]stats.RowStat),
		ErrorSubstitution: make(map[string]map[string]int),
		SeekTimes:         make([]int64, 0),
	}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/... -v -run TestEngine`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/engine.go backend/engine_test.go
git commit -m "feat: engine branches on PracticeMode in StartRound"
```

---

## Task 7: Add Trigram Tracking to Engine

**Files:**
- Modify: `backend/engine.go:118-251` (ProcessKeystrokeWithTiming)

- [ ] **Step 1: Write the failing test**

Add to `backend/engine_test.go`:

```go
func TestEngine_TracksTrigramSeekTime(t *testing.T) {
	cfg := DefaultConfig()
	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	// Type "the" with timing
	engine.ProcessKeystrokeWithTiming("t", 0)   // First letter, no seek time
	engine.ProcessKeystrokeWithTiming("h", 100) // Second letter
	engine.ProcessKeystrokeWithTiming("e", 100) // Third letter - should record trigram

	session := engine.GetSessionStats()

	// Check trigram was recorded
	if _, exists := session.TrigramSeekTime["the"]; !exists {
		t.Error("Expected trigram 'the' to be tracked")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/... -v -run TestEngine_TracksTrigramSeekTime`
Expected: FAIL (trigram not recorded)

- [ ] **Step 3: Write minimal implementation**

Edit `backend/engine.go` in `ProcessKeystrokeWithTiming` function.

Add field to Engine struct (around line 28):
```go
	// Track last two letters for trigram detection
	lastLetter       string
	secondLastLetter string
```

In `ProcessKeystrokeWithTiming`, after bigram recording (around line 190), add:

```go
				// Record trigram timing (3 consecutive correct letters)
				if e.secondLastLetter != "" && e.lastLetter != "" {
					trigram := e.secondLastLetter + e.lastLetter + expectedLetter
					e.session.RecordTrigramSeekTime(trigram, seekTimeMs)
				}

				// Update letter history for next bigram/trigram
				e.secondLastLetter = e.lastLetter
				e.lastLetter = expectedLetter
```

Remove the old `e.lastLetter = expectedLetter` line that was at the end of the isLetter block.

In `StartRound()`, reset secondLastLetter (around line 107):
```go
	e.lastLetter = ""
	e.secondLastLetter = ""
```

In `ProcessSpaceWithTiming`, reset both (around line 289):
```go
		e.lastLetter = ""       // Reset for new word
		e.secondLastLetter = "" // Reset for new word
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/... -v -run TestEngine_TracksTrigramSeekTime`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/engine.go backend/engine_test.go
git commit -m "feat: add trigram seek time tracking to engine"
```

---

## Task 8: Engine Uses Correct Stats Update Method

**Files:**
- Modify: `backend/engine.go:309-327` (SubmitTiming)

- [ ] **Step 1: Write the failing test**

Add to `backend/engine_test.go`:

```go
func TestEngine_NgramMode_UpdatesNgramStats(t *testing.T) {
	cfg := DefaultConfig()
	cfg.PracticeMode = settings.ModeNgrams

	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	// Complete a round
	state := engine.GetGameState()
	for i, ngram := range state.Words {
		for _, char := range ngram {
			engine.ProcessKeystrokeWithTiming(string(char), 50)
		}
		if i < len(state.Words)-1 {
			engine.ProcessSpaceWithTiming(50)
		}
	}

	// Submit timing
	now := time.Now()
	engine.SubmitTiming(now.Add(-time.Minute), now, 60000)

	historical := engine.GetHistoricalStats()
	if historical.NgramTotalSessions != 1 {
		t.Errorf("NgramTotalSessions = %d, want 1", historical.NgramTotalSessions)
	}
	if historical.TotalSessions != 0 {
		t.Errorf("TotalSessions (word mode) = %d, want 0", historical.TotalSessions)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/... -v -run TestEngine_NgramMode_UpdatesNgramStats`
Expected: FAIL (NgramTotalSessions is 0)

- [ ] **Step 3: Write minimal implementation**

Edit `backend/engine.go` in `SubmitTiming` (around line 327):

```go
// SubmitTiming receives final timing data from the frontend and calculates stats.
func (e *Engine) SubmitTiming(startTime, endTime time.Time, durationMs int64) {
	e.session.StartTime = startTime
	e.session.EndTime = endTime
	e.session.Duration = time.Duration(durationMs) * time.Millisecond

	// Calculate WPM and accuracy using the frontend-provided timing
	minutes := e.session.Duration.Minutes()
	if minutes > 0 {
		e.session.WPM = (float64(e.session.CorrectChars) / 5.0) / minutes
	}
	if e.session.TotalCharacters > 0 {
		e.session.Accuracy = (float64(e.session.CorrectChars) / float64(e.session.TotalCharacters)) * 100
	}

	// Update historical stats based on practice mode
	if e.config.PracticeMode == settings.ModeNgrams {
		e.historical.UpdateNgramHistorical(e.session)
	} else {
		e.historical.UpdateHistorical(e.session)
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/... -v -run TestEngine_NgramMode_UpdatesNgramStats`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/engine.go backend/engine_test.go
git commit -m "feat: engine uses UpdateNgramHistorical for n-gram mode"
```

---

## Task 9: Add Practice Mode to Options Screen (TUI)

**Files:**
- Modify: `frontend/views.go` (RenderOptionsScreen)
- Modify: `frontend/model.go` (options handling)

- [ ] **Step 1: Read current options screen implementation**

Run: `grep -n "RenderOptionsScreen\|optionsCursor" frontend/views.go frontend/model.go`

- [ ] **Step 2: Modify RenderOptionsScreen to include Practice Mode**

Edit `frontend/views.go` - find `RenderOptionsScreen` and add Practice Mode section after Perfect Mode:

```go
// Add after Perfect Mode rendering, before the closing help text:

	// Practice Mode section
	practiceModeLabel := "Practice Mode"
	practiceModeLabelStyle := r.styles.ResultLabel
	if cursor == 2 { // Adjust cursor index based on existing options
		practiceModeLabelStyle = r.styles.ResultLabel.Background(lipgloss.Color("236"))
	}
	lines = append(lines, practiceModeLabelStyle.Render(practiceModeLabel))

	wordsOption := "  Words"
	ngramsOption := "  N-grams"

	if s.PracticeMode == settings.ModeWords {
		wordsOption = "● Words"
	}
	if s.PracticeMode == settings.ModeNgrams {
		ngramsOption = "● N-grams"
	}

	wordsStyle := r.styles.ResultLabel
	ngramsStyle := r.styles.ResultLabel
	if cursor == 3 {
		wordsStyle = r.styles.ResultLabel.Background(lipgloss.Color("236"))
	}
	if cursor == 4 {
		ngramsStyle = r.styles.ResultLabel.Background(lipgloss.Color("236"))
	}

	lines = append(lines, wordsStyle.Render(wordsOption))
	lines = append(lines, ngramsStyle.Render(ngramsOption))
```

- [ ] **Step 3: Update model to handle Practice Mode selection**

Edit `frontend/model.go` in the options screen key handling section.

Update the max cursor value and add handling for Practice Mode selection:

```go
// In handleOptionsKey function, update cursor bounds and add:
case "3":
	m.settings.PracticeMode = settings.ModeWords
	m.settings.Save()
case "4":
	m.settings.PracticeMode = settings.ModeNgrams
	m.settings.Save()
```

Also handle Enter/Space on cursor positions 3 and 4:
```go
case 3:
	m.settings.PracticeMode = settings.ModeWords
	m.settings.Save()
case 4:
	m.settings.PracticeMode = settings.ModeNgrams
	m.settings.Save()
```

- [ ] **Step 4: Test manually**

Run: `go build && ./baboon`
Press Ctrl+O to open options, verify Practice Mode appears and can be selected.

- [ ] **Step 5: Commit**

```bash
git add frontend/views.go frontend/model.go
git commit -m "feat: add Practice Mode to TUI options screen"
```

---

## Task 10: Adapt TUI Header and Progress for N-gram Mode

**Files:**
- Modify: `frontend/views.go:36-95` (RenderTypingScreenAnimated)

- [ ] **Step 1: Update header to show mode**

Edit `frontend/views.go` in `RenderTypingScreenAnimated`:

Add parameter to receive practice mode:
```go
func (r *Renderer) RenderTypingScreenAnimated(state backend.GameState, carousel *CarouselAnimator, s *settings.Settings, practiceMode settings.PracticeMode) string {
```

Update progress indicator (around line 84):
```go
	// Progress indicator - adapt text based on mode
	var progress string
	if practiceMode == settings.ModeNgrams {
		progress = fmt.Sprintf("N-gram %d/%d", state.WordNumber, state.TotalWords)
	} else {
		progress = fmt.Sprintf("Word %d/%d", state.WordNumber, state.TotalWords)
	}
```

- [ ] **Step 2: Update header in RenderFullScreen**

Find where header is rendered and update:
```go
	// Header - show mode
	var headerText string
	if practiceMode == settings.ModeNgrams {
		headerText = "🐒 BABOON - N-gram Training"
	} else {
		headerText = "🐒 BABOON - Typing Practice"
	}
```

- [ ] **Step 3: Update caller in model.go**

Edit `frontend/model.go` View() function to pass practice mode:
```go
content = m.renderer.RenderTypingScreenAnimated(state, m.carouselAnimator, m.settings, m.settings.PracticeMode)
```

- [ ] **Step 4: Test manually**

Run: `go build && ./baboon`
Enable N-gram mode in options, start typing, verify header shows "N-gram Training" and progress shows "N-gram X/Y".

- [ ] **Step 5: Commit**

```bash
git add frontend/views.go frontend/model.go
git commit -m "feat: adapt TUI header and progress for n-gram mode"
```

---

## Task 11: Adapt Results Screen for N-gram Mode

**Files:**
- Modify: `frontend/views.go` (RenderResultsScreen)

- [ ] **Step 1: Update results title and stats**

Edit `frontend/views.go` in results rendering:

```go
// Update title based on mode
var title string
if practiceMode == settings.ModeNgrams {
	title = r.styles.Title.Render("N-gram Training Complete!")
} else {
	title = r.styles.Title.Render("Round Complete!")
}

// Use n-gram stats when in n-gram mode
var bestWPM, avgWPM, bestAcc, avgAcc, bestTime, avgTime float64
var totalSessions int

if practiceMode == settings.ModeNgrams {
	bestWPM = historical.NgramBestWPM
	avgWPM = historical.NgramAverageWPM()
	bestAcc = historical.NgramBestAccuracy
	avgAcc = historical.NgramAverageAccuracy()
	bestTime = historical.NgramBestTime
	avgTime = historical.NgramAverageTime()
	totalSessions = historical.NgramTotalSessions
} else {
	bestWPM = historical.BestWPM
	avgWPM = historical.AverageWPM()
	bestAcc = historical.BestAccuracy
	avgAcc = historical.AverageAccuracy()
	bestTime = historical.BestTime
	avgTime = historical.AverageTime()
	totalSessions = historical.TotalSessions
}
```

- [ ] **Step 2: Add slowest n-grams insight for n-gram mode**

Add after other stats in results:

```go
// Show slowest n-grams from this session (n-gram mode only)
if practiceMode == settings.ModeNgrams && len(session.BigramSeekTime) > 0 {
	lines = append(lines, "")
	lines = append(lines, r.styles.ResultLabel.Render("Slowest n-grams:"))

	// Collect and sort by average time
	type ngramTime struct {
		ngram string
		avg   float64
	}
	var slowest []ngramTime
	for ngram, data := range session.BigramSeekTime {
		slowest = append(slowest, ngramTime{ngram, data.AverageMs()})
	}
	for ngram, data := range session.TrigramSeekTime {
		slowest = append(slowest, ngramTime{ngram, data.AverageMs()})
	}
	sort.Slice(slowest, func(i, j int) bool {
		return slowest[i].avg > slowest[j].avg
	})

	// Show top 5
	for i := 0; i < 5 && i < len(slowest); i++ {
		lines = append(lines, fmt.Sprintf("  %s: %.0fms", slowest[i].ngram, slowest[i].avg))
	}
}
```

- [ ] **Step 3: Update function signature and caller**

Add `practiceMode` parameter to results rendering function and update caller.

- [ ] **Step 4: Test manually**

Run n-gram mode and complete a round. Verify results show n-gram-specific stats.

- [ ] **Step 5: Commit**

```bash
git add frontend/views.go
git commit -m "feat: adapt results screen for n-gram mode stats"
```

---

## Task 12: Update SPECIFICATION.md

**Files:**
- Modify: `SPECIFICATION.md`

- [ ] **Step 1: Add user stories**

Add after US-012:

```markdown
### US-013: Practice N-gram Typing
**As a** user wanting to improve specific letter combinations
**I want to** practice typing bigrams and trigrams in isolation
**So that** I can focus on my weakest transitions without full words

### US-014: Adaptive N-gram Selection
**As a** user practicing n-grams
**I want to** automatically practice my slowest combinations more often
**So that** I can efficiently improve my weak areas
```

- [ ] **Step 2: Add functional requirements**

Add after FR-040:

```markdown
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
```

- [ ] **Step 3: Commit**

```bash
git add SPECIFICATION.md
git commit -m "docs: add n-gram training user stories and functional requirements"
```

---

## Task 13: Web Frontend - Add Mode to Settings

**Files:**
- Modify: `web/src/App.js` (add state)
- Modify: `web/src/components/WelcomeScreen.jsx` (add selector)

- [ ] **Step 1: Add practice mode state to App.js**

```jsx
const [practiceMode, setPracticeMode] = useState(() => {
  const saved = localStorage.getItem('baboon_settings');
  if (saved) {
    const settings = JSON.parse(saved);
    return settings.practice_mode || 'words';
  }
  return 'words';
});

// Save when changed
useEffect(() => {
  const settings = JSON.parse(localStorage.getItem('baboon_settings') || '{}');
  settings.practice_mode = practiceMode;
  localStorage.setItem('baboon_settings', JSON.stringify(settings));
}, [practiceMode]);
```

- [ ] **Step 2: Add mode toggle to WelcomeScreen**

```jsx
<Box mt={4}>
  <Text fontSize="sm" color="gray.400" mb={2}>Practice Mode</Text>
  <HStack spacing={4}>
    <Button
      variant={practiceMode === 'words' ? 'solid' : 'outline'}
      colorScheme="orange"
      onClick={() => setPracticeMode('words')}
    >
      Words
    </Button>
    <Button
      variant={practiceMode === 'ngrams' ? 'solid' : 'outline'}
      colorScheme="orange"
      onClick={() => setPracticeMode('ngrams')}
    >
      N-grams
    </Button>
  </HStack>
</Box>
```

- [ ] **Step 3: Pass mode to child components**

Update App.js to pass `practiceMode` to TypingScreen and ResultsScreen.

- [ ] **Step 4: Test in browser**

Run: `cd web && npm start`
Verify mode toggle appears and persists.

- [ ] **Step 5: Commit**

```bash
git add web/src/App.js web/src/components/WelcomeScreen.jsx
git commit -m "feat(web): add practice mode toggle to welcome screen"
```

---

## Task 14: Web Frontend - Adapt TypingScreen for N-grams

**Files:**
- Modify: `web/src/components/TypingScreen.jsx`

- [ ] **Step 1: Update header and progress text**

```jsx
// Receive practiceMode prop
const TypingScreen = ({ practiceMode, ... }) => {

  // Dynamic header
  const headerText = practiceMode === 'ngrams'
    ? '🐒 BABOON - N-gram Training'
    : '🐒 BABOON - Typing Practice';

  // Dynamic progress
  const progressText = practiceMode === 'ngrams'
    ? `N-gram ${wordNumber}/${totalWords}`
    : `Word ${wordNumber}/${totalWords}`;
```

- [ ] **Step 2: Update session creation to pass mode**

When creating session, include practice mode in request if backend supports it.

- [ ] **Step 3: Test in browser**

Verify header and progress adapt to selected mode.

- [ ] **Step 4: Commit**

```bash
git add web/src/components/TypingScreen.jsx
git commit -m "feat(web): adapt TypingScreen header and progress for n-gram mode"
```

---

## Task 15: Web Frontend - Adapt ResultsScreen for N-grams

**Files:**
- Modify: `web/src/components/ResultsScreen.jsx`

- [ ] **Step 1: Update title and stats**

```jsx
const ResultsScreen = ({ practiceMode, sessionStats, historicalStats, ... }) => {

  const title = practiceMode === 'ngrams'
    ? 'N-gram Training Complete!'
    : 'Round Complete!';

  // Use n-gram stats when in n-gram mode
  const bestWPM = practiceMode === 'ngrams'
    ? historicalStats.ngram_best_wpm
    : historicalStats.best_wpm;
  // ... similar for other stats
```

- [ ] **Step 2: Add slowest n-grams section**

```jsx
{practiceMode === 'ngrams' && sessionStats.bigram_seek_time && (
  <Box mt={4}>
    <Text fontSize="sm" color="gray.400">Slowest n-grams this session:</Text>
    {getSlowestNgrams(sessionStats).map(({ ngram, avg }) => (
      <Text key={ngram} fontFamily="mono">{ngram}: {avg.toFixed(0)}ms</Text>
    ))}
  </Box>
)}
```

- [ ] **Step 3: Test in browser**

Complete an n-gram round and verify results show mode-specific stats.

- [ ] **Step 4: Commit**

```bash
git add web/src/components/ResultsScreen.jsx
git commit -m "feat(web): adapt ResultsScreen for n-gram mode stats"
```

---

## Task 16: Final Integration Test

- [ ] **Step 1: Run all tests**

```bash
go test ./... -v
```

- [ ] **Step 2: Manual testing checklist**

- [ ] Options screen shows practice mode toggle (TUI)
- [ ] N-gram mode displays sequences correctly in block letters
- [ ] Space advances between n-grams
- [ ] Header shows "N-gram Training" in n-gram mode
- [ ] Progress shows "N-gram X/Y" in n-gram mode
- [ ] Results show n-gram-specific stats
- [ ] Results show slowest n-grams insight
- [ ] Switching modes preserves separate statistics
- [ ] Web frontend mirrors TUI behaviour
- [ ] Settings persist across sessions

- [ ] **Step 3: Commit final changes**

```bash
git add -A
git commit -m "feat: complete n-gram training mode implementation"
```

---

Plan complete and saved to `docs/superpowers/plans/2026-04-18-ngram-training.md`. Two execution options:

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?
