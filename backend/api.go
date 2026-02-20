// Package backend provides the game engine and API for the typing practice application.
// This package handles all game logic, statistics tracking, and state management.
// The frontend communicates exclusively through the GameAPI interface.
package backend

import (
	"time"

	"github.com/timlinux/baboon/stats"
)

// GameAPI defines the public interface for the typing game backend.
// All frontend interactions with the game engine go through this interface.
type GameAPI interface {
	// Game Lifecycle
	// --------------

	// StartRound initialises a new round with fresh words and resets session stats.
	StartRound()

	// Input Handling
	// ---------------

	// ProcessKeystroke handles a character input from the user (legacy, no timing).
	// Returns the result indicating whether the keystroke was correct and other state changes.
	ProcessKeystroke(char string) KeystrokeResult

	// ProcessKeystrokeWithTiming handles a character input with frontend-measured seek time.
	// The seekTimeMs parameter is the time since the last keystroke, measured on the frontend
	// to avoid network latency affecting timing accuracy.
	ProcessKeystrokeWithTiming(char string, seekTimeMs int64) KeystrokeResult

	// ProcessBackspace removes the last typed character.
	// Returns true if there was a character to remove.
	ProcessBackspace() bool

	// ProcessSpace handles the space key, potentially advancing to the next word (legacy).
	// Returns the result indicating whether the word advanced or round completed.
	ProcessSpace() SpaceResult

	// ProcessSpaceWithTiming handles the space key with frontend-measured seek time.
	ProcessSpaceWithTiming(seekTimeMs int64) SpaceResult

	// SubmitTiming submits the final timing data from the frontend when a round completes.
	// This ensures duration calculations use frontend timestamps, avoiding latency effects.
	SubmitTiming(startTime, endTime time.Time, durationMs int64)

	// State Queries
	// -------------

	// GetGameState returns a snapshot of the current game state for rendering.
	GetGameState() GameState

	// GetSessionStats returns the current session statistics.
	GetSessionStats() *stats.Stats

	// GetHistoricalStats returns the historical statistics.
	GetHistoricalStats() *stats.HistoricalStats

	// Persistence
	// -----------

	// SaveStats persists the current historical stats to disk.
	SaveStats() error
}

// KeystrokeResult contains the outcome of processing a keystroke.
type KeystrokeResult struct {
	// IsCorrect indicates whether the typed character matched the expected character.
	IsCorrect bool

	// TimerStarted indicates whether this keystroke started the game timer.
	TimerStarted bool

	// CharIndex is the position of the character that was typed.
	CharIndex int
}

// SpaceResult contains the outcome of processing the space key.
type SpaceResult struct {
	// Advanced indicates whether the game successfully moved to the next word.
	Advanced bool

	// RoundComplete indicates whether all words have been completed.
	RoundComplete bool

	// TreatedAsError indicates whether the space was treated as an incorrect character
	// (when pressed before the word was fully typed).
	TreatedAsError bool
}

// GameState contains a snapshot of the current game state for rendering.
// This is a read-only view of the game state that the frontend uses for display.
type GameState struct {
	// Words is the list of words in the current round.
	Words []string

	// CurrentWordIdx is the index of the current word being typed.
	CurrentWordIdx int

	// CurrentInput is the user's current typed input for the current word.
	CurrentInput string

	// TimerStarted indicates whether the game timer has started.
	TimerStarted bool

	// PunctuationMode indicates whether punctuation mode is enabled.
	PunctuationMode bool

	// Progress returns current word number and total words.
	WordNumber int
	TotalWords int

	// LiveWPM is the current words-per-minute calculation.
	LiveWPM float64

	// CurrentWord is the word currently being typed.
	CurrentWord string

	// PreviousWord is the word that was just completed (empty if first word).
	PreviousWord string

	// NextWord is the upcoming word (empty if on last word).
	// Deprecated: Use NextWords instead.
	NextWord string

	// NextWords contains the next 3 upcoming words (or fewer if near the end).
	NextWords []string
}

// Config holds configuration options for creating a new game engine.
type Config struct {
	// PunctuationMode enables punctuation between words.
	PunctuationMode bool

	// WordsPerRound is the number of words per round.
	WordsPerRound int

	// CharactersPerRound is the target total characters per round.
	CharactersPerRound int
}

// DefaultConfig returns the default game configuration.
func DefaultConfig() Config {
	return Config{
		PunctuationMode:    false,
		WordsPerRound:      30,
		CharactersPerRound: 150,
	}
}
