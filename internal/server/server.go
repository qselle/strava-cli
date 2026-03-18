package server

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/qselle/strava-cli/internal/api"
	"github.com/qselle/strava-cli/internal/auth"
)

func NewServer() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "strava-cli",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_activities",
		Description: "List recent Strava activities. Returns activity name, type, distance, duration, date, and more.",
	}, makeGetActivities())

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_stats",
		Description: "Get athlete stats: recent, year-to-date, and all-time totals for runs, rides, and swims.",
	}, makeGetStats())

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_streak",
		Description: "Check the athlete's activity streak. Returns active days, rest days, current streak, and a motivational verdict. Use this to check if the user is moving their ass.",
	}, makeGetStreak())

	return s
}

type GetActivitiesInput struct {
	After  string `json:"after,omitempty" jsonschema:"description=Show activities after this date (YYYY-MM-DD)"`
	Before string `json:"before,omitempty" jsonschema:"description=Show activities before this date (YYYY-MM-DD)"`
	Limit  int    `json:"limit,omitempty" jsonschema:"description=Maximum number of activities to return (default 10)"`
}

type GetActivitiesOutput struct {
	Activities []ActivitySummary `json:"activities"`
}

type ActivitySummary struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	DistanceKm float64 `json:"distance_km"`
	Duration   string  `json:"duration"`
	Date       string  `json:"date"`
	Elevation  float64 `json:"elevation_m"`
	Calories   float64 `json:"calories"`
}

func makeGetActivities() func(context.Context, *mcp.CallToolRequest, GetActivitiesInput) (*mcp.CallToolResult, GetActivitiesOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input GetActivitiesInput) (*mcp.CallToolResult, GetActivitiesOutput, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, GetActivitiesOutput{}, err
		}

		params := api.ListActivitiesParams{
			PerPage: input.Limit,
		}
		if params.PerPage <= 0 {
			params.PerPage = 10
		}
		if input.After != "" {
			t, err := time.Parse("2006-01-02", input.After)
			if err != nil {
				return nil, GetActivitiesOutput{}, fmt.Errorf("invalid after date: %w", err)
			}
			params.After = t.Unix()
		}
		if input.Before != "" {
			t, err := time.Parse("2006-01-02", input.Before)
			if err != nil {
				return nil, GetActivitiesOutput{}, fmt.Errorf("invalid before date: %w", err)
			}
			params.Before = t.Unix()
		}

		activities, err := client.ListActivities(ctx, params)
		if err != nil {
			return nil, GetActivitiesOutput{}, err
		}

		summaries := make([]ActivitySummary, len(activities))
		for i, a := range activities {
			summaries[i] = ActivitySummary{
				Name:       a.Name,
				Type:       a.SportType,
				DistanceKm: a.Distance / 1000,
				Duration:   formatDuration(a.MovingTime),
				Date:       a.StartDateLocal[:10],
				Elevation:  a.TotalElevationGain,
				Calories:   a.Calories,
			}
		}

		return nil, GetActivitiesOutput{Activities: summaries}, nil
	}
}

type GetStatsInput struct{}

type GetStatsOutput struct {
	RecentRun  TotalSummary `json:"recent_run"`
	RecentRide TotalSummary `json:"recent_ride"`
	RecentSwim TotalSummary `json:"recent_swim"`
	YTDRun     TotalSummary `json:"ytd_run"`
	YTDRide    TotalSummary `json:"ytd_ride"`
	YTDSwim    TotalSummary `json:"ytd_swim"`
	AllRun     TotalSummary `json:"all_run"`
	AllRide    TotalSummary `json:"all_ride"`
	AllSwim    TotalSummary `json:"all_swim"`
}

type TotalSummary struct {
	Count      int     `json:"count"`
	DistanceKm float64 `json:"distance_km"`
	Duration   string  `json:"duration"`
	ElevationM float64 `json:"elevation_m"`
}

func makeGetStats() func(context.Context, *mcp.CallToolRequest, GetStatsInput) (*mcp.CallToolResult, GetStatsOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input GetStatsInput) (*mcp.CallToolResult, GetStatsOutput, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, GetStatsOutput{}, err
		}

		athlete, err := client.GetAthlete(ctx)
		if err != nil {
			return nil, GetStatsOutput{}, err
		}

		stats, err := client.GetStats(ctx, athlete.ID)
		if err != nil {
			return nil, GetStatsOutput{}, err
		}

		toSummary := func(t api.ActivityTotal) TotalSummary {
			return TotalSummary{
				Count:      t.Count,
				DistanceKm: t.Distance / 1000,
				Duration:   formatDuration(t.MovingTime),
				ElevationM: t.ElevationGain,
			}
		}

		return nil, GetStatsOutput{
			RecentRun:  toSummary(stats.RecentRunTotals),
			RecentRide: toSummary(stats.RecentRideTotals),
			RecentSwim: toSummary(stats.RecentSwimTotals),
			YTDRun:     toSummary(stats.YTDRunTotals),
			YTDRide:    toSummary(stats.YTDRideTotals),
			YTDSwim:    toSummary(stats.YTDSwimTotals),
			AllRun:     toSummary(stats.AllRunTotals),
			AllRide:    toSummary(stats.AllRideTotals),
			AllSwim:    toSummary(stats.AllSwimTotals),
		}, nil
	}
}

type GetStreakInput struct {
	Days int `json:"days,omitempty" jsonschema:"description=Number of days to look back (default 7)"`
}

type GetStreakOutput struct {
	Period        string   `json:"period"`
	ActiveDays    int      `json:"active_days"`
	TotalDays     int      `json:"total_days"`
	RestDays      int      `json:"rest_days"`
	CurrentStreak int      `json:"current_streak"`
	TotalDistance  float64  `json:"total_distance_km"`
	TotalTime     string   `json:"total_time"`
	Activities    int      `json:"activities"`
	Verdict       string   `json:"verdict"`
	Motivation    string   `json:"motivation"`
}

func makeGetStreak() func(context.Context, *mcp.CallToolRequest, GetStreakInput) (*mcp.CallToolResult, GetStreakOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input GetStreakInput) (*mcp.CallToolResult, GetStreakOutput, error) {
		days := input.Days
		if days <= 0 {
			days = 7
		}

		client, err := getClient(ctx)
		if err != nil {
			return nil, GetStreakOutput{}, err
		}

		now := time.Now()
		after := now.AddDate(0, 0, -days)

		activities, err := client.ListActivities(ctx, api.ListActivitiesParams{
			After:   after.Unix(),
			PerPage: 200,
		})
		if err != nil {
			return nil, GetStreakOutput{}, err
		}

		activeDays := make(map[string]bool)
		var totalDistance float64
		var totalTime int

		for _, a := range activities {
			date := a.StartDateLocal[:10]
			activeDays[date] = true
			totalDistance += a.Distance / 1000
			totalTime += a.MovingTime
		}

		currentStreak := 0
		for i := 0; i < days; i++ {
			day := now.AddDate(0, 0, -i).Format("2006-01-02")
			if activeDays[day] {
				currentStreak++
			} else {
				break
			}
		}

		verdict, motivation := getMotivation(len(activeDays), days, currentStreak)

		return nil, GetStreakOutput{
			Period:        fmt.Sprintf("last %d days", days),
			ActiveDays:    len(activeDays),
			TotalDays:     days,
			RestDays:      days - len(activeDays),
			CurrentStreak: currentStreak,
			TotalDistance:  totalDistance,
			TotalTime:     formatDuration(totalTime),
			Activities:    len(activities),
			Verdict:       verdict,
			Motivation:    motivation,
		}, nil
	}
}

func getClient(ctx context.Context) (*api.Client, error) {
	token, err := auth.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("not authenticated — run 'strava-cli auth' first: %w", err)
	}
	return api.NewClient(token.AccessToken), nil
}

func formatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
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
