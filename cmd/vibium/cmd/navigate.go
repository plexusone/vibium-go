package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var navigateTimeout time.Duration

var navigateCmd = &cobra.Command{
	Use:   "go <url>",
	Short: "Navigate to a URL",
	Long: `Navigate the browser to the specified URL.

Examples:
  vibium go https://example.com
  vibium go https://google.com --timeout 30s`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), navigateTimeout)
		defer cancel()

		vibe := mustGetVibe(ctx)

		if err := vibe.Go(ctx, url); err != nil {
			return fmt.Errorf("navigation failed: %w", err)
		}

		title, _ := vibe.Title(ctx)
		fmt.Printf("Navigated to: %s\n", url)
		if title != "" {
			fmt.Printf("Title: %s\n", title)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(navigateCmd)
	navigateCmd.Flags().DurationVar(&navigateTimeout, "timeout", 30*time.Second, "Navigation timeout")
}
