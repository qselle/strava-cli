package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/qselle/strava-cli/internal/auth"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "strava-cli",
	Short: "Strava CLI — check your activities, stats, and streaks",
	Long:  "A command-line interface and MCP server for the Strava API.\nTrack your activities, monitor your streak, and let your AI agents keep you accountable.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}

func getOAuthConfig() auth.OAuthConfig {
	clientID := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Fprintln(os.Stderr, "Error: STRAVA_CLIENT_ID and STRAVA_CLIENT_SECRET environment variables are required.")
		fmt.Fprintln(os.Stderr, "Get them at https://www.strava.com/settings/api")
		os.Exit(1)
	}

	return auth.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}
