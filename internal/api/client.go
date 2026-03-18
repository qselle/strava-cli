package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseURL = "https://www.strava.com/api/v3"

type Client struct {
	httpClient  *http.Client
	accessToken string
}

func NewClient(accessToken string) *Client {
	return &Client{
		httpClient:  &http.Client{},
		accessToken: accessToken,
	}
}

func (c *Client) get(ctx context.Context, path string, params url.Values) ([]byte, error) {
	u := BaseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) getJSON(ctx context.Context, path string, params url.Values, v any) error {
	body, err := c.get(ctx, path, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
