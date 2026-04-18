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
