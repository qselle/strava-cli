package format

import "testing"

func TestDuration(t *testing.T) {
	tests := []struct {
		seconds float64
		want    string
	}{
		{0, "0m"},
		{30, "0m"},
		{60, "1m"},
		{90, "1m"},
		{3600, "1h00m"},
		{3661, "1h01m"},
		{7200, "2h00m"},
		{11446.0, "3h10m"},
		{-60, "-1m"},
	}

	for _, tt := range tests {
		got := Duration(tt.seconds)
		if got != tt.want {
			t.Errorf("Duration(%v) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestDate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2026-03-18T19:00:00Z", "2026-03-18"},
		{"2026-03-18T19:00:00+01:00", "2026-03-18"},
		{"2026-03-18T19:00:00", "2026-03-18"},
		{"2026-03-18", "2026-03-18"},
		{"short", "unknown"},
		{"", "unknown"},
	}

	for _, tt := range tests {
		got := Date(tt.input)
		if got != tt.want {
			t.Errorf("Date(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMotivation(t *testing.T) {
	tests := []struct {
		active, total, streak int
		wantVerdict           string
	}{
		{0, 7, 0, "COUCH POTATO MODE"},
		{1, 7, 0, "BARELY ALIVE"},
		{2, 7, 0, "WARMING UP"},
		{3, 7, 0, "GETTING THERE"},
		{5, 7, 0, "CRUSHING IT"},
		{7, 7, 7, "ABSOLUTE LEGEND"},
		{6, 7, 0, "ON FIRE"},
		{0, 0, 0, "NO DATA"},
		{0, -1, 0, "NO DATA"},
	}

	for _, tt := range tests {
		verdict, _ := Motivation(tt.active, tt.total, tt.streak)
		if verdict != tt.wantVerdict {
			t.Errorf("Motivation(%d, %d, %d) verdict = %q, want %q",
				tt.active, tt.total, tt.streak, verdict, tt.wantVerdict)
		}
	}
}
