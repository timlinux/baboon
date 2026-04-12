// Package filter provides profanity filtering for display names.
package filter

import (
	"regexp"
	"strings"
	"unicode"
)

// ValidationResult contains the result of name validation.
type ValidationResult struct {
	Valid       bool  `json:"valid"`
	InvalidChars []int `json:"invalid_chars,omitempty"` // Indices of invalid characters
	Reason      string `json:"reason,omitempty"`
}

// Leetspeak substitutions for normalization
var leetMap = map[rune]rune{
	'4': 'a',
	'@': 'a',
	'8': 'b',
	'3': 'e',
	'6': 'g',
	'1': 'i',
	'!': 'i',
	'|': 'i',
	'0': 'o',
	'5': 's',
	'$': 's',
	'7': 't',
	'+': 't',
	'2': 'z',
}

// Blocklist of inappropriate words (normalized lowercase)
// This is a minimal list - in production you'd want a more comprehensive one
var blocklist = map[string]bool{
	"fuck":     true,
	"shit":     true,
	"ass":      true,
	"asshole":  true,
	"bitch":    true,
	"bastard":  true,
	"damn":     true,
	"dick":     true,
	"penis":    true,
	"vagina":   true,
	"pussy":    true,
	"cock":     true,
	"cunt":     true,
	"slut":     true,
	"whore":    true,
	"fag":      true,
	"faggot":   true,
	"nigger":   true,
	"nigga":    true,
	"retard":   true,
	"nazi":     true,
	"hitler":   true,
	"porn":     true,
	"sex":      true,
	"rape":     true,
	"kill":     true,
	"murder":   true,
	"piss":     true,
	"crap":     true,
	"tits":     true,
	"boobs":    true,
	"anus":     true,
	"anal":     true,
	"cum":      true,
	"jizz":     true,
	"wank":     true,
	"dildo":    true,
	"erect":    true,
	"horny":    true,
	"milf":     true,
	"nude":     true,
	"naked":    true,
	"xxx":      true,
}

// Allowed characters for display names
var allowedChars = regexp.MustCompile(`^[A-Za-z0-9 _-]+$`)

// normalizeLeetspeak converts leetspeak characters to their letter equivalents.
func normalizeLeetspeak(s string) string {
	var result strings.Builder
	for _, r := range s {
		if mapped, ok := leetMap[r]; ok {
			result.WriteRune(mapped)
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

// containsBlockedWord checks if a normalized string contains any blocked words.
func containsBlockedWord(normalized string) (bool, string) {
	// Remove spaces and special chars for checking
	stripped := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, normalized)

	// Check for exact matches and substrings
	for word := range blocklist {
		if strings.Contains(stripped, word) {
			return true, word
		}
	}

	return false, ""
}

// MaxMessageLength is the maximum length for leaderboard messages.
// Kept at 10 for classic arcade-style high score entries.
const MaxMessageLength = 10

// Validate checks if a display message is valid for the leaderboard.
// Returns a ValidationResult with details about any issues.
func Validate(name string) ValidationResult {
	// Check length
	if len(name) == 0 {
		return ValidationResult{
			Valid:  false,
			Reason: "Message cannot be empty",
		}
	}

	if len(name) > MaxMessageLength {
		return ValidationResult{
			Valid:  false,
			Reason: "Message must be 10 characters or less",
		}
	}

	// Check for valid characters
	if !allowedChars.MatchString(name) {
		invalidChars := make([]int, 0)
		for i, r := range name {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '_' && r != '-' {
				invalidChars = append(invalidChars, i)
			}
		}
		return ValidationResult{
			Valid:        false,
			InvalidChars: invalidChars,
			Reason:       "Name contains invalid characters. Only letters, numbers, spaces, underscores, and hyphens are allowed.",
		}
	}

	// Check for whitespace-only
	if strings.TrimSpace(name) == "" {
		return ValidationResult{
			Valid:  false,
			Reason: "Name cannot be only whitespace",
		}
	}

	// Check for profanity (with leetspeak normalization)
	normalized := normalizeLeetspeak(name)
	if blocked, word := containsBlockedWord(normalized); blocked {
		// Find indices of the problematic characters
		normalizedLower := strings.ToLower(name)
		startIdx := strings.Index(normalizeLeetspeak(normalizedLower), word)
		invalidChars := make([]int, 0)
		if startIdx >= 0 {
			for i := startIdx; i < startIdx+len(word) && i < len(name); i++ {
				invalidChars = append(invalidChars, i)
			}
		}
		return ValidationResult{
			Valid:        false,
			InvalidChars: invalidChars,
			Reason:       "Name contains inappropriate content",
		}
	}

	return ValidationResult{
		Valid: true,
	}
}

// IsValidChar checks if a single character is valid for a display name.
func IsValidChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_' || r == '-'
}

// SanitizeName removes invalid characters from a message.
func SanitizeName(name string) string {
	var result strings.Builder
	for _, r := range name {
		if IsValidChar(r) {
			result.WriteRune(r)
		}
	}
	sanitized := result.String()
	if len(sanitized) > MaxMessageLength {
		sanitized = sanitized[:MaxMessageLength]
	}
	return sanitized
}
