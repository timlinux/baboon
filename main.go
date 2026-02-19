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

func (m model) renderResults() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true).
		MarginBottom(1)

	statLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	statValueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	comparisonBetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	comparisonWorseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		MarginTop(2)

	// Build stats display
	title := titleStyle.Render("Round Complete!")

	wpmStr := fmt.Sprintf("%.1f", m.stats.WPM)
	accuracyStr := fmt.Sprintf("%.1f%%", m.stats.Accuracy)
	durationStr := fmt.Sprintf("%.1f seconds", m.stats.Duration.Seconds())

	statsContent := []string{
		statLabelStyle.Render("Words Per Minute: ") + statValueStyle.Render(wpmStr),
		statLabelStyle.Render("Accuracy: ") + statValueStyle.Render(accuracyStr),
		statLabelStyle.Render("Time: ") + statValueStyle.Render(durationStr),
		statLabelStyle.Render("Characters Typed: ") + statValueStyle.Render(fmt.Sprintf("%d", m.stats.TotalCharacters)),
		"",
	}

	// Historical comparison
	if m.historical.TotalSessions > 1 {
		statsContent = append(statsContent, statLabelStyle.Render("--- Historical Best ---"))

		wpmDiff := m.stats.WPM - m.historical.BestWPM
		var wpmComparison string
		if wpmDiff >= 0 {
			wpmComparison = comparisonBetterStyle.Render(fmt.Sprintf(" (NEW BEST! +%.1f)", wpmDiff))
		} else {
			wpmComparison = comparisonWorseStyle.Render(fmt.Sprintf(" (%.1f from best)", wpmDiff))
		}
		statsContent = append(statsContent, statLabelStyle.Render("Best WPM: ")+statValueStyle.Render(fmt.Sprintf("%.1f", m.historical.BestWPM))+wpmComparison)

		accDiff := m.stats.Accuracy - m.historical.BestAccuracy
		var accComparison string
		if accDiff >= 0 {
			accComparison = comparisonBetterStyle.Render(fmt.Sprintf(" (NEW BEST! +%.1f%%)", accDiff))
		} else {
			accComparison = comparisonWorseStyle.Render(fmt.Sprintf(" (%.1f%% from best)", accDiff))
		}
		statsContent = append(statsContent, statLabelStyle.Render("Best Accuracy: ")+statValueStyle.Render(fmt.Sprintf("%.1f%%", m.historical.BestAccuracy))+accComparison)

		statsContent = append(statsContent, statLabelStyle.Render("Total Sessions: ")+statValueStyle.Render(fmt.Sprintf("%d", m.historical.TotalSessions)))
	} else {
		statsContent = append(statsContent, comparisonBetterStyle.Render("First session complete! Your scores are now the benchmark."))
	}

	help := helpStyle.Render("Press ENTER for a new round | ESC to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		strings.Join(statsContent, "\n"),
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
