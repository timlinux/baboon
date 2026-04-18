// Package ngrams provides n-gram generation for typing practice.
package ngrams

import (
	"sort"

	"github.com/timlinux/baboon/stats"
)

// CommonBigrams contains frequently occurring English bigrams
var CommonBigrams = []string{
	"th", "he", "in", "er", "an", "re", "on", "at", "en", "nd",
	"ti", "es", "or", "te", "of", "ed", "is", "it", "al", "ar",
	"st", "to", "nt", "ng", "se", "ha", "as", "ou", "io", "le",
	"ve", "co", "me", "de", "hi", "ri", "ro", "ic", "ne", "ea",
	"ra", "ce", "li", "ch", "ll", "be", "ma", "si", "om", "ur",
}

// CommonTrigrams contains frequently occurring English trigrams
var CommonTrigrams = []string{
	"the", "and", "ing", "ion", "tio", "ent", "ati", "for", "her", "ter",
	"hat", "tha", "ere", "ate", "his", "con", "res", "ver", "all", "ons",
	"nce", "men", "ith", "ted", "ers", "pro", "thi", "wit", "are", "ess",
	"not", "ive", "was", "ect", "rea", "com", "eve", "per", "int", "est",
	"sta", "cti", "ica", "ist", "ear", "ain", "one", "our", "iti", "rat",
}

// scoredNgram holds an n-gram with its selection score
type scoredNgram struct {
	ngram string
	score float64
}

// GetTrainingSequence generates a sequence of n-grams for training.
// It prioritises n-grams where the user has slower seek times.
// Returns a slice of n-grams that together total approximately targetChars characters.
func GetTrainingSequence(
	targetChars int,
	bigramData map[string]stats.BigramSeekStats,
	trigramData map[string]stats.TrigramSeekStats,
	randFunc func(n int) int,
) []string {
	if randFunc == nil {
		randFunc = defaultRand
	}

	// Build scored list of candidates
	candidates := buildCandidates(bigramData, trigramData)

	// Generate sequence
	result := []string{}
	currentChars := 0

	for currentChars < targetChars {
		// Decide bigram vs trigram (~70/30 split)
		useTrigram := randFunc(10) < 3

		// Select from appropriate candidates
		var selected string
		if useTrigram {
			selected = selectWeighted(filterByLength(candidates, 3), randFunc)
		} else {
			selected = selectWeighted(filterByLength(candidates, 2), randFunc)
		}

		// Fallback if no candidates of desired length
		if selected == "" {
			if useTrigram && len(CommonTrigrams) > 0 {
				selected = CommonTrigrams[randFunc(len(CommonTrigrams))]
			} else if len(CommonBigrams) > 0 {
				selected = CommonBigrams[randFunc(len(CommonBigrams))]
			}
		}

		if selected == "" {
			break
		}

		// Check if adding this would exceed target (with space)
		addedChars := len(selected)
		if len(result) > 0 {
			addedChars++ // space
		}
		if currentChars+addedChars > targetChars+5 {
			break
		}

		result = append(result, selected)
		currentChars += addedChars
	}

	return result
}

// buildCandidates creates a scored list of n-gram candidates
func buildCandidates(
	bigramData map[string]stats.BigramSeekStats,
	trigramData map[string]stats.TrigramSeekStats,
) []scoredNgram {
	candidates := []scoredNgram{}

	// Add bigrams from history
	for ngram, data := range bigramData {
		if data.Count >= 3 {
			candidates = append(candidates, scoredNgram{
				ngram: ngram,
				score: data.AverageMs(),
			})
		}
	}

	// Add trigrams from history
	for ngram, data := range trigramData {
		if data.Count >= 3 {
			candidates = append(candidates, scoredNgram{
				ngram: ngram,
				score: data.AverageMs(),
			})
		}
	}

	// Add common n-grams with medium priority if not already present
	existingNgrams := make(map[string]bool)
	for _, c := range candidates {
		existingNgrams[c.ngram] = true
	}

	mediumScore := 150.0 // Default score for unknown n-grams
	for _, ngram := range CommonBigrams {
		if !existingNgrams[ngram] {
			candidates = append(candidates, scoredNgram{ngram: ngram, score: mediumScore})
		}
	}
	for _, ngram := range CommonTrigrams {
		if !existingNgrams[ngram] {
			candidates = append(candidates, scoredNgram{ngram: ngram, score: mediumScore})
		}
	}

	// Sort by score descending (slower = higher priority)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	return candidates
}

// filterByLength returns only n-grams of the specified length
func filterByLength(candidates []scoredNgram, length int) []scoredNgram {
	filtered := []scoredNgram{}
	for _, c := range candidates {
		if len(c.ngram) == length {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// selectWeighted selects an n-gram with probability proportional to score
func selectWeighted(candidates []scoredNgram, randFunc func(n int) int) string {
	if len(candidates) == 0 {
		return ""
	}

	// Calculate total score
	totalScore := 0.0
	for _, c := range candidates {
		totalScore += c.score
	}

	if totalScore == 0 {
		return candidates[randFunc(len(candidates))].ngram
	}

	// Select weighted random
	target := float64(randFunc(int(totalScore*100))) / 100.0
	cumulative := 0.0
	for _, c := range candidates {
		cumulative += c.score
		if cumulative >= target {
			return c.ngram
		}
	}

	return candidates[len(candidates)-1].ngram
}

// defaultRand is a simple deterministic pseudo-random for testing
var randState int = 42

func defaultRand(n int) int {
	if n <= 0 {
		return 0
	}
	randState = (randState*1103515245 + 12345) & 0x7fffffff
	return randState % n
}
