package frontend

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/timlinux/baboon/backend"
)

// GameState represents the current state of the game UI
type GameState int

const (
	StateTyping GameState = iota
	StateResults
)

// tickMsg is sent periodically to update the WPM display
type tickMsg time.Time

// animTickMsg is sent to update animations
type animTickMsg time.Time

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
}

// NewModel creates a new Model with the given backend API
func NewModel(api backend.GameAPI) Model {
	return Model{
		api:              api,
		state:            StateTyping,
		renderer:         NewRenderer(80, 24), // Default size, will be updated
		carouselAnimator: NewCarouselAnimator(),
		lastWordIdx:      0,
	}
}

// Init initializes the model and returns the initial command
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// Update handles messages and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tickCmd()

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

	case tea.KeyMsg:
		switch m.state {
		case StateTyping:
			return m.handleTypingInput(msg)
		case StateResults:
			return m.handleResultsInput(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.renderer.SetSize(msg.Width, msg.Height)
	}

	return m, nil
}

// View renders the current state
func (m Model) View() string {
	switch m.state {
	case StateTyping:
		gameState := m.api.GetGameState()
		// Override live WPM with locally calculated value (avoids network latency)
		if m.timerStarted && m.correctChars > 0 {
			elapsed := time.Since(m.startTime).Minutes()
			if elapsed > 0 {
				gameState.LiveWPM = (float64(m.correctChars) / 5.0) / elapsed
			}
		}
		gameState.TimerStarted = m.timerStarted
		return m.renderer.RenderTypingScreenAnimated(gameState, m.carouselAnimator)
	case StateResults:
		return m.renderer.RenderResultsScreen(
			m.api.GetSessionStats(),
			m.api.GetHistoricalStats(),
			m.animator,
		)
	}
	return ""
}

// handleTypingInput processes keyboard input during typing
func (m Model) handleTypingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	now := time.Now()

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

	case tea.KeySpace:
		// Calculate seek time locally
		var seekTimeMs int64
		if m.timerStarted && !m.lastKeyTime.IsZero() {
			seekTimeMs = now.Sub(m.lastKeyTime).Milliseconds()
		}
		m.lastKeyTime = now

		result := m.api.ProcessSpaceWithTiming(seekTimeMs)
		if result.RoundComplete {
			// Send final timing to backend
			var durationMs int64
			if m.timerStarted {
				durationMs = now.Sub(m.startTime).Milliseconds()
			}
			m.api.SubmitTiming(m.startTime, now, durationMs)
			m.api.SaveStats()
			m.state = StateResults
			m.animator = NewAnimator()
			// Reset timing state
			m.timerStarted = false
			m.correctChars = 0
			return m, animTickCmd()
		} else if result.Advanced {
			// Trigger carousel animation when moving to next word
			m.carouselAnimator.TriggerTransition()
			return m, animTickCmd()
		}

	case tea.KeyBackspace:
		m.api.ProcessBackspace()

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
		}

		// Track correct chars for local live WPM
		if result.IsCorrect {
			m.correctChars++
		}

		m.lastKeyTime = now
	}

	return m, nil
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
	}

	return m, nil
}

// tickCmd returns a command that sends tick messages
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// animTickCmd returns a command that sends animation tick messages
func animTickCmd() tea.Cmd {
	return tea.Tick(GetAnimationInterval(), func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}
