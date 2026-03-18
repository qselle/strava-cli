package format

import (
	"fmt"
	"time"
)

// Duration formats seconds into a human-readable duration string.
func Duration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// Date extracts YYYY-MM-DD from a Strava date string.
// Returns "unknown" if the input is too short or unparseable.
func Date(dateStr string) string {
	if len(dateStr) < 10 {
		return "unknown"
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
			return dateStr[:10]
		}
	}
	return t.Format("2006-01-02")
}

// Motivation returns a verdict and motivational message based on activity ratio.
func Motivation(activeDays, totalDays, currentStreak int) (string, string) {
	if totalDays <= 0 {
		return "NO DATA", "No period specified."
	}

	ratio := float64(activeDays) / float64(totalDays)

	switch {
	case activeDays == 0:
		return "COUCH POTATO MODE",
			"Your couch misses you... oh wait, you never left. Time to move!"
	case ratio < 0.2:
		return "BARELY ALIVE",
			"One activity is better than none, but your shoes are getting dusty."
	case ratio < 0.4:
		return "WARMING UP",
			"You're showing signs of life! Keep building that momentum."
	case ratio < 0.6:
		return "GETTING THERE",
			"Solid effort! You're building a real habit here."
	case ratio < 0.8:
		return "CRUSHING IT",
			"Beast mode activated! Your consistency is paying off."
	case currentStreak >= totalDays:
		return "ABSOLUTE LEGEND",
			"Perfect streak! You haven't missed a single day. Unstoppable!"
	default:
		return "ON FIRE",
			"You're on fire! Keep this up and nothing can stop you."
	}
}
