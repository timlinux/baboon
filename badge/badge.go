// Package badge generates shareable SVG badges for leaderboard entries.
package badge

import (
	"fmt"
	"html"
	"strings"
)

// Colors matching the Baboon theme
const (
	ColorBackground = "#1a2833"
	ColorOrange     = "#D4922A"
	ColorCyan       = "#4A90A4"
	ColorGreen      = "#4CAF50"
	ColorGold       = "#FFD700"
	ColorWhite      = "#FFFFFF"
	ColorGray       = "#808080"
	ColorDarkGray   = "#2a3843"
)

// Entry contains the data needed to generate a badge.
type Entry struct {
	DisplayName string
	WPM         float64
	Accuracy    float64
	Rank        int
	MonthYear   string
	SiteURL     string
}

// GenerateSVG creates an SVG badge for a leaderboard entry.
func GenerateSVG(entry Entry) string {
	name := html.EscapeString(entry.DisplayName)
	wpm := fmt.Sprintf("%.1f", entry.WPM)
	accuracy := fmt.Sprintf("%.1f%%", entry.Accuracy)

	// Rank star display
	rankStar := ""
	if entry.Rank >= 1 && entry.Rank <= 10 {
		starColor := ColorGold
		if entry.Rank == 1 {
			starColor = ColorGold
		} else if entry.Rank <= 3 {
			starColor = ColorOrange
		} else {
			starColor = ColorGray
		}
		rankStar = fmt.Sprintf(`<text x="350" y="55" fill="%s" font-size="24" font-family="monospace" text-anchor="middle">★ #%d</text>`, starColor, entry.Rank)
	}

	// Format month display
	monthDisplay := entry.MonthYear
	if len(entry.MonthYear) == 7 { // "2026-04"
		parts := strings.Split(entry.MonthYear, "-")
		if len(parts) == 2 {
			months := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
			monthNum := 0
			fmt.Sscanf(parts[1], "%d", &monthNum)
			if monthNum >= 1 && monthNum <= 12 {
				monthDisplay = fmt.Sprintf("%s %s", months[monthNum], parts[0])
			}
		}
	}

	siteURL := entry.SiteURL
	if siteURL == "" {
		siteURL = "baboon.kartoza.com"
	}

	// Simplified baboon mascot (stylized head + keyboard)
	baboonMascot := `
		<!-- Baboon head -->
		<ellipse cx="60" cy="55" rx="35" ry="40" fill="#8B4513"/>
		<ellipse cx="60" cy="60" rx="20" ry="22" fill="#DEB887"/>
		<!-- Eyes -->
		<circle cx="50" cy="45" r="6" fill="white"/>
		<circle cx="70" cy="45" r="6" fill="white"/>
		<circle cx="51" cy="46" r="3" fill="black"/>
		<circle cx="71" cy="46" r="3" fill="black"/>
		<!-- Nose -->
		<ellipse cx="60" cy="60" rx="8" ry="5" fill="#CD853F"/>
		<circle cx="56" cy="60" r="2" fill="#8B4513"/>
		<circle cx="64" cy="60" r="2" fill="#8B4513"/>
		<!-- Mouth -->
		<path d="M 52 70 Q 60 75 68 70" stroke="#8B4513" stroke-width="2" fill="none"/>
		<!-- Ears -->
		<ellipse cx="28" cy="40" rx="8" ry="12" fill="#8B4513"/>
		<ellipse cx="92" cy="40" rx="8" ry="12" fill="#8B4513"/>
		<!-- Keyboard below -->
		<rect x="25" y="100" width="70" height="25" rx="3" fill="#333" stroke="#555" stroke-width="1"/>
		<rect x="30" y="105" width="8" height="6" rx="1" fill="#555"/>
		<rect x="40" y="105" width="8" height="6" rx="1" fill="#555"/>
		<rect x="50" y="105" width="8" height="6" rx="1" fill="#555"/>
		<rect x="60" y="105" width="8" height="6" rx="1" fill="#555"/>
		<rect x="70" y="105" width="8" height="6" rx="1" fill="#555"/>
		<rect x="82" y="105" width="8" height="6" rx="1" fill="#555"/>
		<rect x="30" y="113" width="8" height="6" rx="1" fill="#555"/>
		<rect x="40" y="113" width="8" height="6" rx="1" fill="#555"/>
		<rect x="50" y="113" width="20" height="6" rx="1" fill="#666"/>
		<rect x="72" y="113" width="8" height="6" rx="1" fill="#555"/>
		<rect x="82" y="113" width="8" height="6" rx="1" fill="#555"/>
	`

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 200" width="400" height="200">
	<!-- Background -->
	<rect width="400" height="200" fill="%s" rx="10"/>
	<rect x="5" y="5" width="390" height="190" fill="none" stroke="%s" stroke-width="2" rx="8"/>

	<!-- Mascot -->
	<g transform="translate(0, 10)">
		%s
	</g>

	<!-- Player name -->
	<text x="250" y="45" fill="%s" font-size="28" font-family="'Fira Code', monospace" font-weight="bold" text-anchor="middle">%s</text>

	<!-- Rank -->
	%s

	<!-- WPM -->
	<text x="180" y="90" fill="%s" font-size="36" font-family="'Fira Code', monospace" font-weight="bold">%s</text>
	<text x="290" y="90" fill="%s" font-size="20" font-family="'Fira Code', monospace">WPM</text>

	<!-- Accuracy -->
	<text x="180" y="120" fill="%s" font-size="18" font-family="'Fira Code', monospace">%s accuracy</text>

	<!-- Month -->
	<text x="180" y="145" fill="%s" font-size="14" font-family="'Fira Code', monospace">%s</text>

	<!-- Site URL -->
	<text x="200" y="170" fill="%s" font-size="12" font-family="'Fira Code', monospace" text-anchor="middle">%s</text>

	<!-- Kartoza branding -->
	<text x="200" y="188" fill="%s" font-size="10" font-family="'Fira Code', monospace" text-anchor="middle">Made with ♥ by Kartoza</text>
</svg>`,
		ColorBackground,
		ColorOrange,
		baboonMascot,
		ColorWhite,
		name,
		rankStar,
		ColorGreen,
		wpm,
		ColorCyan,
		ColorWhite,
		accuracy,
		ColorGray,
		monthDisplay,
		ColorCyan,
		siteURL,
		ColorGray,
	)

	return svg
}

// GenerateShareURL creates a shareable URL for a leaderboard entry.
func GenerateShareURL(baseURL, entryID string) string {
	return fmt.Sprintf("%s/leaderboard?highlight=%s", baseURL, entryID)
}
