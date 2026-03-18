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

var (
	activitiesAfter  string
	activitiesBefore string
	activitiesLimit  int
)

var activitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "List your recent Strava activities",
	Long:  "Retrieve and display your recent activities from Strava.\nSupports filtering by date range and limiting results.",
	RunE:  runActivities,
}

func init() {
	activitiesCmd.Flags().StringVar(&activitiesAfter, "after", "", "Show activities after this date (YYYY-MM-DD)")
	activitiesCmd.Flags().StringVar(&activitiesBefore, "before", "", "Show activities before this date (YYYY-MM-DD)")
	activitiesCmd.Flags().IntVarP(&activitiesLimit, "limit", "n", 10, "Maximum number of activities to show")
	rootCmd.AddCommand(activitiesCmd)
}

func runActivities(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cfg := getOAuthConfig()

	token, err := auth.GetValidToken(ctx, cfg)
	if err != nil {
		return err
	}

	client := api.NewClient(token.AccessToken)

	params := api.ListActivitiesParams{
		PerPage: activitiesLimit,
	}

	if activitiesAfter != "" {
		t, err := time.Parse("2006-01-02", activitiesAfter)
		if err != nil {
			return fmt.Errorf("invalid --after date: %w", err)
		}
		params.After = t.Unix()
	}

	if activitiesBefore != "" {
		t, err := time.Parse("2006-01-02", activitiesBefore)
		if err != nil {
			return fmt.Errorf("invalid --before date: %w", err)
		}
		params.Before = t.Unix()
	}

	activities, err := client.ListActivities(ctx, params)
	if err != nil {
		return fmt.Errorf("fetching activities: %w", err)
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(activities)
	}

	if len(activities) == 0 {
		fmt.Println("No activities found.")
		return nil
	}

	for _, a := range activities {
		date := formatDate(a.StartDateLocal)
		distance := a.Distance / 1000
		duration := formatDuration(a.MovingTime)
		fmt.Printf("  %s  %-12s  %6.1f km  %s  %s\n", date, a.SportType, distance, duration, a.Name)
	}

	return nil
}

func formatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
			return dateStr[:10]
		}
	}
	return t.Format("2006-01-02")
}
