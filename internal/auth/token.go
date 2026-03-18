package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func (t *Token) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "strava-cli")
	return dir, os.MkdirAll(dir, 0700)
}

func tokenPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "token.json"), nil
}

func SaveToken(token *Token) error {
	path, err := tokenPath()
	if err != nil {
		return fmt.Errorf("getting token path: %w", err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling token: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

func LoadToken() (*Token, error) {
	path, err := tokenPath()
	if err != nil {
		return nil, fmt.Errorf("getting token path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in — run 'strava auth' first")
		}
		return nil, fmt.Errorf("reading token: %w", err)
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	return &token, nil
}

// GetValidToken loads the token and refreshes it if expired.
// Credentials for refresh are read from the stored token file.
func GetValidToken(ctx context.Context) (*Token, error) {
	token, err := LoadToken()
	if err != nil {
		return nil, err
	}

	if !token.IsExpired() {
		return token, nil
	}

	if token.ClientID == "" || token.ClientSecret == "" {
		return nil, fmt.Errorf("token expired and no client credentials stored — run 'strava-cli auth' again")
	}

	cfg := OAuthConfig{
		ClientID:     token.ClientID,
		ClientSecret: token.ClientSecret,
	}

	refreshed, err := RefreshAccessToken(ctx, cfg, token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("token expired and refresh failed: %w", err)
	}

	// Preserve stored credentials
	refreshed.ClientID = token.ClientID
	refreshed.ClientSecret = token.ClientSecret

	if err := SaveToken(refreshed); err != nil {
		return nil, fmt.Errorf("saving refreshed token: %w", err)
	}

	return refreshed, nil
}

func ClearToken() error {
	path, err := tokenPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
