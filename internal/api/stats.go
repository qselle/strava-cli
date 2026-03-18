package api

import (
	"context"
	"fmt"
)

type ActivityTotal struct {
	Count            int     `json:"count"`
	Distance         float64 `json:"distance"`
	MovingTime       float64 `json:"moving_time"`
	ElapsedTime      float64 `json:"elapsed_time"`
	ElevationGain    float64 `json:"elevation_gain"`
	AchievementCount int     `json:"achievement_count,omitempty"`
}

type Stats struct {
	BiggestRideDistance       float64       `json:"biggest_ride_distance"`
	BiggestClimbElevationGain float64      `json:"biggest_climb_elevation_gain"`
	RecentRideTotals         ActivityTotal `json:"recent_ride_totals"`
	RecentRunTotals          ActivityTotal `json:"recent_run_totals"`
	RecentSwimTotals         ActivityTotal `json:"recent_swim_totals"`
	YTDRideTotals            ActivityTotal `json:"ytd_ride_totals"`
	YTDRunTotals             ActivityTotal `json:"ytd_run_totals"`
	YTDSwimTotals            ActivityTotal `json:"ytd_swim_totals"`
	AllRideTotals            ActivityTotal `json:"all_ride_totals"`
	AllRunTotals             ActivityTotal `json:"all_run_totals"`
	AllSwimTotals            ActivityTotal `json:"all_swim_totals"`
}

func (c *Client) GetStats(ctx context.Context, athleteID int64) (*Stats, error) {
	var stats Stats
	if err := c.getJSON(ctx, fmt.Sprintf("/athletes/%d/stats", athleteID), nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
