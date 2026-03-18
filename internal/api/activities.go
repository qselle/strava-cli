package api

import (
	"context"
	"fmt"
	"net/url"
)

type Activity struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	Type               string  `json:"type"`
	SportType          string  `json:"sport_type"`
	Distance           float64 `json:"distance"`
	MovingTime         int     `json:"moving_time"`
	ElapsedTime        int     `json:"elapsed_time"`
	TotalElevationGain float64 `json:"total_elevation_gain"`
	StartDate          string  `json:"start_date"`
	StartDateLocal     string  `json:"start_date_local"`
	Timezone           string  `json:"timezone"`
	AverageSpeed       float64 `json:"average_speed"`
	MaxSpeed           float64 `json:"max_speed"`
	AverageHeartrate   float64 `json:"average_heartrate"`
	MaxHeartrate       float64 `json:"max_heartrate"`
	HasHeartrate       bool    `json:"has_heartrate"`
	AverageWatts       float64 `json:"average_watts"`
	Calories           float64 `json:"calories"`
	KudosCount         int     `json:"kudos_count"`
	CommentCount       int     `json:"comment_count"`
	AchievementCount   int     `json:"achievement_count"`
	PRCount            int     `json:"pr_count"`
	SufferScore        float64 `json:"suffer_score"`
	GearID             string  `json:"gear_id"`
	DeviceName         string  `json:"device_name"`
	Manual             bool    `json:"manual"`
	Private            bool    `json:"private"`
}

type ListActivitiesParams struct {
	Before  int64
	After   int64
	Page    int
	PerPage int
}

func (c *Client) ListActivities(ctx context.Context, params ListActivitiesParams) ([]Activity, error) {
	q := url.Values{}
	if params.Before > 0 {
		q.Set("before", fmt.Sprintf("%d", params.Before))
	}
	if params.After > 0 {
		q.Set("after", fmt.Sprintf("%d", params.After))
	}
	if params.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", params.Page))
	}
	perPage := params.PerPage
	if perPage <= 0 {
		perPage = 30
	}
	q.Set("per_page", fmt.Sprintf("%d", perPage))

	var activities []Activity
	if err := c.getJSON(ctx, "/athlete/activities", q, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func (c *Client) GetActivity(ctx context.Context, id int64) (*Activity, error) {
	var activity Activity
	if err := c.getJSON(ctx, fmt.Sprintf("/activities/%d", id), nil, &activity); err != nil {
		return nil, err
	}
	return &activity, nil
}
