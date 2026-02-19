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

// Game states
type gameState int

const (
	stateTyping gameState = iota
	stateResults
)

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
}

func initialModel() model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	historical, _ := stats.LoadHistoricalStats()

	m := model{
		state:      stateTyping,
		words:      words.GetRandomWords(wordsPerRound, rng.Intn),
		historical: historical,
		rng:        rng,
		stats: &stats.Stats{
			StartTime: time.Now(),
		},
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		m.currentInput += char
		m.stats.TotalCharacters++

		// Check if character matches
		inputIdx := len(m.currentInput) - 1
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
		m.stats = &stats.Stats{
			StartTime: time.Now(),
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
	help := helpStyle.Render("Type the word, then press SPACE to continue | ESC to quit")

	// Center everything
	content := lipgloss.JoinVertical(
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

	// Center on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
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
