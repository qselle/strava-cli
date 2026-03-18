package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveHTTP string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long:  "Start a Model Context Protocol (MCP) server that exposes Strava tools.\nDefault transport is stdio. Use --http to start an HTTP/SSE server.",
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().StringVar(&serveHTTP, "http", "", "Start HTTP/SSE server on this address (e.g. :8080)")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	// TODO: implement MCP server using modelcontextprotocol/go-sdk
	fmt.Println("MCP server not yet implemented. Coming soon!")
	fmt.Println("Track progress at https://github.com/qselle/strava-cli")
	return nil
}
