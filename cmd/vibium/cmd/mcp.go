package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agentplexus/vibium-go/mcp"
	"github.com/spf13/cobra"
)

var (
	mcpHeadless       bool
	mcpDefaultTimeout time.Duration
	mcpProject        string
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server",
	Long: `Start the Vibium MCP (Model Context Protocol) server.

The MCP server provides browser automation tools for AI assistants.
It communicates via stdio using the MCP protocol.

Examples:
  vibium mcp
  vibium mcp --headless
  vibium mcp --timeout 60s`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle interrupt
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-sigCh
			fmt.Fprintln(os.Stderr, "\nShutting down MCP server...")
			cancel()
		}()

		config := mcp.Config{
			Headless:       mcpHeadless,
			DefaultTimeout: mcpDefaultTimeout,
			Project:        mcpProject,
		}

		server := mcp.NewServer(config)
		defer func() {
			if err := server.Close(context.Background()); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: cleanup error: %v\n", err)
			}
		}()

		if verbose {
			fmt.Fprintln(os.Stderr, "Starting Vibium MCP server...")
			if mcpHeadless {
				fmt.Fprintln(os.Stderr, "Mode: headless")
			}
		}

		return server.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
	mcpCmd.Flags().BoolVar(&mcpHeadless, "headless", false, "Run browser in headless mode")
	mcpCmd.Flags().DurationVar(&mcpDefaultTimeout, "timeout", 30*time.Second, "Default timeout for operations")
	mcpCmd.Flags().StringVar(&mcpProject, "project", "", "Project name for test reports")
}
