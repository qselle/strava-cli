package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/qselle/strava-cli/internal/server"
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
	s := server.NewServer()

	if serveHTTP != "" {
		handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return s
		}, nil)
		fmt.Fprintf(cmd.ErrOrStderr(), "Starting MCP server on %s\n", serveHTTP)
		return http.ListenAndServe(serveHTTP, handler)
	}

	return s.Run(context.Background(), &mcp.StdioTransport{})
}
