package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/timlinux/baboon/font"
	"github.com/timlinux/baboon/stats"
	"github.com/timlinux/baboon/words"
)

// Punctuation characters used in punctuation mode
var punctuationChars = []string{",", ".", ";", ":", "!", "?"}

// Round configuration
const (
	wordsPerRound      = 30  // Fixed number of words per round
	charactersPerRound = 150 // Fixed total characters per round
)

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

// animTickMsg is sent to update animations
type animTickMsg time.Time

// Animation configuration
const (
	numAnimatedRows   = 14 // Number of rows to animate on results screen
	animationInterval = 50 * time.Millisecond
	staggerDelay      = 3 // Frames between each row starting
)

type model struct {
	state           gameState
	words           []string
	currentWordIdx  int
	currentInput    string
	stats           *stats.Stats
	historical      *stats.HistoricalStats
	width           int
	height          int
	rng             *rand.Rand
	timerStarted    bool
	punctuationMode bool // When true, words are separated by punctuation + space

	// Animation state for results screen
	animSprings    []harmonica.Spring
	animPositions  []float64
	animVelocities []float64
	animFrame      int
}

// getLetterData extracts letter frequency and accuracy data from historical stats
func getLetterData(historical *stats.HistoricalStats) words.LetterData {
	data := make(words.LetterData)
	if historical == nil || historical.LetterAccuracy == nil {
		return data
	}
	for letter, letterStats := range historical.LetterAccuracy {
		data[letter] = words.LetterStats{
			Presented: letterStats.Presented,
			Correct:   letterStats.Correct,
		}
	}
	return data
}

func initialModel(punctuationMode bool) model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	historical, _ := stats.LoadHistoricalStats()

	// Get letter data for weighted word selection (frequency + accuracy)
	letterData := getLetterData(historical)
	wordList := words.GetRandomWordsFixedCount(wordsPerRound, charactersPerRound, rng.Intn, letterData)
	sessionStats := &stats.Stats{
		LetterAccuracy: make(map[string]stats.LetterStats),
		LetterSeekTime: make(map[string]stats.LetterSeekStats),
		BigramSeekTime: make(map[string]stats.BigramSeekStats),
	}

	// Record all letters in all words as presented (only a-z, not punctuation)
	for _, word := range wordList {
		for _, char := range word {
			if char >= 'a' && char <= 'z' {
				sessionStats.RecordLetterPresented(string(char))
			}
		}
	}

	// In punctuation mode, append random punctuation to each word (except the last)
	if punctuationMode {
		for i := 0; i < len(wordList)-1; i++ {
			punct := punctuationChars[rng.Intn(len(punctuationChars))]
			wordList[i] = wordList[i] + punct
		}
	}

	m := model{
		state:           stateTyping,
		words:           wordList,
		historical:      historical,
		rng:             rng,
		timerStarted:    false,
		stats:           sessionStats,
		punctuationMode: punctuationMode,
	}
	return m
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func animTickCmd() tea.Cmd {
	return tea.Tick(animationInterval, func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}

// initAnimations sets up springs for animating results rows
func (m *model) initAnimations() {
	m.animSprings = make([]harmonica.Spring, numAnimatedRows)
	m.animPositions = make([]float64, numAnimatedRows)
	m.animVelocities = make([]float64, numAnimatedRows)
	m.animFrame = 0

	// Create springs with a bouncy feel
	for i := range m.animSprings {
		m.animSprings[i] = harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5)
		m.animPositions[i] = 0 // Start at 0 (invisible)
		m.animVelocities[i] = 0
	}
}

// updateAnimations advances all active springs
func (m *model) updateAnimations() {
	m.animFrame++

	for i := range m.animSprings {
		// Stagger the start: each row starts staggerDelay frames after the previous
		startFrame := i * staggerDelay
		if m.animFrame >= startFrame {
			// Target is 1.0 (fully visible)
			m.animPositions[i], m.animVelocities[i] = m.animSprings[i].Update(
				m.animPositions[i], m.animVelocities[i], 1.0,
			)
		}
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Continue ticking to update WPM display
		return m, tickCmd()

	case animTickMsg:
		// Update animations on results screen
		if m.state == stateResults {
			m.updateAnimations()
			return m, animTickCmd()
		}
		return m, nil

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
		// Only advance to next word if all letters have been typed
		if len(m.currentInput) >= len(currentWord) {
			m.stats.WordsCompleted++
			m.currentInput = ""
			m.currentWordIdx++
			// Reset last letter for bigram tracking (new word = fresh start)
			m.stats.LastLetter = ""

			if m.currentWordIdx >= len(m.words) {
				// Round complete - show results
				m.stats.Calculate()
				m.historical.UpdateHistorical(m.stats)
				stats.SaveHistoricalStats(m.historical)
				m.state = stateResults
				m.initAnimations()
				return m, animTickCmd()
			}
		} else if len(m.currentInput) > 0 || m.timerStarted {
			// Treat space as an incorrect character if word is not complete
			// Don't record seek time for incorrect keystrokes
			m.stats.LastKeyTime = time.Now()
			m.currentInput += " "
			m.stats.TotalCharacters++
			m.stats.IncorrectChars++
		}

	case tea.KeyBackspace:
		if len(m.currentInput) > 0 {
			m.currentInput = m.currentInput[:len(m.currentInput)-1]
		}

	case tea.KeyRunes:
		char := string(msg.Runes)
		inputIdx := len(m.currentInput)
		now := time.Now()

		// Start timer on first correct character of first word
		if !m.timerStarted && m.currentWordIdx == 0 && inputIdx == 0 {
			if len(currentWord) > 0 && char == string(currentWord[0]) {
				m.timerStarted = true
				m.stats.StartTime = now
				m.stats.LastKeyTime = now
			}
		}

		m.currentInput += char
		m.stats.TotalCharacters++

		// Check if character matches
		isCorrect := inputIdx < len(currentWord) && m.currentInput[inputIdx] == currentWord[inputIdx]
		if isCorrect {
			m.stats.CorrectChars++
			expectedChar := currentWord[inputIdx]
			expectedLetter := string(expectedChar)

			// Only record letter stats for actual letters (a-z), not punctuation
			isLetter := expectedChar >= 'a' && expectedChar <= 'z'
			if isLetter {
				// Record letter as correctly typed
				m.stats.RecordLetterCorrect(expectedLetter)

				// Record seek time only for correct keystrokes
				// Exclude first letter of each word (includes word-reading time)
				if m.timerStarted && inputIdx > 0 && !m.stats.LastKeyTime.IsZero() {
					seekTimeMs := now.Sub(m.stats.LastKeyTime).Milliseconds()
					// Only record reasonable seek times (< 5 seconds to filter pauses)
					if seekTimeMs > 0 && seekTimeMs < 5000 {
						// Record against expected letter, not typed letter
						m.stats.RecordLetterSeekTime(expectedLetter, seekTimeMs)

						// Record bigram timing (letter pair)
						if m.stats.LastLetter != "" {
							bigram := m.stats.LastLetter + expectedLetter
							m.stats.RecordBigramSeekTime(bigram, seekTimeMs)
						}
					}
				}
				// Update last letter for bigram tracking (only for correct letters)
				m.stats.LastLetter = expectedLetter
			}
		} else {
			m.stats.IncorrectChars++
			// Don't update LastLetter for incorrect keystrokes
		}
		m.stats.LastKeyTime = now
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
		// Get letter data for weighted word selection (use updated historical stats)
		letterData := getLetterData(m.historical)
		m.words = words.GetRandomWordsFixedCount(wordsPerRound, charactersPerRound, m.rng.Intn, letterData)
		m.currentWordIdx = 0
		m.currentInput = ""
		m.timerStarted = false
		m.stats = &stats.Stats{
			LetterAccuracy: make(map[string]stats.LetterStats),
			LetterSeekTime: make(map[string]stats.LetterSeekStats),
			BigramSeekTime: make(map[string]stats.BigramSeekStats),
		}
		// Record all letters in all words as presented (before adding punctuation)
		for _, word := range m.words {
			for _, char := range word {
				if char >= 'a' && char <= 'z' {
					m.stats.RecordLetterPresented(string(char))
				}
			}
		}
		// In punctuation mode, append random punctuation to each word (except the last)
		if m.punctuationMode {
			for i := 0; i < len(m.words)-1; i++ {
				punct := punctuationChars[m.rng.Intn(len(punctuationChars))]
				m.words[i] = m.words[i] + punct
			}
		}
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

	// Previous word (top left) and Next word (top right)
	prevNextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	var prevWord, nextWord string
	if m.currentWordIdx > 0 {
		prevWord = m.words[m.currentWordIdx-1]
	}
	if m.currentWordIdx < len(m.words)-1 {
		nextWord = m.words[m.currentWordIdx+1]
	}

	// Create top bar with prev/next words
	prevLabel := prevNextStyle.Render(prevWord)
	nextLabel := prevNextStyle.Render(nextWord)

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
	topPadding := (m.height - contentHeight - wpmBarHeight - 6) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Build full screen layout
	var fullContent strings.Builder

	// Top bar with previous and next words
	topBarWidth := m.width - 4 // Leave some margin
	if topBarWidth < 20 {
		topBarWidth = 20
	}
	spaceBetween := topBarWidth - len(prevWord) - len(nextWord)
	if spaceBetween < 1 {
		spaceBetween = 1
	}
	topBar := "  " + prevLabel + strings.Repeat(" ", spaceBetween) + nextLabel + "  "
	fullContent.WriteString(topBar)
	fullContent.WriteString("\n")

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
	currentHeight := topPadding + contentHeight + 2 + wpmBarHeight + 2
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

// getAccuracyColour returns a colour on red-yellow-green gradient based on accuracy (0-100)
func getAccuracyColour(accuracy float64) string {
	// Gradient from red (0%) through yellow (50%) to green (100%)
	// Using 256-colour palette
	if accuracy >= 95 {
		return "46" // Bright green
	} else if accuracy >= 90 {
		return "82"
	} else if accuracy >= 85 {
		return "118"
	} else if accuracy >= 80 {
		return "154"
	} else if accuracy >= 75 {
		return "190"
	} else if accuracy >= 70 {
		return "226" // Yellow
	} else if accuracy >= 65 {
		return "220"
	} else if accuracy >= 60 {
		return "214"
	} else if accuracy >= 50 {
		return "208"
	} else if accuracy >= 40 {
		return "202"
	}
	return "196" // Red
}

// getFrequencyColour returns a colour on red-yellow-green gradient based on frequency (0-1)
func getFrequencyColour(frequency float64) string {
	// Same gradient as accuracy
	return getAccuracyColour(frequency * 100)
}

// renderLetterHeaderRow renders a row of 26 letters as column headers
func (m model) renderLetterHeaderRow() string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true)

	for i, letter := range letters {
		row.WriteString(headerStyle.Render(string(letter)))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderLetterAccuracyRow renders a row of 26 filled circles coloured by accuracy
func (m model) renderLetterAccuracyRow() string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i, letter := range letters {
		lowerLetter := string(letter + 32) // Convert to lowercase for lookup
		letterStats := m.historical.LetterAccuracy[lowerLetter]

		var accuracy float64
		if letterStats.Presented > 0 {
			accuracy = (float64(letterStats.Correct) / float64(letterStats.Presented)) * 100
		}

		colour := getAccuracyColour(accuracy)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

		row.WriteString(style.Render("●"))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderLetterFrequencyRow renders a row of 26 filled circles coloured by frequency
func (m model) renderLetterFrequencyRow() string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Find max presented count for normalization
	var maxPresented int
	for _, letter := range letters {
		lowerLetter := string(letter + 32)
		if count := m.historical.LetterAccuracy[lowerLetter].Presented; count > maxPresented {
			maxPresented = count
		}
	}

	for i, letter := range letters {
		lowerLetter := string(letter + 32)
		letterStats := m.historical.LetterAccuracy[lowerLetter]

		var frequency float64
		if maxPresented > 0 {
			frequency = float64(letterStats.Presented) / float64(maxPresented)
		}

		colour := getFrequencyColour(frequency)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

		row.WriteString(style.Render("●"))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// getSeekTimeColour returns a colour based on seek time (inverted: fast=green, slow=red)
func getSeekTimeColour(seekTimeMs, maxSeekTimeMs float64) string {
	if maxSeekTimeMs == 0 {
		return "46" // Green if no data
	}
	// Normalize and invert (fast = low time = green, slow = high time = red)
	normalized := seekTimeMs / maxSeekTimeMs
	// Invert: low time should be green (high accuracy equivalent)
	accuracy := (1.0 - normalized) * 100
	return getAccuracyColour(accuracy)
}

// renderLetterSeekTimeRow renders a row of 26 filled circles coloured by seek time
func (m model) renderLetterSeekTimeRow() string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Find max average seek time for normalization
	var maxSeekTime float64
	for _, letter := range letters {
		lowerLetter := string(letter + 32)
		if seekStats, exists := m.historical.LetterSeekTime[lowerLetter]; exists {
			avgTime := seekStats.AverageMs()
			if avgTime > maxSeekTime {
				maxSeekTime = avgTime
			}
		}
	}

	for i, letter := range letters {
		lowerLetter := string(letter + 32)
		seekStats := m.historical.LetterSeekTime[lowerLetter]
		avgTime := seekStats.AverageMs()

		colour := getSeekTimeColour(avgTime, maxSeekTime)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

		row.WriteString(style.Render("●"))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// applyRowAnimation applies slide-in animation to a row based on position (0-1)
func (m model) applyRowAnimation(row string, animIdx int) string {
	if animIdx >= len(m.animPositions) {
		return row
	}

	pos := m.animPositions[animIdx]
	if pos >= 0.99 {
		return row // Fully visible
	}
	if pos <= 0.01 {
		return "" // Hidden
	}

	// Slide in from right: offset decreases as position increases
	maxOffset := 40
	offset := int(float64(maxOffset) * (1.0 - pos))
	if offset > 0 {
		return strings.Repeat(" ", offset) + row
	}
	return row
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

	// Build animated rows
	animIdx := 0
	var statsLines []string

	// WPM section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"WPM this run:", fmt.Sprintf("%.1f", m.stats.WPM),
		renderStatBar(m.stats.WPM, maxWPMDisplay, barWidth, isNewBestWPM),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"WPM best:", fmt.Sprintf("%.1f", m.historical.BestWPM),
		renderStatBar(m.historical.BestWPM, maxWPMDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"WPM average:", fmt.Sprintf("%.1f", m.historical.AverageWPM()),
		renderStatBar(m.historical.AverageWPM(), maxWPMDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++

	// Time section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"Time this run:", fmt.Sprintf("%.1fs", m.stats.Duration.Seconds()),
		renderTimeBar(m.stats.Duration.Seconds(), maxTimeDisplay, barWidth, isNewBestTime),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"Time best:", fmt.Sprintf("%.1fs", m.historical.BestTime),
		renderTimeBar(m.historical.BestTime, maxTimeDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"Time average:", fmt.Sprintf("%.1fs", m.historical.AverageTime()),
		renderTimeBar(m.historical.AverageTime(), maxTimeDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++

	// Accuracy section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"Accuracy this run:", fmt.Sprintf("%.1f%%", m.stats.Accuracy),
		renderStatBar(m.stats.Accuracy, maxAccuracy, barWidth, isNewBestAccuracy),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"Accuracy best:", fmt.Sprintf("%.1f%%", m.historical.BestAccuracy),
		renderStatBar(m.historical.BestAccuracy, maxAccuracy, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(formatStatRow(
		"Accuracy average:", fmt.Sprintf("%.1f%%", m.historical.AverageAccuracy()),
		renderStatBar(m.historical.AverageAccuracy(), maxAccuracy, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++

	// Sessions - use same grid alignment
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, m.applyRowAnimation(
		sessionLabelStyle.Render("Total sessions:")+" "+sessionValueStyle.Render(fmt.Sprintf("%d", m.historical.TotalSessions)),
		animIdx))
	animIdx++

	// Letter statistics matrix
	letterLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Width(labelWidth).
		Align(lipgloss.Right)
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, m.applyRowAnimation(
		letterLabelStyle.Render("")+" "+m.renderLetterHeaderRow(), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(
		letterLabelStyle.Render("Accuracy:")+" "+m.renderLetterAccuracyRow(), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(
		letterLabelStyle.Render("Frequency:")+" "+m.renderLetterFrequencyRow(), animIdx))
	animIdx++
	statsLines = append(statsLines, m.applyRowAnimation(
		letterLabelStyle.Render("Seek time:")+" "+m.renderLetterSeekTimeRow(), animIdx))

	// Legend (only show if user achieved a personal best)
	if isNewBestWPM || isNewBestTime || isNewBestAccuracy {
		statsLines = append(statsLines, "")
		statsLines = append(statsLines, newBestStyle.Render("* = New personal best!"))
	}

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
	punctuationMode := flag.Bool("p", false, "Enable punctuation mode (words separated by punctuation + space)")
	flag.Parse()

	p := tea.NewProgram(initialModel(*punctuationMode), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
