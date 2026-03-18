package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/qselle/strava-cli/internal/api"
	"github.com/qselle/strava-cli/internal/auth"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show your all-time and year-to-date Strava stats",
	Long:  "Display your lifetime, year-to-date, and recent activity totals.\nIncludes distance, time, elevation, and activity counts for runs, rides, and swims.",
	RunE:  runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cfg := getOAuthConfig()

	token, err := auth.GetValidToken(ctx, cfg)
	if err != nil {
		return err
	}

	client := api.NewClient(token.AccessToken)

	athlete, err := client.GetAthlete(ctx)
	if err != nil {
		return fmt.Errorf("fetching athlete: %w", err)
	}

	stats, err := client.GetStats(ctx, athlete.ID)
	if err != nil {
		return fmt.Errorf("fetching stats: %w", err)
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(stats)
	}

	fmt.Printf("Stats for %s %s\n\n", athlete.FirstName, athlete.LastName)

	printTotals := func(label string, run, ride, swim api.ActivityTotal) {
		fmt.Printf("  %s:\n", label)
		if run.Count > 0 {
			fmt.Printf("    Run:  %3d activities  %8.1f km  %s  %5.0f m elevation\n",
				run.Count, run.Distance/1000, formatDuration(run.MovingTime), run.ElevationGain)
		}
		if ride.Count > 0 {
			fmt.Printf("    Ride: %3d activities  %8.1f km  %s  %5.0f m elevation\n",
				ride.Count, ride.Distance/1000, formatDuration(ride.MovingTime), ride.ElevationGain)
		}
		if swim.Count > 0 {
			fmt.Printf("    Swim: %3d activities  %8.1f km  %s\n",
				swim.Count, swim.Distance/1000, formatDuration(swim.MovingTime))
		}
		if run.Count == 0 && ride.Count == 0 && swim.Count == 0 {
			fmt.Printf("    No activities\n")
		}
		fmt.Println()
	}

	printTotals("Recent (last 4 weeks)", stats.RecentRunTotals, stats.RecentRideTotals, stats.RecentSwimTotals)
	printTotals("Year to date", stats.YTDRunTotals, stats.YTDRideTotals, stats.YTDSwimTotals)
	printTotals("All time", stats.AllRunTotals, stats.AllRideTotals, stats.AllSwimTotals)

	return nil
}
