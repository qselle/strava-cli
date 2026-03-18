package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/qselle/strava-cli/internal/auth"
)

var authManual bool

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Strava",
	Long:  "Log in to Strava via OAuth2.\nBy default, opens your browser for authorization.\nUse --manual on headless servers to paste the code manually.",
	RunE:  runAuth,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored Strava credentials",
	RunE:  runLogout,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE:  runAuthStatus,
}

func init() {
	authCmd.Flags().BoolVar(&authManual, "manual", false, "Manually paste the authorization code (for headless servers)")
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(authCmd)
}

func runAuth(cmd *cobra.Command, args []string) error {
	cfg := getOAuthConfig()
	ctx := context.Background()

	var token *auth.Token
	var err error
	if authManual {
		token, err = auth.LoginManual(ctx, cfg)
	} else {
		token, err = auth.LoginBrowser(ctx, cfg)
	}
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Store credentials alongside token for auto-refresh
	token.ClientID = cfg.ClientID
	token.ClientSecret = cfg.ClientSecret

	if err := auth.SaveToken(token); err != nil {
		return fmt.Errorf("saving token: %w", err)
	}

	if jsonOutput {
		out := map[string]any{"status": "authenticated"}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	fmt.Println("Authenticated successfully!")
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := auth.ClearToken(); err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	if jsonOutput {
		out := map[string]any{"status": "logged_out"}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	fmt.Println("Logged out.")
	return nil
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	token, err := auth.LoadToken()
	if err != nil {
		if jsonOutput {
			out := map[string]any{"authenticated": false}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(out)
		}
		fmt.Println("Not authenticated. Run 'strava auth' to log in.")
		return nil
	}

	if jsonOutput {
		out := map[string]any{
			"authenticated": true,
			"expired":       token.IsExpired(),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	status := "valid"
	if token.IsExpired() {
		status = "expired (will auto-refresh)"
	}
	fmt.Printf("Authenticated (token: %s)\n", status)
	return nil
}
