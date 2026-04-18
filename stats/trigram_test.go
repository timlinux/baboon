package stats

import "testing"

func TestTrigramSeekStats_AverageMs(t *testing.T) {
	tests := []struct {
		name  string
		stats TrigramSeekStats
		want  float64
	}{
		{"empty", TrigramSeekStats{}, 0},
		{"single", TrigramSeekStats{TotalTimeMs: 100, Count: 1}, 100},
		{"multiple", TrigramSeekStats{TotalTimeMs: 300, Count: 3}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stats.AverageMs(); got != tt.want {
				t.Errorf("AverageMs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStats_RecordTrigramSeekTime(t *testing.T) {
	s := &Stats{
		TrigramSeekTime: make(map[string]TrigramSeekStats),
	}

	s.RecordTrigramSeekTime("the", 150)
	s.RecordTrigramSeekTime("the", 200)
	s.RecordTrigramSeekTime("ing", 100)

	if got := s.TrigramSeekTime["the"].Count; got != 2 {
		t.Errorf("the count = %v, want 2", got)
	}
	if got := s.TrigramSeekTime["the"].TotalTimeMs; got != 350 {
		t.Errorf("the total = %v, want 350", got)
	}
	if got := s.TrigramSeekTime["ing"].Count; got != 1 {
		t.Errorf("ing count = %v, want 1", got)
	}
}

func TestHistoricalStats_NgramAverages(t *testing.T) {
	h := &HistoricalStats{
		NgramTotalWPM:      200,
		NgramTotalAccuracy: 190,
		NgramTotalTime:     300,
		NgramTotalSessions: 2,
	}

	if got := h.NgramAverageWPM(); got != 100 {
		t.Errorf("NgramAverageWPM() = %v, want 100", got)
	}
	if got := h.NgramAverageAccuracy(); got != 95 {
		t.Errorf("NgramAverageAccuracy() = %v, want 95", got)
	}
	if got := h.NgramAverageTime(); got != 150 {
		t.Errorf("NgramAverageTime() = %v, want 150", got)
	}
}

func TestHistoricalStats_UpdateNgramHistorical(t *testing.T) {
	h := &HistoricalStats{
		TrigramSeekTime: make(map[string]TrigramSeekStats),
	}
	session := &Stats{
		WPM:      80,
		Accuracy: 95,
		Duration: 60 * 1000000000, // 60 seconds in nanoseconds
		TrigramSeekTime: map[string]TrigramSeekStats{
			"the": {TotalTimeMs: 200, Count: 2},
		},
	}

	h.UpdateNgramHistorical(session)

	if h.NgramTotalSessions != 1 {
		t.Errorf("NgramTotalSessions = %v, want 1", h.NgramTotalSessions)
	}
	if h.NgramBestWPM != 80 {
		t.Errorf("NgramBestWPM = %v, want 80", h.NgramBestWPM)
	}
	if h.TrigramSeekTime["the"].Count != 2 {
		t.Errorf("TrigramSeekTime[the].Count = %v, want 2", h.TrigramSeekTime["the"].Count)
	}
}
