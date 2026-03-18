package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/qselle/strava-cli/internal/api"
	"github.com/qselle/strava-cli/internal/auth"
)

var streakDays int

var streakCmd = &cobra.Command{
	Use:   "streak",
	Short: "Check if you've been moving your ass",
	Long:  "Analyze your recent activity streak.\nShows consecutive active days, rest days, and motivational feedback.\nPerfect for AI agents to keep you accountable.",
	RunE:  runStreak,
}

func init() {
	streakCmd.Flags().IntVarP(&streakDays, "days", "d", 7, "Number of days to look back")
	rootCmd.AddCommand(streakCmd)
}

type StreakResult struct {
	Period         string   `json:"period"`
	TotalDays      int      `json:"total_days"`
	ActiveDays     int      `json:"active_days"`
	RestDays       int      `json:"rest_days"`
	CurrentStreak  int      `json:"current_streak"`
	TotalDistance   float64  `json:"total_distance_km"`
	TotalTime      int      `json:"total_time_seconds"`
	Activities     int      `json:"activities"`
	ActiveDaysList []string `json:"active_days_list"`
	Verdict        string   `json:"verdict"`
	Motivation     string   `json:"motivation"`
}

func runStreak(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	token, err := auth.GetValidToken(ctx)
	if err != nil {
		return err
	}

	client := api.NewClient(token.AccessToken)

	now := time.Now()
	after := now.AddDate(0, 0, -streakDays)

	activities, err := client.ListActivities(ctx, api.ListActivitiesParams{
		After:   after.Unix(),
		PerPage: 200,
	})
	if err != nil {
		return fmt.Errorf("fetching activities: %w", err)
	}

	activeDays := make(map[string]bool)
	var totalDistance float64
	var totalTime int

	for _, a := range activities {
		date := formatDate(a.StartDateLocal)
		activeDays[date] = true
		totalDistance += a.Distance / 1000
		totalTime += a.MovingTime
	}

	activeDaysList := make([]string, 0, len(activeDays))
	for day := range activeDays {
		activeDaysList = append(activeDaysList, day)
	}

	currentStreak := 0
	for i := 0; i < streakDays; i++ {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		if activeDays[day] {
			currentStreak++
		} else {
			break
		}
	}

	restDays := streakDays - len(activeDays)
	verdict, motivation := getMotivation(len(activeDays), streakDays, currentStreak)

	result := StreakResult{
		Period:         fmt.Sprintf("last %d days", streakDays),
		TotalDays:      streakDays,
		ActiveDays:     len(activeDays),
		RestDays:       restDays,
		CurrentStreak:  currentStreak,
		TotalDistance:   totalDistance,
		TotalTime:      totalTime,
		Activities:     len(activities),
		ActiveDaysList: activeDaysList,
		Verdict:        verdict,
		Motivation:     motivation,
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Printf("Streak report (last %d days)\n\n", streakDays)
	fmt.Printf("  Active days:    %d / %d\n", result.ActiveDays, result.TotalDays)
	fmt.Printf("  Current streak: %d days\n", result.CurrentStreak)
	fmt.Printf("  Activities:     %d\n", result.Activities)
	fmt.Printf("  Total distance: %.1f km\n", result.TotalDistance)
	fmt.Printf("  Total time:     %s\n\n", formatDuration(result.TotalTime))
	fmt.Printf("  %s\n", result.Verdict)
	fmt.Printf("  %s\n", result.Motivation)

	return nil
}

func getMotivation(activeDays, totalDays, currentStreak int) (string, string) {
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
