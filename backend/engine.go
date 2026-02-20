package backend

import (
	"math/rand"
	"time"

	"github.com/timlinux/baboon/stats"
	"github.com/timlinux/baboon/words"
)

// Punctuation characters used in punctuation mode
var punctuationChars = []string{",", ".", ";", ":", "!", "?"}

// Engine implements the GameAPI interface and manages all game logic.
type Engine struct {
	config     Config
	rng        *rand.Rand
	historical *stats.HistoricalStats
	session    *stats.Stats
	words      []string
	wordIdx    int
	input      string
	started    bool

	// lastLetter is tracked here for bigram/SFB detection (not timing related)
	lastLetter string
}

// NewEngine creates a new game engine with the given configuration.
func NewEngine(config Config) (*Engine, error) {
	historical, err := stats.LoadHistoricalStats()
	if err != nil {
		return nil, err
	}

	e := &Engine{
		config:     config,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
		historical: historical,
	}

	e.StartRound()
	return e, nil
}

// StartRound initialises a new round with fresh words and resets session stats.
func (e *Engine) StartRound() {
	// Get letter data for weighted word selection
	letterData := e.getLetterData()
	e.words = words.GetRandomWordsFixedCount(
		e.config.WordsPerRound,
		e.config.CharactersPerRound,
		e.rng.Intn,
		letterData,
	)

	// Create new session stats
	e.session = &stats.Stats{
		LetterAccuracy:    make(map[string]stats.LetterStats),
		LetterSeekTime:    make(map[string]stats.LetterSeekStats),
		BigramSeekTime:    make(map[string]stats.BigramSeekStats),
		FingerStats:       make(map[int]stats.FingerStat),
		HandStats:         make(map[int]stats.HandStat),
		RowStats:          make(map[int]stats.RowStat),
		ErrorSubstitution: make(map[string]map[string]int),
		SeekTimes:         make([]int64, 0),
	}

	// Record all letters as presented (before adding punctuation)
	for _, word := range e.words {
		for _, char := range word {
			if char >= 'a' && char <= 'z' {
				e.session.RecordLetterPresented(string(char))
				if finger := stats.GetFinger(char); finger >= 0 {
					e.session.RecordFingerPresented(finger)
				}
				if hand := stats.GetHand(char); hand >= 0 {
					e.session.RecordHandPresented(hand)
				}
				if row := stats.GetRow(char); row >= 0 {
					e.session.RecordRowPresented(row)
				}
			}
		}
	}

	// Add punctuation if enabled
	if e.config.PunctuationMode {
		for i := 0; i < len(e.words)-1; i++ {
			punct := punctuationChars[e.rng.Intn(len(punctuationChars))]
			e.words[i] = e.words[i] + punct
		}
	}

	e.wordIdx = 0
	e.input = ""
	e.started = false
	e.lastLetter = ""
}

// ProcessKeystroke handles a character input from the user (legacy, no timing).
// This calls ProcessKeystrokeWithTiming with 0 seek time.
func (e *Engine) ProcessKeystroke(char string) KeystrokeResult {
	return e.ProcessKeystrokeWithTiming(char, 0)
}

// ProcessKeystrokeWithTiming handles a character input with frontend-measured seek time.
// All timing is done on the frontend to avoid network latency affecting measurements.
func (e *Engine) ProcessKeystrokeWithTiming(char string, seekTimeMs int64) KeystrokeResult {
	if e.wordIdx >= len(e.words) {
		return KeystrokeResult{}
	}

	currentWord := e.words[e.wordIdx]
	inputIdx := len(e.input)
	result := KeystrokeResult{CharIndex: inputIdx}

	// Check if this should start the timer (first correct character of first word)
	if !e.started && e.wordIdx == 0 && inputIdx == 0 {
		if len(currentWord) > 0 && char == string(currentWord[0]) {
			e.started = true
			result.TimerStarted = true
		}
	}

	e.input += char
	e.session.TotalCharacters++

	// Check if character matches
	isCorrect := inputIdx < len(currentWord) && e.input[inputIdx] == currentWord[inputIdx]
	result.IsCorrect = isCorrect

	if isCorrect {
		e.session.CorrectChars++
		expectedChar := currentWord[inputIdx]
		expectedLetter := string(expectedChar)

		// Only record letter stats for actual letters (a-z), not punctuation
		isLetter := expectedChar >= 'a' && expectedChar <= 'z'
		if isLetter {
			e.session.RecordLetterCorrect(expectedLetter)

			finger := stats.GetFinger(rune(expectedChar))
			hand := stats.GetHand(rune(expectedChar))
			row := stats.GetRow(rune(expectedChar))

			// Record seek time only for correct keystrokes
			// Exclude first letter of each word (inputIdx > 0)
			// Use frontend-provided seek time
			if e.started && inputIdx > 0 && seekTimeMs > 0 && seekTimeMs < 5000 {
				e.session.RecordLetterSeekTime(expectedLetter, seekTimeMs)
				e.session.RecordSeekTime(seekTimeMs)

				// Record bigram timing
				if e.lastLetter != "" {
					bigram := e.lastLetter + expectedLetter
					e.session.RecordBigramSeekTime(bigram, seekTimeMs)

					// Check for same-finger bigram
					lastChar := rune(e.lastLetter[0])
					if stats.IsSameFingerBigram(lastChar, rune(expectedChar)) {
						e.session.RecordSFB(seekTimeMs)
					}

					// Track hand alternation
					lastHand := stats.GetHand(lastChar)
					if lastHand >= 0 && hand >= 0 {
						e.session.RecordHandTransition(lastHand != hand)
					}
				}

				// Record finger, hand, row stats with timing
				if finger >= 0 {
					e.session.RecordFingerCorrect(finger, seekTimeMs)
				}
				if hand >= 0 {
					e.session.RecordHandCorrect(hand, seekTimeMs)
				}
				if row >= 0 {
					e.session.RecordRowCorrect(row, seekTimeMs)
				}
			} else {
				// Record finger, hand, row stats without timing (first char of word)
				if finger >= 0 {
					e.session.RecordFingerCorrect(finger, 0)
				}
				if hand >= 0 {
					e.session.RecordHandCorrect(hand, 0)
				}
				if row >= 0 {
					e.session.RecordRowCorrect(row, 0)
				}
			}

			e.lastLetter = expectedLetter
		}
	} else {
		e.session.IncorrectChars++
		// Track error substitution pattern
		if inputIdx < len(currentWord) {
			expectedChar := currentWord[inputIdx]
			typedChar := e.input[inputIdx]
			if expectedChar >= 'a' && expectedChar <= 'z' && typedChar >= 'a' && typedChar <= 'z' {
				e.session.RecordErrorSubstitution(string(expectedChar), string(typedChar))
			}
		}
	}

	return result
}

// ProcessBackspace removes the last typed character.
func (e *Engine) ProcessBackspace() bool {
	if len(e.input) > 0 {
		e.input = e.input[:len(e.input)-1]
		return true
	}
	return false
}

// ProcessSpace handles the space key (legacy, no timing).
func (e *Engine) ProcessSpace() SpaceResult {
	return e.ProcessSpaceWithTiming(0)
}

// ProcessSpaceWithTiming handles the space key with frontend-measured seek time.
func (e *Engine) ProcessSpaceWithTiming(seekTimeMs int64) SpaceResult {
	if e.wordIdx >= len(e.words) {
		return SpaceResult{}
	}

	currentWord := e.words[e.wordIdx]

	// Only advance if all letters have been typed
	if len(e.input) >= len(currentWord) {
		e.session.WordsCompleted++
		e.input = ""
		e.wordIdx++
		e.lastLetter = "" // Reset for new word

		if e.wordIdx >= len(e.words) {
			// Round complete - don't calculate yet, wait for SubmitTiming
			return SpaceResult{Advanced: true, RoundComplete: true}
		}
		return SpaceResult{Advanced: true}
	}

	// Treat space as incorrect if word not complete
	if len(e.input) > 0 || e.started {
		e.input += " "
		e.session.TotalCharacters++
		e.session.IncorrectChars++
		return SpaceResult{TreatedAsError: true}
	}

	return SpaceResult{}
}

// SubmitTiming receives final timing data from the frontend and calculates stats.
// This ensures duration calculations use frontend timestamps, avoiding latency effects.
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

	// Update historical stats
	e.historical.UpdateHistorical(e.session)
}

// GetGameState returns a snapshot of the current game state.
func (e *Engine) GetGameState() GameState {
	state := GameState{
		Words:           e.words,
		CurrentWordIdx:  e.wordIdx,
		CurrentInput:    e.input,
		TimerStarted:    e.started,
		PunctuationMode: e.config.PunctuationMode,
		WordNumber:      e.wordIdx + 1,
		TotalWords:      len(e.words),
	}

	if e.wordIdx < len(e.words) {
		state.CurrentWord = e.words[e.wordIdx]
	}
	if e.wordIdx > 0 {
		state.PreviousWord = e.words[e.wordIdx-1]
	}
	if e.wordIdx < len(e.words)-1 {
		state.NextWord = e.words[e.wordIdx+1]
	}

	// Note: LiveWPM is now calculated on the frontend to avoid network latency
	// The frontend tracks timing locally and overrides this value
	state.LiveWPM = 0

	return state
}

// GetSessionStats returns the current session statistics.
func (e *Engine) GetSessionStats() *stats.Stats {
	return e.session
}

// GetHistoricalStats returns the historical statistics.
func (e *Engine) GetHistoricalStats() *stats.HistoricalStats {
	return e.historical
}

// SaveStats persists the current historical stats to disk.
func (e *Engine) SaveStats() error {
	return stats.SaveHistoricalStats(e.historical)
}

// getLetterData extracts letter frequency and accuracy data for word selection.
func (e *Engine) getLetterData() words.LetterData {
	data := make(words.LetterData)
	if e.historical == nil || e.historical.LetterAccuracy == nil {
		return data
	}
	for letter, letterStats := range e.historical.LetterAccuracy {
		data[letter] = words.LetterStats{
			Presented: letterStats.Presented,
			Correct:   letterStats.Correct,
		}
	}
	return data
}

// Ensure Engine implements GameAPI
var _ GameAPI = (*Engine)(nil)
