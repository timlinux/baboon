package frontend

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/timlinux/baboon/backend"
	"github.com/timlinux/baboon/settings"
)

// GameState represents the current state of the game UI
type GameState int

const (
	StateTyping GameState = iota
	StateCelebration
	StateResults
	StateOptions
	StateAbout
)

// tickMsg is sent periodically to update the WPM display
type tickMsg time.Time

// animTickMsg is sent to update animations
type animTickMsg time.Time

// celebrationTickMsg is sent to update celebration animations
type celebrationTickMsg time.Time

// Model is the Bubble Tea model for the typing game
type Model struct {
	// Backend API
	api backend.GameAPI

	// UI state
	state    GameState
	width    int
	height   int
	renderer *Renderer
	animator *Animator

	// Carousel animation for typing screen
	carouselAnimator *CarouselAnimator
	lastWordIdx      int // Track word changes to trigger animations

	// Local timing state (to avoid network latency affecting measurements)
	timerStarted bool
	startTime    time.Time
	lastKeyTime  time.Time
	correctChars int // For live WPM calculation

	// Settings
	settings          *settings.Settings
	optionsCursor     int  // Current selection in options menu
	optionsFromTyping bool // Whether options was opened from typing screen

	// Celebration state
	celebrationState *CelebrationState
}

// NewModel creates a new Model with the given backend API
func NewModel(api backend.GameAPI) Model {
	// Load settings (use defaults if error)
	s, _ := settings.Load()
	return Model{
		api:              api,
		state:            StateTyping,
		renderer:         NewRenderer(80, 24), // Default size, will be updated
		carouselAnimator: NewCarouselAnimator(),
		lastWordIdx:      0,
		settings:         s,
	}
}

// Init initializes the model and returns the initial command
func (m Model) Init() tea.Cmd {
	// Don't start tick loop immediately - only start when timer starts
	// This reduces unnecessary screen updates and flicker
	return nil
}

// Update handles messages and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Only continue ticking if timer is running (for live WPM updates)
		if m.timerStarted && m.state == StateTyping {
			return m, tickCmd()
		}
		// Stop ticking when timer not running to reduce flicker
		return m, nil

	case animTickMsg:
		// Handle results screen animations
		if m.state == StateResults && m.animator != nil {
			m.animator.Update()
			// Stop animation loop once all animations are complete
			if !m.animator.IsComplete() {
				return m, animTickCmd()
			}
		}
		// Handle typing screen carousel animations
		if m.state == StateTyping && m.carouselAnimator != nil && m.carouselAnimator.IsAnimating {
			m.carouselAnimator.Update()
			if m.carouselAnimator.IsAnimating {
				return m, animTickCmd()
			}
		}
		return m, nil

	case celebrationTickMsg:
		return m.updateCelebration()

	case tea.KeyMsg:
		switch m.state {
		case StateTyping:
			return m.handleTypingInput(msg)
		case StateCelebration:
			return m.handleCelebrationInput(msg)
		case StateResults:
			return m.handleResultsInput(msg)
		case StateOptions:
			return m.handleOptionsInput(msg)
		case StateAbout:
			return m.handleAboutInput(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.renderer.SetSize(msg.Width, msg.Height)
	}

	return m, nil
}

// View renders the current state
// Flicker is reduced by:
// 1. Rounding WPM to integers (no constant small changes)
// 2. Using 500ms tick interval instead of 100ms
// 3. Bubble Tea's built-in view diffing
func (m Model) View() string {
	switch m.state {
	case StateTyping:
		gameState := m.api.GetGameState()

		// Calculate live WPM locally, rounded to integer for stable display
		if m.timerStarted && m.correctChars > 0 {
			elapsed := time.Since(m.startTime).Minutes()
			if elapsed > 0 {
				// Round to integer to prevent constant tiny updates
				gameState.LiveWPM = float64(int((float64(m.correctChars) / 5.0) / elapsed))
			}
		}
		gameState.TimerStarted = m.timerStarted

		return m.renderer.RenderTypingScreenAnimated(gameState, m.carouselAnimator, m.settings)

	case StateCelebration:
		return m.renderer.RenderCelebrationScreen(m.celebrationState)

	case StateResults:
		return m.renderer.RenderResultsScreen(
			m.api.GetSessionStats(),
			m.api.GetHistoricalStats(),
			m.animator,
		)
	case StateOptions:
		return m.renderer.RenderOptionsScreen(m.settings, m.optionsCursor)
	case StateAbout:
		return m.renderer.RenderAboutScreen()
	}
	return ""
}

// handleTypingInput processes keyboard input during typing
func (m Model) handleTypingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	now := time.Now()

	// Check if this key should advance to the next word
	isAdvanceKey := func() bool {
		switch m.settings.AdvanceKey {
		case settings.AdvanceKeySpace:
			return msg.Type == tea.KeySpace
		case settings.AdvanceKeyEnter:
			return msg.Type == tea.KeyEnter
		case settings.AdvanceKeyEither:
			return msg.Type == tea.KeySpace || msg.Type == tea.KeyEnter
		}
		return msg.Type == tea.KeySpace // default
	}

	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyTab:
		// Restart the current round
		m.api.StartRound()
		m.carouselAnimator = NewCarouselAnimator()
		m.timerStarted = false
		m.startTime = time.Time{}
		m.lastKeyTime = time.Time{}
		m.correctChars = 0
		return m, nil

	case tea.KeySpace, tea.KeyEnter:
		// Only process if this is the configured advance key
		if !isAdvanceKey() {
			return m, nil
		}

		// Calculate seek time locally
		var seekTimeMs int64
		if m.timerStarted && !m.lastKeyTime.IsZero() {
			seekTimeMs = now.Sub(m.lastKeyTime).Milliseconds()
		}
		m.lastKeyTime = now

		result := m.api.ProcessSpaceWithTiming(seekTimeMs)
		if result.RoundComplete {
			// Capture historical bests BEFORE updating stats
			historical := m.api.GetHistoricalStats()
			oldBestWPM := historical.BestWPM
			oldBestAccuracy := historical.BestAccuracy
			oldBestTime := historical.BestTime

			// Send final timing to backend (this updates historical stats)
			var durationMs int64
			if m.timerStarted {
				durationMs = now.Sub(m.startTime).Milliseconds()
			}
			m.api.SubmitTiming(m.startTime, now, durationMs)
			m.api.SaveStats()
			return m.handleRoundComplete(oldBestWPM, oldBestAccuracy, oldBestTime)
		} else if result.Advanced {
			// Trigger carousel animation when moving to next word
			m.carouselAnimator.TriggerTransition()
			return m, animTickCmd()
		}

	case tea.KeyBackspace:
		m.api.ProcessBackspace()

	case tea.KeyCtrlW:
		// Ctrl+W or Ctrl+Backspace clears all typed characters
		m.api.ClearInput()

	case tea.KeyCtrlO:
		// Open options with Ctrl+O (only before timer starts)
		if !m.timerStarted {
			m.state = StateOptions
			m.optionsCursor = 0
			m.optionsFromTyping = true
			return m, nil
		}

	case tea.KeyCtrlA:
		// Open about screen with Ctrl+A (only before timer starts)
		if !m.timerStarted {
			m.state = StateAbout
			return m, nil
		}

	case tea.KeyRunes:
		char := string(msg.Runes)

		// Calculate seek time locally before sending to backend
		var seekTimeMs int64
		if m.timerStarted && !m.lastKeyTime.IsZero() {
			seekTimeMs = now.Sub(m.lastKeyTime).Milliseconds()
		}

		result := m.api.ProcessKeystrokeWithTiming(char, seekTimeMs)

		// Start timer on first correct character (tracked locally)
		if result.TimerStarted && !m.timerStarted {
			m.timerStarted = true
			m.startTime = now
			m.lastKeyTime = now
			// Start tick loop for live WPM updates now that timer is running
			return m, tickCmd()
		}

		// Track correct chars for local live WPM
		if result.IsCorrect {
			m.correctChars++
		} else if m.settings.PerfectMode && m.timerStarted {
			// Perfect mode: restart round on first mistake
			m.api.StartRound()
			m.carouselAnimator = NewCarouselAnimator()
			m.timerStarted = false
			m.startTime = time.Time{}
			m.lastKeyTime = time.Time{}
			m.correctChars = 0
			return m, nil
		}

		m.lastKeyTime = now

		// Check if round completed (last correct char of last word)
		if result.RoundComplete {
			// Capture historical bests BEFORE updating stats
			historical := m.api.GetHistoricalStats()
			oldBestWPM := historical.BestWPM
			oldBestAccuracy := historical.BestAccuracy
			oldBestTime := historical.BestTime

			// Send final timing to backend (this updates historical stats)
			var durationMs int64
			if m.timerStarted {
				durationMs = now.Sub(m.startTime).Milliseconds()
			}
			m.api.SubmitTiming(m.startTime, now, durationMs)
			m.api.SaveStats()
			return m.handleRoundComplete(oldBestWPM, oldBestAccuracy, oldBestTime)
		}
	}

	return m, nil
}

// handleRoundComplete processes the end of a typing round, checking for personal bests
// oldBest* parameters contain the historical values BEFORE the session was recorded
func (m Model) handleRoundComplete(oldBestWPM, oldBestAccuracy, oldBestTime float64) (tea.Model, tea.Cmd) {
	session := m.api.GetSessionStats()

	// Check for new bests against the OLD historical values (before this session updated them)
	var bestTypes BestType

	// Check WPM: session WPM > old best WPM
	if session.WPM > oldBestWPM {
		bestTypes |= BestWPM
	}

	// Check Accuracy: session accuracy > old best accuracy
	if session.Accuracy > oldBestAccuracy {
		bestTypes |= BestAccuracy
	}

	// Check Time: session time < old best time (lower is better)
	// Also triggers if no best time recorded yet (oldBestTime == 0)
	if oldBestTime == 0 || session.Duration.Seconds() < oldBestTime {
		bestTypes |= BestTime
	}

	// Reset timing state
	m.timerStarted = false
	m.correctChars = 0

	if bestTypes != 0 {
		// New best! Start celebration
		m.celebrationState = NewCelebrationState(
			m.width, m.height, bestTypes,
			session.WPM, session.Accuracy, session.Duration.Seconds(),
		)
		m.state = StateCelebration
		return m, celebrationTickCmd()
	}

	// No new best, go directly to results
	m.state = StateResults
	m.animator = NewAnimator()
	return m, animTickCmd()
}

// handleResultsInput processes keyboard input on results screen
func (m Model) handleResultsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyEnter, tea.KeyTab:
		m.api.StartRound()
		m.state = StateTyping
		m.animator = nil
		// Reset carousel animator for new round
		m.carouselAnimator = NewCarouselAnimator()
		// Reset timing state for new round
		m.timerStarted = false
		m.startTime = time.Time{}
		m.lastKeyTime = time.Time{}
		m.correctChars = 0

	case tea.KeyCtrlO:
		// Open options with Ctrl+O
		m.state = StateOptions
		m.optionsCursor = 0
		m.optionsFromTyping = false
		return m, nil
	}

	return m, nil
}

// handleOptionsInput processes keyboard input on options screen
func (m Model) handleOptionsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	const totalOptions = 5 // 3 advance key options + 2 perfect mode options

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEsc:
		// Return to previous screen
		if m.optionsFromTyping {
			m.state = StateTyping
		} else {
			m.state = StateResults
		}
		return m, nil

	case tea.KeyUp, tea.KeyShiftTab:
		// Move cursor up (wrap around)
		if m.optionsCursor > 0 {
			m.optionsCursor--
		} else {
			m.optionsCursor = totalOptions - 1 // Wrap to last option
		}

	case tea.KeyDown, tea.KeyTab:
		// Move cursor down (wrap around)
		if m.optionsCursor < totalOptions-1 {
			m.optionsCursor++
		} else {
			m.optionsCursor = 0 // Wrap to first option
		}

	case tea.KeyEnter, tea.KeySpace:
		// Select current option based on cursor position
		if m.optionsCursor < 3 {
			// Advance key options (0, 1, 2)
			m.settings.AdvanceKey = settings.AdvanceKey(m.optionsCursor)
		} else {
			// Perfect mode options (3 = Off, 4 = On)
			m.settings.PerfectMode = m.optionsCursor == 4
		}
		// Save settings
		_ = m.settings.Save()
		// Return to previous screen
		if m.optionsFromTyping {
			m.state = StateTyping
		} else {
			m.state = StateResults
		}
		return m, nil

	case tea.KeyRunes:
		char := string(msg.Runes)
		// Quick select with number keys
		switch char {
		case "1":
			m.settings.AdvanceKey = settings.AdvanceKeySpace
			_ = m.settings.Save()
			if m.optionsFromTyping {
				m.state = StateTyping
			} else {
				m.state = StateResults
			}
			return m, nil
		case "2":
			m.settings.AdvanceKey = settings.AdvanceKeyEnter
			_ = m.settings.Save()
			if m.optionsFromTyping {
				m.state = StateTyping
			} else {
				m.state = StateResults
			}
			return m, nil
		case "3":
			m.settings.AdvanceKey = settings.AdvanceKeyEither
			_ = m.settings.Save()
			if m.optionsFromTyping {
				m.state = StateTyping
			} else {
				m.state = StateResults
			}
			return m, nil
		case "4":
			m.settings.PerfectMode = false
			_ = m.settings.Save()
			if m.optionsFromTyping {
				m.state = StateTyping
			} else {
				m.state = StateResults
			}
			return m, nil
		case "5":
			m.settings.PerfectMode = true
			_ = m.settings.Save()
			if m.optionsFromTyping {
				m.state = StateTyping
			} else {
				m.state = StateResults
			}
			return m, nil
		}
	}

	return m, nil
}

// updateCelebration handles celebration animation updates
func (m Model) updateCelebration() (tea.Model, tea.Cmd) {
	if m.celebrationState == nil {
		m.state = StateResults
		m.animator = NewAnimator()
		return m, animTickCmd()
	}

	dt := time.Since(m.celebrationState.LastUpdate).Seconds()
	m.celebrationState.Update(dt)

	if m.celebrationState.IsComplete() {
		m.state = StateResults
		m.celebrationState = nil
		m.animator = NewAnimator()
		return m, animTickCmd()
	}

	return m, celebrationTickCmd()
}

// handleCelebrationInput processes keyboard input during celebration
func (m Model) handleCelebrationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	default:
		// Any other key skips to results
		m.state = StateResults
		m.celebrationState = nil
		m.animator = NewAnimator()
		return m, animTickCmd()
	}
}

// handleAboutInput processes keyboard input on about screen
func (m Model) handleAboutInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc, tea.KeyEnter, tea.KeySpace:
		// Return to typing screen
		m.state = StateTyping
		return m, nil
	}
	return m, nil
}

// tickCmd returns a command that sends tick messages
// Using 500ms interval since WPM is displayed as integer (no need for faster updates)
func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// animTickCmd returns a command that sends animation tick messages
func animTickCmd() tea.Cmd {
	return tea.Tick(GetAnimationInterval(), func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}

// celebrationTickCmd returns a command that sends celebration tick messages
func celebrationTickCmd() tea.Cmd {
	return tea.Tick(CelebrationTickInterval, func(t time.Time) tea.Msg {
		return celebrationTickMsg(t)
	})
}
