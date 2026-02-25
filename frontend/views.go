package frontend

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/timlinux/baboon/backend"
	"github.com/timlinux/baboon/font"
	"github.com/timlinux/baboon/settings"
	"github.com/timlinux/baboon/stats"
)

// Renderer handles all view rendering for the application
type Renderer struct {
	styles Styles
	width  int
	height int
}

// NewRenderer creates a new renderer with the given dimensions
func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		styles: NewStyles(),
		width:  width,
		height: height,
	}
}

// SetSize updates the renderer dimensions
func (r *Renderer) SetSize(width, height int) {
	r.width = width
	r.height = height
}

// RenderTypingScreenAnimated renders the main typing interface with smooth carousel animations
func (r *Renderer) RenderTypingScreenAnimated(state backend.GameState, carousel *CarouselAnimator, s *settings.Settings) string {
	if state.CurrentWordIdx >= len(state.Words) {
		return ""
	}

	currentWord := state.CurrentWord
	if len(currentWord) == 0 {
		return "Loading..."
	}

	// Render current word using custom block font
	letterLines := font.RenderWord(currentWord)

	// Build colored output for each line
	coloredLines := make([]string, font.LetterHeight)

	for lineIdx := 0; lineIdx < font.LetterHeight; lineIdx++ {
		var lineBuilder strings.Builder

		for charIdx, letterLine := range letterLines[lineIdx] {
			var style lipgloss.Style

			if charIdx < len(state.CurrentInput) {
				// Character has been typed
				if charIdx < len(currentWord) && state.CurrentInput[charIdx] == currentWord[charIdx] {
					style = r.styles.Correct
				} else {
					style = r.styles.Incorrect
				}
			} else {
				// Character not yet typed
				style = r.styles.Untyped
			}

			lineBuilder.WriteString(style.Render(letterLine))
			if charIdx < len(letterLines[lineIdx])-1 {
				lineBuilder.WriteString(style.Render(" "))
			}
		}

		coloredLines[lineIdx] = lineBuilder.String()
	}

	coloredWord := strings.Join(coloredLines, "\n")

	// Progress indicator
	progress := fmt.Sprintf("Word %d/%d", state.WordNumber, state.TotalWords)

	// Get animation values (default to fully visible if no animator)
	prevOpacity := 0.5
	currentOffset := 0
	nextOpacity := 0.6
	if carousel != nil {
		prevOpacity = carousel.GetPrevOpacity()
		currentOffset = carousel.GetCurrentOffset()
		nextOpacity = carousel.GetNextOpacity()
	}

	// Carousel style: Previous word above (animated opacity via colour intensity)
	prevWordDisplay := ""
	if state.PreviousWord != "" {
		// Map opacity to greyscale colour (232-255 range in 256-colour palette)
		// Lower opacity = darker colour
		greyLevel := 232 + int(prevOpacity*23) // 232 (darkest) to 255 (lightest)
		if greyLevel > 255 {
			greyLevel = 255
		}
		prevStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("%d", greyLevel))).
			Italic(true)
		prevWordDisplay = prevStyle.Render("· · · " + state.PreviousWord + " · · ·")
	}

	// Carousel style: Next words below (up to 3, with decreasing opacity)
	var nextWordsDisplay []string
	// Use NextWords if available, otherwise fall back to NextWord for backwards compatibility
	nextWords := state.NextWords
	if len(nextWords) == 0 && state.NextWord != "" {
		nextWords = []string{state.NextWord}
	}
	for i, word := range nextWords {
		// Decrease opacity for words further ahead
		wordOpacity := nextOpacity * (1.0 - float64(i)*0.2)
		if wordOpacity < 0.2 {
			wordOpacity = 0.2
		}
		greyLevel := 232 + int(wordOpacity*23)
		if greyLevel > 255 {
			greyLevel = 255
		}
		nextStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("%d", greyLevel))).
			Align(lipgloss.Center)
		if i == 0 {
			nextWordsDisplay = append(nextWordsDisplay, nextStyle.Render("▼  "+word+"  ▼"))
		} else {
			nextWordsDisplay = append(nextWordsDisplay, nextStyle.Render(word))
		}
	}

	// WPM Bar
	wpmBar := r.renderWPMBar(state.LiveWPM)

	// Build the carousel layout vertically (main content only)
	var carouselElements []string

	// Progress at top of main content
	carouselElements = append(carouselElements, r.styles.Progress.Render(progress))
	carouselElements = append(carouselElements, "")

	// Previous word (above current, animated)
	if prevWordDisplay != "" {
		// Add vertical offset lines for animation (previous scrolls up)
		prevOffset := 0
		if carousel != nil {
			prevOffset = carousel.GetPrevOffset()
		}
		for i := 0; i < prevOffset; i++ {
			carouselElements = append(carouselElements, "")
		}
		carouselElements = append(carouselElements, prevWordDisplay)
		carouselElements = append(carouselElements, "")
	}

	// Decorative separator before main word
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	carouselElements = append(carouselElements, separatorStyle.Render("─────────────────────────────────"))

	// Add offset lines for current word animation (slides up from below)
	for i := 0; i < currentOffset; i++ {
		carouselElements = append(carouselElements, "")
	}
	carouselElements = append(carouselElements, "")

	// Current word (large block letters) - the main focus
	carouselElements = append(carouselElements, coloredWord)

	// Decorative separator after main word
	carouselElements = append(carouselElements, "")
	carouselElements = append(carouselElements, separatorStyle.Render("─────────────────────────────────"))

	// Next words (below current, animated) - show up to 3 upcoming words
	if len(nextWordsDisplay) > 0 {
		nextOffset := 0
		if carousel != nil {
			nextOffset = carousel.GetNextOffset()
		}
		for i := 0; i < nextOffset; i++ {
			carouselElements = append(carouselElements, "")
		}
		carouselElements = append(carouselElements, "")
		for _, nextWord := range nextWordsDisplay {
			carouselElements = append(carouselElements, nextWord)
		}
	}

	carouselElements = append(carouselElements, "")
	carouselElements = append(carouselElements, wpmBar)

	// Center the main content horizontally
	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		carouselElements...,
	)

	// Fixed header at top
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)
	header := lipgloss.PlaceHorizontal(r.width, lipgloss.Center,
		headerStyle.Render("🐒 BABOON - Typing Practice"))

	// Fixed footer at bottom
	var helpText string
	advanceKeyHint := "SPACE"
	if s != nil {
		advanceKeyHint = s.AdvanceKey.KeyHint()
	}
	if !state.TimerStarted {
		helpText = "Type the first letter to start | Tab to restart | 'o' for options | ESC to quit"
	} else {
		helpText = fmt.Sprintf("Type the word, then press %s to continue | Tab to restart | ESC to quit", advanceKeyHint)
	}
	footer := lipgloss.PlaceHorizontal(r.width, lipgloss.Center, r.styles.Help.Render(helpText))

	// Calculate heights
	headerHeight := 1
	footerHeight := 1
	contentHeight := strings.Count(mainContent, "\n") + 1
	availableHeight := r.height - headerHeight - footerHeight - 2 // -2 for spacing

	// Calculate top padding to center main content in available space
	topPadding := (availableHeight - contentHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Build full screen layout
	var fullContent strings.Builder

	// Header at top (line 0)
	fullContent.WriteString(header)
	fullContent.WriteString("\n")

	// Top padding to center content
	for i := 0; i < topPadding; i++ {
		fullContent.WriteString("\n")
	}

	// Main content (centered horizontally)
	centeredMain := lipgloss.PlaceHorizontal(r.width, lipgloss.Center, mainContent)
	fullContent.WriteString(centeredMain)

	// Bottom padding to push footer to the bottom
	currentHeight := headerHeight + 1 + topPadding + contentHeight
	for i := currentHeight; i < r.height-footerHeight; i++ {
		fullContent.WriteString("\n")
	}

	// Footer at bottom (last line)
	fullContent.WriteString(footer)

	return fullContent.String()
}

// renderWPMBar creates a beautiful gradient progress bar for WPM
func (r *Renderer) renderWPMBar(wpm float64) string {
	const maxWPM = 120.0
	const barWidth = 50

	fillPercent := wpm / maxWPM
	if fillPercent > 1.0 {
		fillPercent = 1.0
	}
	if fillPercent < 0 {
		fillPercent = 0
	}

	filledWidth := int(float64(barWidth) * fillPercent)
	emptyWidth := barWidth - filledWidth

	var bar strings.Builder

	// Build the filled portion with gradient
	for i := 0; i < filledWidth; i++ {
		colourIdx := int(float64(i) / float64(barWidth) * float64(len(GradientColours)-1))
		if colourIdx >= len(GradientColours) {
			colourIdx = len(GradientColours) - 1
		}
		colour := GradientColours[colourIdx]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		bar.WriteString(style.Render("█"))
	}

	// Build the empty portion
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourEmptyBar))
	for i := 0; i < emptyWidth; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}

	// WPM label
	wpmLabel := fmt.Sprintf(" %.0f WPM", wpm)
	var labelStyle lipgloss.Style
	if wpm >= 60 {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	} else if wpm >= 40 {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	} else {
		labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
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

// RenderResultsScreen renders the results screen with animations
func (r *Renderer) RenderResultsScreen(
	session *stats.Stats,
	historical *stats.HistoricalStats,
	animator *Animator,
) string {
	const labelWidth = 18
	const valueWidth = 8
	const barWidth = 30
	const maxWPMDisplay = 120.0
	const maxTimeDisplay = 180.0
	const maxAccuracy = 100.0

	title := r.styles.Title.Render("Round Complete!")

	// Check for new bests
	isNewBestWPM := session.WPM >= historical.BestWPM
	isNewBestTime := historical.TotalSessions == 1 || session.Duration.Seconds() <= historical.BestTime
	isNewBestAccuracy := session.Accuracy >= historical.BestAccuracy

	// Build animated rows
	animIdx := 0
	var statsLines []string

	// WPM section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"WPM this run:", fmt.Sprintf("%.1f", session.WPM),
		r.renderStatBar(session.WPM, maxWPMDisplay, barWidth, isNewBestWPM),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"WPM best:", fmt.Sprintf("%.1f", historical.BestWPM),
		r.renderStatBar(historical.BestWPM, maxWPMDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"WPM average:", fmt.Sprintf("%.1f", historical.AverageWPM()),
		r.renderStatBar(historical.AverageWPM(), maxWPMDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++

	// Time section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"Time this run:", fmt.Sprintf("%.1fs", session.Duration.Seconds()),
		r.renderTimeBar(session.Duration.Seconds(), maxTimeDisplay, barWidth, isNewBestTime),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"Time best:", fmt.Sprintf("%.1fs", historical.BestTime),
		r.renderTimeBar(historical.BestTime, maxTimeDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"Time average:", fmt.Sprintf("%.1fs", historical.AverageTime()),
		r.renderTimeBar(historical.AverageTime(), maxTimeDisplay, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++

	// Accuracy section
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"Accuracy this run:", fmt.Sprintf("%.1f%%", session.Accuracy),
		r.renderStatBar(session.Accuracy, maxAccuracy, barWidth, isNewBestAccuracy),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"Accuracy best:", fmt.Sprintf("%.1f%%", historical.BestAccuracy),
		r.renderStatBar(historical.BestAccuracy, maxAccuracy, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.formatStatRow(
		"Accuracy average:", fmt.Sprintf("%.1f%%", historical.AverageAccuracy()),
		r.renderStatBar(historical.AverageAccuracy(), maxAccuracy, barWidth, false),
		labelWidth, valueWidth), animIdx))
	animIdx++

	// Sessions
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, animator.ApplyAnimation(
		r.styles.SessionLabel.Render("Total sessions:")+" "+r.styles.SessionValue.Render(fmt.Sprintf("%d", historical.TotalSessions)),
		animIdx))
	animIdx++

	// Typing theory stats
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, animator.ApplyAnimation(r.renderFingerRow(historical, labelWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.renderRowAccuracyRow(historical, labelWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.renderHandStats(historical, labelWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.renderRhythmStats(session, historical, labelWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.renderSFBStats(session, historical, labelWidth), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(r.renderTopErrors(historical, labelWidth), animIdx))
	animIdx++

	// Letter statistics matrix
	statsLines = append(statsLines, "")
	statsLines = append(statsLines, animator.ApplyAnimation(
		r.styles.LetterLabel.Render("")+" "+r.renderLetterHeaderRow(), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(
		r.styles.LetterLabel.Render("Accuracy:")+" "+r.renderLetterAccuracyRow(historical), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(
		r.styles.LetterLabel.Render("Frequency:")+" "+r.renderLetterFrequencyRow(historical), animIdx))
	animIdx++
	statsLines = append(statsLines, animator.ApplyAnimation(
		r.styles.LetterLabel.Render("Seek time:")+" "+r.renderLetterSeekTimeRow(historical), animIdx))

	// Legend (only show if user achieved a personal best)
	if isNewBestWPM || isNewBestTime || isNewBestAccuracy {
		statsLines = append(statsLines, "")
		statsLines = append(statsLines, r.styles.NewBest.Render("* = New personal best!"))
	}

	// Main content (title + stats)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		strings.Join(statsLines, "\n"),
	)

	// Fixed header at top
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)
	header := lipgloss.PlaceHorizontal(r.width, lipgloss.Center,
		headerStyle.Render("🐒 BABOON - Typing Practice"))

	// Fixed footer at bottom
	footer := lipgloss.PlaceHorizontal(r.width, lipgloss.Center,
		r.styles.Help.Render("Press ENTER for a new round | 'o' for options | ESC to quit"))

	// Calculate heights
	headerHeight := 1
	footerHeight := 1
	contentHeight := strings.Count(mainContent, "\n") + 1
	availableHeight := r.height - headerHeight - footerHeight - 2 // -2 for spacing

	// Calculate top padding to center main content in available space
	topPadding := (availableHeight - contentHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Build full screen layout
	var fullContent strings.Builder

	// Header at top (line 0)
	fullContent.WriteString(header)
	fullContent.WriteString("\n")

	// Top padding to center content
	for i := 0; i < topPadding; i++ {
		fullContent.WriteString("\n")
	}

	// Main content (centered horizontally)
	centeredMain := lipgloss.PlaceHorizontal(r.width, lipgloss.Center, mainContent)
	fullContent.WriteString(centeredMain)

	// Bottom padding to push footer to the bottom
	currentHeight := headerHeight + 1 + topPadding + contentHeight
	for i := currentHeight; i < r.height-footerHeight; i++ {
		fullContent.WriteString("\n")
	}

	// Footer at bottom (last line)
	fullContent.WriteString(footer)

	return fullContent.String()
}

// formatStatRow creates a perfectly aligned row with label, value, and bar
func (r *Renderer) formatStatRow(label string, value string, bar string, labelWidth int, valueWidth int) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourLabel)).
		Width(labelWidth).
		Align(lipgloss.Right)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourValue)).
		Width(valueWidth).
		Align(lipgloss.Right)

	return labelStyle.Render(label) + " " + valueStyle.Render(value) + " " + bar
}

// renderStatBar creates a gradient bar for a stat value
func (r *Renderer) renderStatBar(value, maxValue float64, width int, isNewBest bool) string {
	fillPercent := value / maxValue
	if fillPercent > 1.0 {
		fillPercent = 1.0
	}
	if fillPercent < 0 {
		fillPercent = 0
	}

	filledWidth := int(float64(width) * fillPercent)
	emptyWidth := width - filledWidth

	var bar strings.Builder

	for i := 0; i < filledWidth; i++ {
		colourIdx := int(float64(i) / float64(width) * float64(len(GradientColours)-1))
		if colourIdx >= len(GradientColours) {
			colourIdx = len(GradientColours) - 1
		}
		colour := GradientColours[colourIdx]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		bar.WriteString(style.Render("█"))
	}

	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourEmptyBar))
	for i := 0; i < emptyWidth; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}

	if isNewBest {
		starStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourNewBest)).Bold(true)
		bar.WriteString(starStyle.Render(" *"))
	} else {
		bar.WriteString("  ")
	}

	return bar.String()
}

// renderTimeBar creates a bar for time (lower is better, so inverted)
func (r *Renderer) renderTimeBar(value, maxValue float64, width int, isNewBest bool) string {
	fillPercent := 1.0 - (value / maxValue)
	if fillPercent > 1.0 {
		fillPercent = 1.0
	}
	if fillPercent < 0 {
		fillPercent = 0
	}

	filledWidth := int(float64(width) * fillPercent)
	emptyWidth := width - filledWidth

	var bar strings.Builder

	for i := 0; i < filledWidth; i++ {
		colourIdx := int(float64(i) / float64(width) * float64(len(GradientColours)-1))
		if colourIdx >= len(GradientColours) {
			colourIdx = len(GradientColours) - 1
		}
		colour := GradientColours[colourIdx]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		bar.WriteString(style.Render("█"))
	}

	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourEmptyBar))
	for i := 0; i < emptyWidth; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}

	if isNewBest {
		starStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourNewBest)).Bold(true)
		bar.WriteString(starStyle.Render(" *"))
	} else {
		bar.WriteString("  ")
	}

	return bar.String()
}

// renderLetterHeaderRow renders a row of 26 letters as column headers
func (r *Renderer) renderLetterHeaderRow() string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i, letter := range letters {
		row.WriteString(r.styles.LetterHeader.Render(string(letter)))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderLetterAccuracyRow renders a row of 26 filled circles coloured by accuracy
func (r *Renderer) renderLetterAccuracyRow(historical *stats.HistoricalStats) string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i, letter := range letters {
		lowerLetter := string(letter + 32)
		letterStats := historical.LetterAccuracy[lowerLetter]

		var accuracy float64
		if letterStats.Presented > 0 {
			accuracy = (float64(letterStats.Correct) / float64(letterStats.Presented)) * 100
		}

		colour := GetAccuracyColour(accuracy)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

		row.WriteString(style.Render("●"))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderLetterFrequencyRow renders a row of 26 filled circles coloured by frequency
func (r *Renderer) renderLetterFrequencyRow(historical *stats.HistoricalStats) string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var maxPresented int
	for _, letter := range letters {
		lowerLetter := string(letter + 32)
		if count := historical.LetterAccuracy[lowerLetter].Presented; count > maxPresented {
			maxPresented = count
		}
	}

	for i, letter := range letters {
		lowerLetter := string(letter + 32)
		letterStats := historical.LetterAccuracy[lowerLetter]

		var frequency float64
		if maxPresented > 0 {
			frequency = float64(letterStats.Presented) / float64(maxPresented)
		}

		colour := GetFrequencyColour(frequency)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

		row.WriteString(style.Render("●"))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderLetterSeekTimeRow renders a row of 26 filled circles coloured by seek time
func (r *Renderer) renderLetterSeekTimeRow(historical *stats.HistoricalStats) string {
	var row strings.Builder
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var maxSeekTime float64
	for _, letter := range letters {
		lowerLetter := string(letter + 32)
		if seekStats, exists := historical.LetterSeekTime[lowerLetter]; exists {
			avgTime := seekStats.AverageMs()
			if avgTime > maxSeekTime {
				maxSeekTime = avgTime
			}
		}
	}

	for i, letter := range letters {
		lowerLetter := string(letter + 32)
		seekStats := historical.LetterSeekTime[lowerLetter]
		avgTime := seekStats.AverageMs()

		colour := GetSeekTimeColour(avgTime, maxSeekTime)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

		row.WriteString(style.Render("●"))
		if i < len(letters)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderFingerRow renders finger accuracy
func (r *Renderer) renderFingerRow(historical *stats.HistoricalStats, labelWidth int) string {
	var row strings.Builder
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourLabel)).
		Width(labelWidth).
		Align(lipgloss.Right)

	row.WriteString(labelStyle.Render("Finger accuracy:"))
	row.WriteString(" ")

	fingers := []int{0, 1, 2, 3, 6, 7, 8, 9}
	fingerLabels := []string{"LP", "LR", "LM", "LI", "RI", "RM", "RR", "RP"}
	labelStyleSmall := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourHelp))

	for i, finger := range fingers {
		stat := historical.FingerStats[finger]
		accuracy := stat.Accuracy()
		colour := GetAccuracyColour(accuracy)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		row.WriteString(labelStyleSmall.Render(fingerLabels[i]))
		row.WriteString(style.Render("●"))
		if i < len(fingers)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderRowAccuracyRow renders keyboard row accuracy
func (r *Renderer) renderRowAccuracyRow(historical *stats.HistoricalStats, labelWidth int) string {
	var row strings.Builder
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourLabel)).
		Width(labelWidth).
		Align(lipgloss.Right)

	row.WriteString(labelStyle.Render("Row accuracy:"))
	row.WriteString(" ")

	rows := []int{0, 1, 2}
	rowLabels := []string{"Top", "Home", "Bot"}
	labelStyleSmall := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourHelp))

	for i, rowIdx := range rows {
		stat := historical.RowStats[rowIdx]
		accuracy := stat.Accuracy()
		colour := GetAccuracyColour(accuracy)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		row.WriteString(labelStyleSmall.Render(rowLabels[i]))
		row.WriteString(style.Render("●"))
		if i < len(rows)-1 {
			row.WriteString(" ")
		}
	}

	return row.String()
}

// renderHandStats renders hand balance and alternation stats
func (r *Renderer) renderHandStats(historical *stats.HistoricalStats, labelWidth int) string {
	var row strings.Builder
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourLabel)).
		Width(labelWidth).
		Align(lipgloss.Right)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourValue))

	row.WriteString(labelStyle.Render("Hand balance:"))
	row.WriteString(" ")

	leftStat := historical.HandStats[0]
	rightStat := historical.HandStats[1]
	total := leftStat.Correct + rightStat.Correct
	if total > 0 {
		leftPct := float64(leftStat.Correct) / float64(total) * 100
		rightPct := float64(rightStat.Correct) / float64(total) * 100
		row.WriteString(valueStyle.Render(fmt.Sprintf("L:%.0f%% R:%.0f%%", leftPct, rightPct)))
	} else {
		row.WriteString(valueStyle.Render("N/A"))
	}

	totalTransitions := historical.HandAlternations + historical.SameHandRuns
	if totalTransitions > 0 {
		altRate := float64(historical.HandAlternations) / float64(totalTransitions) * 100
		colour := GetAccuracyColour(altRate)
		altStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))
		row.WriteString(valueStyle.Render("  Alt:"))
		row.WriteString(altStyle.Render(fmt.Sprintf("%.0f%%", altRate)))
	}

	return row.String()
}

// renderRhythmStats renders typing rhythm consistency
func (r *Renderer) renderRhythmStats(session *stats.Stats, historical *stats.HistoricalStats, labelWidth int) string {
	var row strings.Builder
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourLabel)).
		Width(labelWidth).
		Align(lipgloss.Right)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourValue))

	row.WriteString(labelStyle.Render("Rhythm:"))
	row.WriteString(" ")

	sessionStdDev := session.CalculateRhythmStdDev()
	histStdDev := historical.RhythmStats.StdDev()

	maxStdDev := 200.0
	consistency := 100.0 - (sessionStdDev/maxStdDev)*100
	if consistency < 0 {
		consistency = 0
	}
	colour := GetAccuracyColour(consistency)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

	row.WriteString(valueStyle.Render("StdDev: "))
	row.WriteString(style.Render(fmt.Sprintf("%.0fms", sessionStdDev)))
	row.WriteString(valueStyle.Render(fmt.Sprintf(" (avg: %.0fms)", histStdDev)))

	return row.String()
}

// renderSFBStats renders same-finger bigram statistics
func (r *Renderer) renderSFBStats(session *stats.Stats, historical *stats.HistoricalStats, labelWidth int) string {
	var row strings.Builder
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourLabel)).
		Width(labelWidth).
		Align(lipgloss.Right)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColourValue))

	row.WriteString(labelStyle.Render("Same-finger:"))
	row.WriteString(" ")

	sfbCount := session.SFBCount
	var sfbAvg float64
	if sfbCount > 0 {
		sfbAvg = float64(session.SFBTotalTime) / float64(sfbCount)
	}

	histAvg := historical.SFBStats.AverageMs()

	maxSFBTime := 400.0
	performance := 100.0 - (sfbAvg/maxSFBTime)*100
	if performance < 0 {
		performance = 0
	}
	colour := GetAccuracyColour(performance)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(colour))

	row.WriteString(valueStyle.Render(fmt.Sprintf("%d SFBs", sfbCount)))
	if sfbCount > 0 {
		row.WriteString(style.Render(fmt.Sprintf(" @%.0fms", sfbAvg)))
	}
	if histAvg > 0 {
		row.WriteString(valueStyle.Render(fmt.Sprintf(" (avg: %.0fms)", histAvg)))
	}

	return row.String()
}

// renderTopErrors renders top error substitution patterns
func (r *Renderer) renderTopErrors(historical *stats.HistoricalStats, labelWidth int) string {
	var row strings.Builder
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColourLabel)).
		Width(labelWidth).
		Align(lipgloss.Right)

	row.WriteString(labelStyle.Render("Common errors:"))
	row.WriteString(" ")

	type errorPair struct {
		expected string
		typed    string
		count    int
	}
	var errors []errorPair
	for expected, typedMap := range historical.ErrorSubstitution {
		for typed, count := range typedMap {
			errors = append(errors, errorPair{expected, typed, count})
		}
	}

	// Sort by count
	for i := 0; i < len(errors); i++ {
		for j := i + 1; j < len(errors); j++ {
			if errors[j].count > errors[i].count {
				errors[i], errors[j] = errors[j], errors[i]
			}
		}
	}

	shown := 0
	for _, e := range errors {
		if shown >= 5 {
			break
		}
		if shown > 0 {
			row.WriteString(" ")
		}
		row.WriteString(r.styles.ErrorStyle.Render(e.expected + "→" + e.typed))
		row.WriteString(r.styles.CountStyle.Render(fmt.Sprintf("(%d)", e.count)))
		shown++
	}

	if shown == 0 {
		row.WriteString(r.styles.CountStyle.Render("none"))
	}

	return row.String()
}

// RenderOptionsScreen renders the options/settings screen
func (r *Renderer) RenderOptionsScreen(s *settings.Settings, cursor int) string {
	title := r.styles.Title.Render("Options")

	// Options for advance key
	options := []struct {
		key         settings.AdvanceKey
		label       string
		description string
	}{
		{settings.AdvanceKeySpace, "Space", "Press Space to advance to the next word (default)"},
		{settings.AdvanceKeyEnter, "Enter", "Press Enter to advance to the next word"},
		{settings.AdvanceKeyEither, "Either", "Press Space or Enter to advance to the next word"},
	}

	var optionLines []string
	optionLines = append(optionLines, "")
	optionLines = append(optionLines, r.styles.SessionLabel.Render("Advance to next word with:"))
	optionLines = append(optionLines, "")

	for i, opt := range options {
		// Build the option line
		var line strings.Builder

		// Number prefix
		numStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		line.WriteString(numStyle.Render(fmt.Sprintf(" %d. ", i+1)))

		// Selection indicator and label
		isSelected := s.AdvanceKey == opt.key
		isCursor := cursor == i

		var labelStyle lipgloss.Style
		if isCursor {
			// Cursor position - highlighted
			labelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("39")).
				Bold(true).
				Padding(0, 1)
		} else if isSelected {
			// Currently selected option
			labelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46")).
				Bold(true)
		} else {
			// Normal option
			labelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))
		}

		// Checkmark for selected option
		if isSelected {
			checkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
			line.WriteString(checkStyle.Render("✓ "))
		} else {
			line.WriteString("  ")
		}

		line.WriteString(labelStyle.Render(opt.label))

		// Description
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
		line.WriteString("  ")
		line.WriteString(descStyle.Render(opt.description))

		optionLines = append(optionLines, line.String())
	}

	// Main content (title + options)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		strings.Join(optionLines, "\n"),
	)

	// Fixed header at top
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)
	header := lipgloss.PlaceHorizontal(r.width, lipgloss.Center,
		headerStyle.Render("🐒 BABOON - Typing Practice"))

	// Fixed footer at bottom
	footer := lipgloss.PlaceHorizontal(r.width, lipgloss.Center,
		r.styles.Help.Render("↑/↓ to navigate | Enter/Space to select | 1-3 quick select | ESC to go back"))

	// Calculate heights
	headerHeight := 1
	footerHeight := 1
	contentHeight := strings.Count(mainContent, "\n") + 1
	availableHeight := r.height - headerHeight - footerHeight - 2 // -2 for spacing

	// Calculate top padding to center main content in available space
	topPadding := (availableHeight - contentHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Build full screen layout
	var fullContent strings.Builder

	// Header at top (line 0)
	fullContent.WriteString(header)
	fullContent.WriteString("\n")

	// Top padding to center content
	for i := 0; i < topPadding; i++ {
		fullContent.WriteString("\n")
	}

	// Main content (centered horizontally)
	centeredMain := lipgloss.PlaceHorizontal(r.width, lipgloss.Center, mainContent)
	fullContent.WriteString(centeredMain)

	// Bottom padding to push footer to the bottom
	currentHeight := headerHeight + 1 + topPadding + contentHeight
	for i := currentHeight; i < r.height-footerHeight; i++ {
		fullContent.WriteString("\n")
	}

	// Footer at bottom (last line)
	fullContent.WriteString(footer)

	return fullContent.String()
}
