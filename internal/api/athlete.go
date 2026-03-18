package api

import "context"

type Athlete struct {
	ID            int64   `json:"id"`
	Username      string  `json:"username"`
	FirstName     string  `json:"firstname"`
	LastName      string  `json:"lastname"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Country       string  `json:"country"`
	Sex           string  `json:"sex"`
	Premium       bool    `json:"premium"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Weight        float64 `json:"weight"`
	ProfileMedium string  `json:"profile_medium"`
	Profile       string  `json:"profile"`
}

func (c *Client) GetAthlete(ctx context.Context) (*Athlete, error) {
	var athlete Athlete
	if err := c.getJSON(ctx, "/athlete", nil, &athlete); err != nil {
		return nil, err
	}
	return &athlete, nil
}
