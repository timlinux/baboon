package ngrams

import (
	"testing"

	"github.com/timlinux/baboon/stats"
)

func TestGetTrainingSequence_ReturnsCorrectCharCount(t *testing.T) {
	bigramData := map[string]stats.BigramSeekStats{}
	trigramData := map[string]stats.TrigramSeekStats{}

	result := GetTrainingSequence(150, bigramData, trigramData, nil)

	// Count total characters (including spaces between n-grams)
	totalChars := 0
	for i, ngram := range result {
		totalChars += len(ngram)
		if i < len(result)-1 {
			totalChars++ // space between n-grams
		}
	}

	// Should be within 10 chars of target (allow for n-gram boundaries)
	if totalChars < 140 || totalChars > 160 {
		t.Errorf("Total chars = %d, want 140-160", totalChars)
	}
}

func TestGetTrainingSequence_HandlesEmptyHistory(t *testing.T) {
	result := GetTrainingSequence(150, nil, nil, nil)

	if len(result) == 0 {
		t.Error("Expected non-empty result with empty history")
	}

	// Verify all results are valid n-grams (2-3 lowercase letters)
	for _, ngram := range result {
		if len(ngram) < 2 || len(ngram) > 3 {
			t.Errorf("Invalid n-gram length: %q", ngram)
		}
		for _, c := range ngram {
			if c < 'a' || c > 'z' {
				t.Errorf("Invalid character in n-gram: %q", ngram)
			}
		}
	}
}

func TestGetTrainingSequence_PrioritisesSlowNgrams(t *testing.T) {
	// Create data where "zz" is very slow
	bigramData := map[string]stats.BigramSeekStats{
		"zz": {TotalTimeMs: 1000, Count: 5}, // 200ms average - slow
		"th": {TotalTimeMs: 250, Count: 5},  // 50ms average - fast
		"er": {TotalTimeMs: 250, Count: 5},  // 50ms average - fast
	}

	// Generate multiple sequences and count "zz" occurrences
	zzCount := 0
	thCount := 0
	for i := 0; i < 100; i++ {
		result := GetTrainingSequence(150, bigramData, nil, func(n int) int { return i % n })
		for _, ngram := range result {
			if ngram == "zz" {
				zzCount++
			}
			if ngram == "th" {
				thCount++
			}
		}
	}

	// "zz" should appear more often than "th" due to higher priority
	if zzCount <= thCount {
		t.Errorf("Slow n-gram 'zz' count (%d) should exceed fast 'th' count (%d)", zzCount, thCount)
	}
}

func TestGetTrainingSequence_MixesBigramsAndTrigrams(t *testing.T) {
	result := GetTrainingSequence(150, nil, nil, nil)

	bigramCount := 0
	trigramCount := 0
	for _, ngram := range result {
		if len(ngram) == 2 {
			bigramCount++
		} else if len(ngram) == 3 {
			trigramCount++
		}
	}

	// Should have mix of both
	if bigramCount == 0 {
		t.Error("Expected some bigrams in result")
	}
	if trigramCount == 0 {
		t.Error("Expected some trigrams in result")
	}

	// Trigrams should be minority (~30%)
	total := bigramCount + trigramCount
	trigramRatio := float64(trigramCount) / float64(total)
	if trigramRatio > 0.5 {
		t.Errorf("Trigram ratio %.2f is too high, expected ~0.3", trigramRatio)
	}
}
