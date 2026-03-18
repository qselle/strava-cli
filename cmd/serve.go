package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

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

		srv := &http.Server{Addr: serveHTTP, Handler: handler}

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		fmt.Fprintf(cmd.ErrOrStderr(), "Starting MCP server on %s\n", serveHTTP)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	}

	return s.Run(cmd.Context(), &mcp.StdioTransport{})
}
