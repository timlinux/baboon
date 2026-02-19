package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/timlinux/baboon/font"
	"github.com/timlinux/baboon/stats"
	"github.com/timlinux/baboon/words"
)

const wordsPerRound = 30

// WPM bar configuration
const (
	maxWPM   = 120.0 // Maximum WPM for the bar scale
	barWidth = 50    // Width of the WPM bar in characters
)

// Game states
type gameState int

const (
	stateTyping gameState = iota
	stateResults
)

// tickMsg is sent periodically to update the WPM display
type tickMsg time.Time

type model struct {
	state          gameState
	words          []string
	currentWordIdx int
	currentInput   string
	stats          *stats.Stats
	historical     *stats.HistoricalStats
	width          int
	height         int
	rng            *rand.Rand
	timerStarted   bool
}

func initialModel() model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	historical, _ := stats.LoadHistoricalStats()

	m := model{
		state:        stateTyping,
		words:        words.GetRandomWords(wordsPerRound, rng.Intn),
		historical:   historical,
		rng:          rng,
		timerStarted: false,
		stats:        &stats.Stats{},
	}
	return m
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Continue ticking to update WPM display
		return m, tickCmd()

	case tea.KeyMsg:
		switch m.state {
		case stateTyping:
			return m.handleTypingInput(msg)
		case stateResults:
			return m.handleResultsInput(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) handleTypingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	currentWord := m.words[m.currentWordIdx]

	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeySpace:
		// Check if word is complete (correct or not, move to next)
		if len(m.currentInput) > 0 {
			m.stats.WordsCompleted++
			m.currentInput = ""
			m.currentWordIdx++

			if m.currentWordIdx >= len(m.words) {
				// Round complete - show results
				m.stats.Calculate()
				m.historical.UpdateHistorical(m.stats)
				stats.SaveHistoricalStats(m.historical)
				m.state = stateResults
			}
		}

	case tea.KeyBackspace:
		if len(m.currentInput) > 0 {
			m.currentInput = m.currentInput[:len(m.currentInput)-1]
		}

	case tea.KeyRunes:
		char := string(msg.Runes)
		inputIdx := len(m.currentInput)

		// Start timer on first correct character of first word
		if !m.timerStarted && m.currentWordIdx == 0 && inputIdx == 0 {
			if len(currentWord) > 0 && char == string(currentWord[0]) {
				m.timerStarted = true
				m.stats.StartTime = time.Now()
			}
		}

		m.currentInput += char
		m.stats.TotalCharacters++

		// Check if character matches
		if inputIdx < len(currentWord) && m.currentInput[inputIdx] == currentWord[inputIdx] {
			m.stats.CorrectChars++
		} else {
			m.stats.IncorrectChars++
		}
	}

	return m, nil
}

func (m model) handleResultsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyEnter:
		// Start new round
		m.state = stateTyping
		m.words = words.GetRandomWords(wordsPerRound, m.rng.Intn)
		m.currentWordIdx = 0
		m.currentInput = ""
		m.timerStarted = false
		m.stats = &stats.Stats{}
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateTyping:
		return m.renderTyping()
	case stateResults:
		return m.renderResults()
	}
	return ""
}

// calculateCurrentWPM calculates WPM based on current progress
func (m model) calculateCurrentWPM() float64 {
	if !m.timerStarted || m.stats.CorrectChars == 0 {
		return 0
	}

	elapsed := time.Since(m.stats.StartTime).Minutes()
	if elapsed <= 0 {
		return 0
	}

	// WPM = (correct characters / 5) / minutes
	return (float64(m.stats.CorrectChars) / 5.0) / elapsed
}

// renderWPMBar creates a beautiful gradient progress bar for WPM
func (m model) renderWPMBar() string {
	wpm := m.calculateCurrentWPM()

	// Calculate fill percentage
	fillPercent := wpm / maxWPM
	if fillPercent > 1.0 {
		fillPercent = 1.0
	}
	if fillPercent < 0 {
		fillPercent = 0
	}

	filledWidth := int(float64(barWidth) * fillPercent)
	emptyWidth := barWidth - filledWidth

	// Gradient colours from red (slow) through yellow to green (fast)
	// Using 256-colour palette for smooth gradient
	gradientColours := []string{
		"196", "202", "208", "214", "220", "226", // Red to yellow
		"190", "154", "118", "82", "46", "47", // Yellow to green
	}

	var bar strings.Builder

	// Build the filled portion with gradient
	for i := 0; i < filledWidth; i++ {
		// Calculate which colour to use based on position
		colourIdx := int(float64(i) / float64(barWidth) * float64(len(gradientColours)-1))
		if colourIdx >= len(gradientColours) {
			colourIdx = len(gradientColours) - 1
		}
		colour := gradientColours[colourIdx]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		bar.WriteString(style.Render("█"))
	}

	// Build the empty portion
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
	for i := 0; i < emptyWidth; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}

	// WPM label
	wpmLabel := fmt.Sprintf(" %.0f WPM", wpm)
	var labelStyle lipgloss.Style
	if wpm >= 60 {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true) // Green
	} else if wpm >= 40 {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true) // Yellow
	} else {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red
	}

	// Scale markers
	scaleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	scale := scaleStyle.Render("0                        60                       120")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		bar.String()+labelStyle.Render(wpmLabel),
		scale,
	)
}

func (m model) renderTyping() string {
	if m.currentWordIdx >= len(m.words) {
		return ""
	}

	currentWord := m.words[m.currentWordIdx]

	// Safety check: skip empty words (shouldn't happen but prevents lockup)
	if len(currentWord) == 0 {
		return "Loading..."
	}

	// Render word using custom block font
	letterLines := font.RenderWord(currentWord)

	// Define styles
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green - correct
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))    // Red - incorrect
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))   // Gray - not yet typed

	// Build colored output for each line
	coloredLines := make([]string, font.LetterHeight)

	for lineIdx := 0; lineIdx < font.LetterHeight; lineIdx++ {
		var lineBuilder strings.Builder

		for charIdx, letterLine := range letterLines[lineIdx] {
			var style lipgloss.Style

			if charIdx < len(m.currentInput) {
				// Character has been typed
				if charIdx < len(currentWord) && m.currentInput[charIdx] == currentWord[charIdx] {
					style = greenStyle
				} else {
					style = redStyle
				}
			} else {
				// Character not yet typed
				style = grayStyle
			}

			lineBuilder.WriteString(style.Render(letterLine))
			// Add spacing between letters
			if charIdx < len(letterLines[lineIdx])-1 {
				lineBuilder.WriteString(style.Render(" "))
			}
		}

		coloredLines[lineIdx] = lineBuilder.String()
	}

	coloredWord := strings.Join(coloredLines, "\n")

	// Progress indicator
	progress := fmt.Sprintf("Word %d/%d", m.currentWordIdx+1, len(m.words))
	progressStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)

	// Instructions
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	var helpText string
	if !m.timerStarted {
		helpText = "Type the first letter to start the timer | ESC to quit"
	} else {
		helpText = "Type the word, then press SPACE to continue | ESC to quit"
	}
	help := helpStyle.Render(helpText)

	// WPM Bar
	wpmBar := m.renderWPMBar()

	// Center the main content
	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		progressStyle.Render(progress),
		"",
		"",
		coloredWord,
		"",
		"",
		help,
	)

	// Calculate vertical positioning
	// Main content goes in center, WPM bar at bottom
	contentHeight := strings.Count(mainContent, "\n") + 1
	wpmBarHeight := strings.Count(wpmBar, "\n") + 1
	topPadding := (m.height - contentHeight - wpmBarHeight - 4) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Build full screen layout
	var fullContent strings.Builder

	// Top padding
	for i := 0; i < topPadding; i++ {
		fullContent.WriteString("\n")
	}

	// Center main content horizontally
	centeredMain := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, mainContent)
	fullContent.WriteString(centeredMain)

	// Spacer before WPM bar
	fullContent.WriteString("\n\n")

	// Bottom WPM bar (centered)
	centeredBar := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, wpmBar)
	fullContent.WriteString(centeredBar)

	// Fill remaining space
	currentHeight := topPadding + contentHeight + 2 + wpmBarHeight
	for i := currentHeight; i < m.height; i++ {
		fullContent.WriteString("\n")
	}

	return fullContent.String()
}

// renderStatBar creates a gradient bar for a stat value
// Returns a fixed-width string (barWidth + 2 for the star column)
func renderStatBar(value, maxValue float64, width int, isNewBest bool) string {
	fillPercent := value / maxValue
	if fillPercent > 1.0 {
		fillPercent = 1.0
	}
	if fillPercent < 0 {
		fillPercent = 0
	}

	filledWidth := int(float64(width) * fillPercent)
	emptyWidth := width - filledWidth

	// Gradient colours
	gradientColours := []string{
		"196", "202", "208", "214", "220", "226",
		"190", "154", "118", "82", "46", "47",
	}

	var bar strings.Builder

	for i := 0; i < filledWidth; i++ {
		colourIdx := int(float64(i) / float64(width) * float64(len(gradientColours)-1))
		if colourIdx >= len(gradientColours) {
			colourIdx = len(gradientColours) - 1
		}
		colour := gradientColours[colourIdx]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		bar.WriteString(style.Render("█"))
	}

	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
	for i := 0; i < emptyWidth; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}

	// Fixed-width star column (2 chars: space + star or two spaces)
	if isNewBest {
		starStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
		bar.WriteString(starStyle.Render(" *"))
	} else {
		bar.WriteString("  ")
	}

	return bar.String()
}

// renderTimeBar creates a bar for time (lower is better, so inverted)
// Returns a fixed-width string (barWidth + 2 for the star column)
func renderTimeBar(value, maxValue float64, width int, isNewBest bool) string {
	// Invert the percentage since lower time is better
	fillPercent := 1.0 - (value / maxValue)
	if fillPercent > 1.0 {
		fillPercent = 1.0
	}
	if fillPercent < 0 {
		fillPercent = 0
	}

	filledWidth := int(float64(width) * fillPercent)
	emptyWidth := width - filledWidth

	// Gradient colours (inverted - green for fast, red for slow)
	gradientColours := []string{
		"196", "202", "208", "214", "220", "226",
		"190", "154", "118", "82", "46", "47",
	}

	var bar strings.Builder

	for i := 0; i < filledWidth; i++ {
		colourIdx := int(float64(i) / float64(width) * float64(len(gradientColours)-1))
		if colourIdx >= len(gradientColours) {
			colourIdx = len(gradientColours) - 1
		}
		colour := gradientColours[colourIdx]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		bar.WriteString(style.Render("█"))
	}

	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
	for i := 0; i < emptyWidth; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}

	// Fixed-width star column (2 chars: space + star or two spaces)
	if isNewBest {
		starStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
		bar.WriteString(starStyle.Render(" *"))
	} else {
		bar.WriteString("  ")
	}

	return bar.String()
}

// formatStatRow creates a perfectly aligned row with label, value, and bar
func formatStatRow(label string, value string, bar string, labelWidth int, valueWidth int) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Width(labelWidth).
		Align(lipgloss.Right)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(valueWidth).
		Align(lipgloss.Right)

	return labelStyle.Render(label) + " " + valueStyle.Render(value) + " " + bar
}

func (m model) renderResults() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	newBestStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	sessionLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Width(18).
		Align(lipgloss.Right)

	sessionValueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(8).
		Align(lipgloss.Right)

	// Grid dimensions - must be consistent for perfect alignment
	const labelWidth = 18
	const valueWidth = 8
	const barWidth = 30
	const maxWPMDisplay = 120.0
	const maxTimeDisplay = 180.0 // 3 minutes max for time bar
	const maxAccuracy = 100.0

	title := titleStyle.Render("Round Complete!")

	// Check for new bests
	isNewBestWPM := m.stats.WPM >= m.historical.BestWPM
	isNewBestTime := m.historical.TotalSessions == 1 || m.stats.Duration.Seconds() <= m.historical.BestTime
	isNewBestAccuracy := m.stats.Accuracy >= m.historical.BestAccuracy

	var statsLines []string

	// WPM section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, formatStatRow(
		"WPM this run:", fmt.Sprintf("%.1f", m.stats.WPM),
		renderStatBar(m.stats.WPM, maxWPMDisplay, barWidth, isNewBestWPM),
		labelWidth, valueWidth))
	statsLines = append(statsLines, formatStatRow(
		"WPM best:", fmt.Sprintf("%.1f", m.historical.BestWPM),
		renderStatBar(m.historical.BestWPM, maxWPMDisplay, barWidth, false),
		labelWidth, valueWidth))
	statsLines = append(statsLines, formatStatRow(
		"WPM average:", fmt.Sprintf("%.1f", m.historical.AverageWPM()),
		renderStatBar(m.historical.AverageWPM(), maxWPMDisplay, barWidth, false),
		labelWidth, valueWidth))

	// Time section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, formatStatRow(
		"Time this run:", fmt.Sprintf("%.1fs", m.stats.Duration.Seconds()),
		renderTimeBar(m.stats.Duration.Seconds(), maxTimeDisplay, barWidth, isNewBestTime),
		labelWidth, valueWidth))
	statsLines = append(statsLines, formatStatRow(
		"Time best:", fmt.Sprintf("%.1fs", m.historical.BestTime),
		renderTimeBar(m.historical.BestTime, maxTimeDisplay, barWidth, false),
		labelWidth, valueWidth))
	statsLines = append(statsLines, formatStatRow(
		"Time average:", fmt.Sprintf("%.1fs", m.historical.AverageTime()),
		renderTimeBar(m.historical.AverageTime(), maxTimeDisplay, barWidth, false),
		labelWidth, valueWidth))

	// Accuracy section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, formatStatRow(
		"Accuracy this run:", fmt.Sprintf("%.1f%%", m.stats.Accuracy),
		renderStatBar(m.stats.Accuracy, maxAccuracy, barWidth, isNewBestAccuracy),
		labelWidth, valueWidth))
	statsLines = append(statsLines, formatStatRow(
		"Accuracy best:", fmt.Sprintf("%.1f%%", m.historical.BestAccuracy),
		renderStatBar(m.historical.BestAccuracy, maxAccuracy, barWidth, false),
		labelWidth, valueWidth))
	statsLines = append(statsLines, formatStatRow(
		"Accuracy average:", fmt.Sprintf("%.1f%%", m.historical.AverageAccuracy()),
		renderStatBar(m.historical.AverageAccuracy(), maxAccuracy, barWidth, false),
		labelWidth, valueWidth))

	// Sessions - use same grid alignment
	statsLines = append(statsLines, "")
	statsLines = append(statsLines,
		sessionLabelStyle.Render("Total sessions:")+" "+sessionValueStyle.Render(fmt.Sprintf("%d", m.historical.TotalSessions)))

	// Legend
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, newBestStyle.Render("* = New personal best!"))

	help := helpStyle.Render("Press ENTER for a new round | ESC to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		strings.Join(statsLines, "\n"),
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
